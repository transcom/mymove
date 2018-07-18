package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// S3 implements the file storage API using S3.
type S3 struct {
	bucket       string
	keyNamespace string
	logger       *zap.Logger
	client       *s3.S3
}

// NewS3 creates a new S3 using the provided AWS session.
func NewS3(bucket string, keyNamespace string, logger *zap.Logger, session *session.Session) *S3 {
	client := s3.New(session)
	return &S3{bucket, keyNamespace, logger, client}
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
func (s *S3) Fetch(key string) (string, error) {
	namespacedKey := path.Join(s.keyNamespace, key)

	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &namespacedKey,
	}

	getObjectOutput, err := s.client.GetObject(input)
	if err != nil {
		return "", errors.Wrap(err, "get object on S3 failed")
	}

	outputFile, err := ioutil.TempFile(os.TempDir(), "s3")
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer outputFile.Close()

	io.Copy(outputFile, getObjectOutput.Body)

	return outputFile.Name(), nil
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
