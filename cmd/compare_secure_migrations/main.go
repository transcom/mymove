package main

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

type errInvalidComparison struct {
	Comparison string
}

func (e *errInvalidComparison) Error() string {
	return fmt.Sprintf("invalid comparison %s", e.Comparison)
}

func initFlags(flag *pflag.FlagSet) {
	flag.String("aws-s3-region", "", "AWS region used for S3 file storage")
	flag.String("comparison", "size", "Comparison used against files, either 'size' or 'md5'")
	flag.Int64("max-object-size", 10, "The maximum size of files to download in MB for use with md5 comparison")
}

func hashObjectMd5(buff *aws.WriteAtBuffer) string {

	// Sum the bytes in the buffer
	// #nosec
	hashInBytes := md5.Sum(buff.Bytes())

	// Convert the bytes to a string
	return hex.EncodeToString(hashInBytes[:16])
}

func stringSliceContains(stringSlice []string, value string) bool {
	for _, x := range stringSlice {
		if value == x {
			return true
		}
	}
	return false
}

func checkRegion(v *viper.Viper) error {

	regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.S3ServiceID)
	if !ok {
		return fmt.Errorf("could not find regions for service %s", endpoints.S3ServiceID)
	}

	r := v.GetString("aws-s3-region")
	if len(r) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
	}

	if _, ok := regions[r]; !ok {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
	}

	return nil
}

func checkComparison(v *viper.Viper) error {
	if c := v.GetString("comparison"); len(c) == 0 || !stringSliceContains([]string{"size", "md5"}, c) {
		return errors.Wrap(&errInvalidComparison{Comparison: c}, fmt.Sprintf("%s is invalid", "comparison"))
	}
	return nil
}

func checkConfig(v *viper.Viper) error {

	err := checkRegion(v)
	if err != nil {
		return err
	}

	err = checkComparison(v)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err := checkConfig(v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	compareMap := map[string]map[string]string{}

	var wg sync.WaitGroup
	var errs []error

	comparison := v.GetString("comparison")
	maxObjectSize := v.GetInt64("max-object-size")

	for _, bucket := range bucketNames {
		resp, err := s3Service.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Marker: aws.String("secure-migrations"),
		})
		if err != nil {
			errs = append(errs, err)
		}

		// Download and compare all objects concurrently for this bucket
		for _, item := range resp.Contents {
			key := *item.Key

			// Ignore directories and migrations larger than maxObjectSize
			if comparison == "md5" && (*item.Size == 0 || *item.Size >= maxObjectSize*1024*1024) {
				fmt.Println("SKIP", bucket, key, "size", *item.Size, "bytes")
				continue
			}

			if _, ok := compareMap[key]; !ok {
				compareMap[key] = map[string]string{}
			}

			if comparison == "md5" {
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
				}(downloader, bucket, key, compareMap, &wg)
			} else {
				compareMap[key][bucket] = strconv.FormatInt(*item.Size, 10)
			}
		}
	}

	wg.Wait()

	// Compare files and print differences
	for migration, v := range compareMap {
		if v[bucketNames[0]] == v[bucketNames[1]] && v[bucketNames[1]] == v[bucketNames[2]] {
			fmt.Println("Migration:", migration)
			fmt.Println("\t", fmt.Sprintf("%s:", comparison), v[bucketNames[0]])
		} else {
			fmt.Println("Migration:", migration)
			for _, bucket := range bucketNames {
				comp := v[bucket]
				if comp != "" {
					fmt.Println("\t", bucket, "\t", comp)
				} else {
					fmt.Println("\t", bucket, "\tNot Found")
				}
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
