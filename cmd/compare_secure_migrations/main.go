package main

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initFlags(flag *pflag.FlagSet) {
	flag.String("aws_s3_region", "", "AWS region used for S3 file storage")
}

func exitError(msg string, err error) {
	fmt.Println(msg)
	fmt.Println(err)
	os.Exit(1)
}

func hashObjectMd5(buff *aws.WriteAtBuffer) string {

	// Sum the bytes in the buffer
	// #nosec
	hashInBytes := md5.Sum(buff.Bytes())

	// Convert the bytes to a string
	returnMD5String := hex.EncodeToString(hashInBytes[:16])

	return returnMD5String
}

func downloadAndHash(downloader *s3manager.Downloader, bucket string, objectName string) (string, error) {
	var hash string

	// Create an in-memory buffer to write the object
	buff := &aws.WriteAtBuffer{}

	// Download the object to the buffer
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectName),
		})
	if err != nil {
		return hash, err
	}

	// Calculate the hash for the object
	hash = hashObjectMd5(buff)

	return hash, err
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	sess := awssession.Must(awssession.NewSession(&aws.Config{
		Region: aws.String(v.GetString("aws-s3-region")),
	}))

	s3Service := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	bucketNames := []string{
		"transcom-ppp-app-experimental-us-west-2",
		"transcom-ppp-app-staging-us-west-2",
		"transcom-ppp-app-prod-us-west-2",
	}

	hashCompare := map[string]map[string]string{}

	var wg sync.WaitGroup
	for _, bucket := range bucketNames {
		resp, err := s3Service.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Marker: aws.String("secure-migrations"),
		})
		if err != nil {
			exitError("Failed to List Objects", err)
		}

		// Download and hash all objects concurrently for this bucket
		for _, item := range resp.Contents {
			key := *item.Key
			if _, ok := hashCompare[key]; !ok {
				hashCompare[key] = map[string]string{}
			}
			wg.Add(1)
			go func(downloader *s3manager.Downloader, bucket string, objectName string, compare map[string]map[string]string) {
				hash, err := downloadAndHash(downloader, bucket, objectName)
				compare[objectName][bucket] = hash
				if err != nil {
					exitError("Unable to download or hash file", err)
				}
				fmt.Println(bucket, objectName, hash)
			}(downloader, bucket, key, hashCompare)
		}
	}
	wg.Wait()

	fmt.Println(hashCompare)
}
