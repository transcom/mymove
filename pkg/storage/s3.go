package storage

import (
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type S3 struct {
	bucket string
	logger *zap.Logger
	client *s3.S3
}

func NewS3(bucket string, logger *zap.Logger, session *session.Session) *S3 {
	client := s3.New(session)
	return &S3{bucket, logger, client}
}

func (s *S3) Store(key string, data io.ReadSeeker, md5 string) (*StoreResult, error) {
	input := &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   data,
	}

	_, err := s.client.PutObject(input)
	if err != nil {
		return nil, errors.Wrap(err, "put to S3 failed")
	}
	return &StoreResult{}, nil
}

func (s *S3) Key(args ...string) string {
	namespace := os.Getenv("AWS_S3_KEY_NAMESPACE")
	args = append([]string{namespace}, args...)
	return path.Join(args...)
}

func (s *S3) PresignedURL(key string) (string, error) {
	return "", nil
}
