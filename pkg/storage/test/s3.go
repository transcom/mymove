package test

import (
	"fmt"
	"io"
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
	PutFiles    []PutFile
	willSucceed bool
}

// Key returns a key built using the specified path fragments.
func (fake *FakeS3Storage) Key(args ...string) string {
	return path.Join(args...)
}

// Delete removes a file.
func (fake *FakeS3Storage) Delete(key string) error {
	itemIndex := -1
	for i, f := range fake.PutFiles {
		if f.Key == key {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		return errors.New("can't delete item that doesn't exist")
	}
	// Remove file from putFiles
	fake.PutFiles = append(fake.PutFiles[:itemIndex], fake.PutFiles[itemIndex+1:]...)
	return nil
}

// Store stores a file.
func (fake *FakeS3Storage) Store(key string, data io.ReadSeeker, md5 string) (*storage.StoreResult, error) {
	file := PutFile{
		Key:      key,
		Body:     data,
		Checksum: md5,
	}
	fake.PutFiles = append(fake.PutFiles, file)
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

// PresignedURL returns a URL that can be used to retrieve a file.
func (fake *FakeS3Storage) PresignedURL(key string, contentType string) (string, error) {
	url := fmt.Sprintf("https://example.com/dir/%s?contentType=%s&signed=test", key, contentType)
	return url, nil
}

// NewFakeS3Storage creates a new FakeS3Storage for testing purposes.
func NewFakeS3Storage(willSucceed bool) *FakeS3Storage {
	return &FakeS3Storage{
		willSucceed: willSucceed,
	}
}
