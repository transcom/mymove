package storage

import (
	"io"
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
	input := &s3.PutObjectInput{
		Bucket:     &s.bucket,
		Key:        &key,
		Body:       data,
		ContentMD5: &checksum,
	}

	if _, err := s.client.PutObject(input); err != nil {
		return nil, errors.Wrap(err, "put to S3 failed")
	}

	return &StoreResult{}, nil
}

// Delete deletes an object at a specified key
func (s *S3) Delete(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}

	_, err := s.client.DeleteObject(input)
	if err != nil {
		return errors.Wrap(err, "delete to S3 failed")
	}

	return nil
}

// Key returns a joined key plus any global namespace
func (s *S3) Key(args ...string) string {
	args = append([]string{s.keyNamespace}, args...)
	return path.Join(args...)
}

// PresignedURL returns a URL that provides access to a file for 15 minutes.
func (s *S3) PresignedURL(key string, contentType string) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              &s.bucket,
		Key:                 &key,
		ResponseContentType: &contentType,
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "could not generate presigned URL")
	}
	return url, nil
}
