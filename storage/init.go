package storage

import (
	"cloud.google.com/go/storage"
	"errors"
	"io"
)

type Storage interface {
	// Write returns a new writer for a specific file on cloud storage
	Write(filepath string, options ...Option) (io.WriteCloser, error)
	// Read reads a file from online storage
	Read(filepath string, options ...Option) (io.ReadCloser, error)
	// Close closes the storage connection
	Close() error
	// Checks if a file exists on the cloud storage
	IsExists(filepath string, options ...Option) (bool, error)
	// GetSignedURL returns a URL for the specified object. Signed URLs allow
	// the users access to a restricted resource for a limited time without signing in
	GetSignedURL(filepath string, opts *SignedURLOptions) (string, error)
	// ListObject returns list of object name inside the bucket
	// query can be used to filter object name
	ListObject(query *storage.Query, options ...Option) ([]string, error)
}

func NewStorage(impl Implementation, options ...Option) (Storage, error) {
	c := &config{}
	for _, o := range options {
		o(c)
	}

	switch impl {
	case GCP:
		return newGCPStorage(c)
	default:
		return nil, errors.New("implementation not found")
	}
}
