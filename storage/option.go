package storage

import (
	"context"
	"time"

	"gocloud.dev/blob"
)

// Option initializes configs
type Option func(c *config)

type SignedURLOptions struct {
	HTTPMethod  string
	ContentType string
	ExpiryTime  time.Time
}

type config struct {
	// context for initializing connection to cloud storage
	initCtx       context.Context
	// context for each of operation (Read, Write, Exists, etc...)
	operationCtx  context.Context
	storageSecret string
	storageBucket string

	// GCP specific settings
	gcpWriterOptions *blob.WriterOptions
	gcpReaderOptions *blob.ReaderOptions
}

// context used for reading, writing, etc...
func OperationCtx(opCtx context.Context) Option {
	return func(c *config) {
		c.operationCtx = opCtx
	}
}

func GCPStorage(initCtx context.Context, storageBucket, storageSecretFilepath string) Option {
	return func(c *config) {
		c.initCtx = initCtx
		c.storageSecret = storageSecretFilepath
		c.storageBucket = storageBucket
	}
}

func GCPReaderOptions(gcpReaderOptions *blob.ReaderOptions) Option {
	return func(c *config) {
		c.gcpReaderOptions = gcpReaderOptions
	}
}

func GCPWriterOptions(gcpWriterOptions *blob.WriterOptions) Option {
	return func(c *config) {
		c.gcpWriterOptions = gcpWriterOptions
	}
}
