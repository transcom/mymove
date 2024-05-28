package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/afero"
)

// FileHelper is an afero filesystem struct
type FileHelper struct {
	fs afero.Fs
}

type S3ListObjectsV2API interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// NewFileHelper creates and returns a new File Helper
func NewFileHelper() *FileHelper {
	fs := afero.NewOsFs()
	return &FileHelper{
		fs: fs,
	}
}

// SetFileSystem sets the file system for the file helper
func (fh *FileHelper) SetFileSystem(fs afero.Fs) {
	fh.fs = fs
}

// ListFiles lists the files in a given directory.
func (fh *FileHelper) ListFiles(p string, s3Client S3ListObjectsV2API) ([]string, error) {
	if strings.HasPrefix(p, "file://") {
		f, err := fh.fs.Open(p[len("file://"):])
		if err != nil {
			return []string{}, err
		}
		return f.Readdirnames(-1)
	} else if strings.HasPrefix(p, "s3://") {
		if s3Client == nil {
			return make([]string, 0), fmt.Errorf("No s3Client provided to list files at %s", p)
		}
		parts := strings.SplitN(p[len("s3://"):], "/", 2)
		if len(parts) == 2 {
			bucket := parts[0]
			prefix := parts[1]
			filenames := make([]string, 0)
			var continuationToken *string
			for {
				listObjectsOutput, err := s3Client.ListObjectsV2(
					context.Background(),
					&s3.ListObjectsV2Input{
						Bucket:            aws.String(bucket),
						Prefix:            aws.String(prefix),
						ContinuationToken: continuationToken,
					})
				if err != nil {
					return filenames, err
				}
				for _, obj := range listObjectsOutput.Contents {
					key := *obj.Key
					filenames = append(filenames, key[len(prefix)+1:])
				}
				if listObjectsOutput.IsTruncated == nil || !*listObjectsOutput.IsTruncated {
					break
				}
				continuationToken = listObjectsOutput.NextContinuationToken
			}
			return filenames, nil
		}
	}
	return make([]string, 0), nil
}
