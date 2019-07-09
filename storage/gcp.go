package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
)

const (
	APISCOPE = "https://www.googleapis.com/auth/cloud-platform"
)

type accountInfo struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

type gcpStorage struct {
	bucket     *blob.Bucket
	bucketName string
	gcsClient  *storage.Client
	*accountInfo
}

func newGCPStorage(cfg *config) (*gcpStorage, error) {
	jsonFile, err := os.Open(cfg.storageSecret)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var info accountInfo
	err = json.Unmarshal(jsonData, &info)
	if err != nil {
		return nil, err
	}

	// validate permission + bucket existence
	gcsClient, err := storage.NewClient(cfg.initCtx, option.WithCredentialsFile(cfg.storageSecret))
	if err != nil {
		return nil, err
	}
	if _, err = gcsClient.Bucket(cfg.storageBucket).Attrs(cfg.initCtx); err != nil {
		return nil, err
	}

	gcpCredentials, err := google.CredentialsFromJSON(cfg.initCtx, jsonData, APISCOPE)
	if err != nil {
		return nil, err
	}

	httpClient, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(gcpCredentials))
	if err != nil {
		return nil, err
	}

	bucket, err := gcsblob.OpenBucket(cfg.initCtx, httpClient, cfg.storageBucket, nil)
	if err != nil {
		return nil, err
	}

	return &gcpStorage{
		bucket:      bucket,
		bucketName:  cfg.storageBucket,
		accountInfo: &info,
		gcsClient:   gcsClient,
	}, nil
}

func (m *gcpStorage) Read(filepath string, options ...Option) (io.ReadCloser, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	var ctx context.Context
	if c.operationCtx == nil {
		ctx = context.Background()
	} else {
		ctx = c.operationCtx
	}

	reader, err := m.bucket.NewReader(ctx, filepath, c.gcpReaderOptions)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (m *gcpStorage) Close() error {
	return m.bucket.Close()
}

func (m *gcpStorage) Write(filepath string, options ...Option) (io.WriteCloser, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	var ctx context.Context
	if c.operationCtx == nil {
		ctx = context.Background()
	} else {
		ctx = c.operationCtx
	}

	writer, err := m.bucket.NewWriter(ctx, filepath, c.gcpWriterOptions)

	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (m *gcpStorage) IsExists(filepath string, options ...Option) (bool, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	var ctx context.Context
	if c.operationCtx == nil {
		ctx = context.Background()
	} else {
		ctx = c.operationCtx
	}

	exists, err := m.bucket.Exists(ctx, filepath)
	return exists, err
}

func (m *gcpStorage) GetSignedURL(filepath string, opts *SignedURLOptions) (string, error) {
	if opts == nil {
		return "", errors.New("must specify signed URL options")
	}

	signedURL, err := storage.SignedURL(m.bucketName, filepath, &storage.SignedURLOptions{
		GoogleAccessID: m.accountInfo.ClientEmail,
		PrivateKey:     []byte(m.accountInfo.PrivateKey),
		Method:         opts.HTTPMethod,
		Expires:        opts.ExpiryTime,
		ContentType:    opts.ContentType,
	})
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (m *gcpStorage) ListObject(query *storage.Query, options ...Option) ([]string, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	var ctx context.Context
	if c.operationCtx == nil {
		ctx = context.Background()
	} else {
		ctx = c.operationCtx
	}

	it := m.gcsClient.Bucket(m.bucketName).Objects(ctx, query)
	var res []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, attrs.Name)
	}
	return res, nil
}
