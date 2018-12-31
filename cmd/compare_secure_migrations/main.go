package main

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"errors"
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
	flag.String("aws-s3-region", "", "AWS region used for S3 file storage")
	flag.Int64("max-object-size", 10, "The maximum size of files to download in MB")
}

func hashObjectMd5(buff *aws.WriteAtBuffer) string {

	// Sum the bytes in the buffer
	// #nosec
	hashInBytes := md5.Sum(buff.Bytes())

	// Convert the bytes to a string
	return hex.EncodeToString(hashInBytes[:16])
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
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 128 * 1024 * 1024 // 128MB per part
		d.Concurrency = 100
	})

	bucketNames := []string{
		"transcom-ppp-app-experimental-us-west-2",
		"transcom-ppp-app-staging-us-west-2",
		"transcom-ppp-app-prod-us-west-2",
	}

	hashCompare := map[string]map[string]string{}

	var wg sync.WaitGroup
	var errs []error

	maxObjectSize := v.GetInt64("max-object-size")

	for _, bucket := range bucketNames {
		resp, err := s3Service.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Marker: aws.String("secure-migrations"),
		})
		if err != nil {
			errs = append(errs, err)
		}

		// Download and hash all objects concurrently for this bucket
		for _, item := range resp.Contents {
			key := *item.Key

			// Ignore directories and migrations larger than maxObjectSize
			if *item.Size == 0 || *item.Size >= maxObjectSize*1024*1024 {
				fmt.Println("SKIP", bucket, key, "size", *item.Size, "bytes")
				continue
			}

			if _, ok := hashCompare[key]; !ok {
				hashCompare[key] = map[string]string{}
			}

			wg.Add(1)
			go func(downloader *s3manager.Downloader, bucket string, objectName string, compare map[string]map[string]string, wg *sync.WaitGroup) {
				var hash string

				// Create an in-memory buffer to write the object
				buff := &aws.WriteAtBuffer{}

				// Download the object to the buffer
				_, err := downloader.Download(buff,
					&s3.GetObjectInput{
						Bucket: aws.String(bucket),
						Key:    aws.String(objectName),
					})

				// Save errors for later
				if err != nil {
					errs = append(errs, err)
					fmt.Println("ERROR", bucket, objectName, err)
				}

				// Calculate the hash for the object
				hash = hashObjectMd5(buff)
				compare[objectName][bucket] = hash
				wg.Done()
			}(downloader, bucket, key, hashCompare, &wg)
		}
	}

	wg.Wait()

	// Compare hashses and print differences
	for migration, v := range hashCompare {
		if v[bucketNames[0]] == v[bucketNames[1]] && v[bucketNames[1]] == v[bucketNames[2]] {
			fmt.Println("Migration:", migration)
			fmt.Println("\t", "Hash:", v[bucketNames[0]])
		} else {
			fmt.Println("Migration:", migration)
			for _, bucket := range bucketNames {
				hash := v[bucket]
				fmt.Println("\t", bucket, "\t", hash)
			}
		}
	}

	// Print the errors and use exit code as length of list
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(len(errs))
	}
}
