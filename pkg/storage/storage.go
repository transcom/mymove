package storage

import (
	/*
		#nosec - we use md5 because it's required by the S3 API for
		validating data integrity.
		https://aws.amazon.com/premiumsupport/knowledge-center/data-integrity-s3/
	*/
	"crypto/md5"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// StoreResult represents the result of a call to Store().
type StoreResult struct{}

// FileStorer is the set of methods needed to store and retrieve objects.
type FileStorer interface {
	Store(string, io.ReadSeeker, string) (*StoreResult, error)
	Fetch(string) (io.ReadCloser, error)
	Delete(string) error
	PresignedURL(string, string) (string, error)
	FileSystem() *afero.Afero
	TempFileSystem() *afero.Afero
}

// ComputeChecksum calculates the MD% checksum for the provided data. It expects that
// the passed io object will be seeked to its beginning and will seek back to the
// beginning after reading its content.
func ComputeChecksum(data io.ReadSeeker) (string, error) {
	/*
		#nosec - we use md5 because it's required by the S3 API for
		validating data integrity.
		https://aws.amazon.com/premiumsupport/knowledge-center/data-integrity-s3/
	*/
	hash := md5.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", errors.Wrap(err, "could not read file")
	}

	if _, err := data.Seek(0, io.SeekStart); err != nil { // seek back to beginning of file
		return "", errors.Wrap(err, "could not seek to beginning of file")
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// DetectContentType leverages http.DetectContentType to identify the content type
// of the provided data. It expects that the passed io object will be seeked to its
// beginning and will seek back to the beginning after reading its content.
func DetectContentType(data io.ReadSeeker) (string, error) {
	// Start by seeking to beginning
	data.Seek(0, io.SeekStart)

	buffer := make([]byte, 512)
	if _, err := data.Read(buffer); err != nil {
		return "", errors.Wrap(err, "could not read first bytes of file")
	}

	contentType := http.DetectContentType(buffer)

	if _, err := data.Seek(0, io.SeekStart); err != nil { // seek back to beginning of file
		return "", errors.Wrap(err, "could not seek to beginning of file")
	}
	return contentType, nil
}
