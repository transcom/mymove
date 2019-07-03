package migrate

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ListFiles lists the files in a given directory.
func ListFiles(p string, s3Client *s3.S3) ([]string, error) {
	if strings.HasPrefix(p, "file://") {
		f, err := os.Open(p[len("file://"):])
		if err != nil {
			return []string{}, err
		}
		return f.Readdirnames(-1)
	} else if strings.HasPrefix(p, "s3://") {
		parts := strings.SplitN(p[len("s3://"):], "/", 2)
		if len(parts) == 2 {
			bucket := parts[0]
			prefix := parts[1]
			filenames := make([]string, 0)
			var marker *string
			for {
				listObjectsOutput, err := s3Client.ListObjects(&s3.ListObjectsInput{
					Bucket: aws.String(bucket),
					Prefix: aws.String(prefix),
					Marker: marker,
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
				marker = listObjectsOutput.Contents[len(listObjectsOutput.Contents)-1].Key
			}
			return filenames, nil
		}
	}
	return make([]string, 0), nil
}
