package storage

import (
	"io"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// S3 implements the file storage API using S3.
type S3 struct {
	bucket       string
	keyNamespace string
	logger       Logger
	client       *s3.S3
	fs           *afero.Afero
	tempFs       *afero.Afero
}

// NewS3 creates a new S3 using the provided AWS session.
func NewS3(bucket string, keyNamespace string, logger Logger, session *session.Session) *S3 {
	var fs = afero.NewMemMapFs()
	var tempFs = afero.NewMemMapFs()
	client := s3.New(session)
	return &S3{
		bucket:       bucket,
		keyNamespace: keyNamespace,
		logger:       logger,
		client:       client,
		fs:           &afero.Afero{Fs: fs},
		tempFs:       &afero.Afero{Fs: tempFs},
	}
}

// Store stores the content from an io.ReadSeeker at the specified key.
func (s *S3) Store(key string, data io.ReadSeeker, checksum string) (*StoreResult, error) {
	if key == "" {
		return nil, errors.New("A valid StorageKey must be set before data can be uploaded")
	}

	namespacedKey := path.Join(s.keyNamespace, key)

	input := &s3.PutObjectInput{
		Bucket:     &s.bucket,
		Key:        &namespacedKey,
		Body:       data,
		ContentMD5: &checksum,
	}

	if _, err := s.client.PutObject(input); err != nil {
		return nil, errors.Wrap(err, "put on S3 failed")
	}

	return &StoreResult{}, nil
}

// Delete deletes an object at a specified key
func (s *S3) Delete(key string) error {
	namespacedKey := path.Join(s.keyNamespace, key)

	input := &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &namespacedKey,
	}

	_, err := s.client.DeleteObject(input)
	if err != nil {
		return errors.Wrap(err, "delete on S3 failed")
	}

	return nil
}

// Fetch retrieves an object at a specified key and stores it in a tempfile. The
// path to this file is returned.
//
// It is the caller's responsibility to cleanup this file.
func (s *S3) Fetch(key string) (io.ReadCloser, error) {
	namespacedKey := path.Join(s.keyNamespace, key)

	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &namespacedKey,
	}

	getObjectOutput, err := s.client.GetObject(input)
	if err != nil {
		return nil, errors.Wrap(err, "get object on S3 failed")
	}

	return getObjectOutput.Body, nil
}

// FileSystem returns the underlying afero filesystem
func (s *S3) FileSystem() *afero.Afero {
	return s.fs
}

// TempFileSystem returns the temporary afero filesystem
func (s *S3) TempFileSystem() *afero.Afero {
	return s.tempFs
}

// PresignedURL returns a URL that provides access to a file for 15 minutes.
func (s *S3) PresignedURL(key string, contentType string) (string, error) {
	namespacedKey := path.Join(s.keyNamespace, key)

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              &s.bucket,
		Key:                 &namespacedKey,
		ResponseContentType: &contentType,
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "could not generate presigned URL")
	}
	return url, nil
}
