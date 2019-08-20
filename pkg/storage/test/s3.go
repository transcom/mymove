package test

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/storage"
)

// FakeS3Storage is used for local testing to stub out calls to S3.
type FakeS3Storage struct {
	willSucceed bool
	fs          *afero.Afero
	tempFs      *afero.Afero
}

// Delete removes a file.
func (fake *FakeS3Storage) Delete(key string) error {
	f, err := fake.fs.Open(key)
	if err != nil {
		return err
	}

	return f.Close()
}

// Store stores a file.
func (fake *FakeS3Storage) Store(key string, data io.ReadSeeker, md5 string) (*storage.StoreResult, error) {
	f, err := fake.fs.Create(key)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, data)
	if err != nil {
		return nil, err
	}

	if fake.willSucceed {
		return &storage.StoreResult{}, nil
	}
	return nil, errors.New("failed to push")
}

// Fetch returns the file at the given key
func (fake *FakeS3Storage) Fetch(key string) (io.ReadCloser, error) {
	return fake.fs.Open(key)
}

// PresignedURL returns a URL that can be used to retrieve a file.
func (fake *FakeS3Storage) PresignedURL(key string, contentType string) (string, error) {
	url := fmt.Sprintf("https://example.com/dir/%s?contentType=%s&signed=test", key, contentType)
	return url, nil
}

// FileSystem returns the underlying afero filesystem
func (fake *FakeS3Storage) FileSystem() *afero.Afero {
	return fake.fs
}

// TempFileSystem returns the underlying afero filesystem
func (fake *FakeS3Storage) TempFileSystem() *afero.Afero {
	return fake.tempFs
}

// NewFakeS3Storage creates a new FakeS3Storage for testing purposes.
func NewFakeS3Storage(willSucceed bool) *FakeS3Storage {
	var fs = afero.NewMemMapFs()
	var tempFs = afero.NewMemMapFs()

	return &FakeS3Storage{
		willSucceed: willSucceed,
		fs:          &afero.Afero{Fs: fs},
		tempFs:      &afero.Afero{Fs: tempFs},
	}
}
