package storage

import (
	"io"
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
	bucket string
	logger *zap.Logger
	client *s3.S3
}

// NewS3 creates a new S3 using the provided AWS session.
func NewS3(bucket string, logger *zap.Logger, session *session.Session) *S3 {
	client := s3.New(session)
	return &S3{bucket, logger, client}
}

// Store stores the content from an io.ReadSeeker at the specified key.
func (s *S3) Store(key string, data io.ReadSeeker, md5 string) (*StoreResult, error) {
	input := &s3.PutObjectInput{
		Bucket:     &s.bucket,
		Key:        &key,
		Body:       data,
		ContentMD5: &md5,
	}

	_, err := s.client.PutObject(input)
	if err != nil {
		return nil, errors.Wrap(err, "put to S3 failed")
	}
	return &StoreResult{}, nil
}

// Key returns a joined key plus any global namespace
func (s *S3) Key(args ...string) string {
	namespace := os.Getenv("AWS_S3_KEY_NAMESPACE")
	args = append([]string{namespace}, args...)
	return path.Join(args...)
}

// PresignedURL returns a URL that provides access to a file for 15 mintes.
func (s *S3) PresignedURL(key string) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "could not generate presigned URL")
	}
	return url, nil
}
