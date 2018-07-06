package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/storage"
)

// PutFile represents a file that has been uploaded to a FakeS3Storage.
type PutFile struct {
	Key      string
	Body     io.ReadSeeker
	Checksum string
}

// FakeS3Storage is used for local testing to stub out calls to S3.
type FakeS3Storage struct {
	PutFiles    map[string]PutFile
	willSucceed bool
}

// Key returns a key built using the specified path fragments.
func (fake *FakeS3Storage) Key(args ...string) string {
	return path.Join(args...)
}

// Delete removes a file.
func (fake *FakeS3Storage) Delete(key string) error {
	if _, ok := fake.PutFiles[key]; !ok {
		return errors.New("can't delete item that doesn't exist")
	}

	delete(fake.PutFiles, key)
	return nil
}

// Store stores a file.
func (fake *FakeS3Storage) Store(key string, data io.ReadSeeker, md5 string) (*storage.StoreResult, error) {
	file := PutFile{
		Key:      key,
		Body:     data,
		Checksum: md5,
	}
	fake.PutFiles[key] = file
	buf := []byte{}
	_, err := data.Read(buf)
	if err != nil {
		return nil, err
	}
	if fake.willSucceed {
		return &storage.StoreResult{}, nil
	}
	return nil, errors.New("failed to push")
}

// Fetch retrieves a copy of a file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (fake *FakeS3Storage) Fetch(key string) (string, error) {
	outputFile, err := ioutil.TempFile(os.TempDir(), "filesystem")
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer outputFile.Close()

	file := fake.PutFiles[key]
	_, err = io.Copy(outputFile, file.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return outputFile.Name(), nil
}

// PresignedURL returns a URL that can be used to retrieve a file.
func (fake *FakeS3Storage) PresignedURL(key string, contentType string) (string, error) {
	url := fmt.Sprintf("https://example.com/dir/%s?contentType=%s&signed=test", key, contentType)
	return url, nil
}

// NewFakeS3Storage creates a new FakeS3Storage for testing purposes.
func NewFakeS3Storage(willSucceed bool) *FakeS3Storage {
	return &FakeS3Storage{
		willSucceed: willSucceed,
		PutFiles:    make(map[string]PutFile),
	}
}
