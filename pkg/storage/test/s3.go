package test

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/text/encoding/charmap"

	"github.com/transcom/mymove/pkg/storage"
)

// FakeS3Storage is used for local testing to stub out calls to S3.
type FakeS3Storage struct {
	willSucceed  bool
	fs           *afero.Afero
	tempFs       *afero.Afero
	tagsToReturn map[string]string
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
func (fake *FakeS3Storage) Store(key string, data io.ReadSeeker, _ string, _ *string) (*storage.StoreResult, error) {
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
	if !fake.willSucceed {
		return nil, errors.New("failed to fetch file")
	}

	return fake.fs.Open(key)
}

// PresignedURL returns a URL that can be used to retrieve a file.
func (fake *FakeS3Storage) PresignedURL(key string, contentType string, filename string) (string, error) {
	filenameBuffer := make([]byte, 0)

	// Double quote the filename to be able to handle filenames with commas in them
	quotedFilename := strconv.Quote(filename)

	for _, r := range quotedFilename {
		if encodedRune, ok := charmap.ISO8859_1.EncodeRune(r); ok {
			filenameBuffer = append(filenameBuffer, encodedRune)
		}
	}

	contentDisposition := "attachment"
	if len(filenameBuffer) > 0 {
		contentDisposition += "; filename=" + string(filenameBuffer)
	}

	values := url.Values{}
	values.Add("response-content-type", contentType)
	values.Add("response-content-disposition", contentDisposition)
	values.Add("signed", "test")
	url := fmt.Sprintf("https://example.com/dir/%s?", key) + values.Encode()
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

// Tags returns the tags for a specified key
// The following tags must be present for S3 user uploading
// "av-status": "CLEAN"
// OR
// "av-status": "INFECTED"
// The upload function will hang, attempting to poll the uploaded file until one of these
// tags are found
func (fake *FakeS3Storage) Tags(_ string) (map[string]string, error) {
	if fake.tagsToReturn == nil {
		// Default to clean file upload
		return map[string]string{"av-status": "CLEAN"}, nil
	}
	return fake.tagsToReturn, nil
}

func (_ *FakeS3Storage) StorageType() string {
	return "S3"
}

// NewFakeS3Storage creates a new FakeS3Storage for testing purposes.
func NewFakeS3Storage(willSucceed bool, tagsToReturn map[string]string) *FakeS3Storage {
	var fs = afero.NewMemMapFs()
	var tempFs = afero.NewMemMapFs()

	return &FakeS3Storage{
		willSucceed:  willSucceed,
		fs:           &afero.Afero{Fs: fs},
		tempFs:       &afero.Afero{Fs: tempFs},
		tagsToReturn: tagsToReturn,
	}
}
