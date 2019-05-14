package main

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

type errInvalidComparison struct {
	Comparison string
}

func (e *errInvalidComparison) Error() string {
	return fmt.Sprintf("invalid comparison %s", e.Comparison)
}

func initFlags(flag *pflag.FlagSet) {

	// AWS Flags
	cli.InitAWSFlags(flag)

	// Vault Flags
	cli.InitVaultFlags(flag)

	flag.String("comparison", "size", "Comparison used against files, either 'size' or 'md5'")
	flag.Int64("max-object-size", 10, "The maximum size of files to download in MB for use with md5 comparison")

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
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

func checkComparison(v *viper.Viper) error {
	if c := v.GetString("comparison"); len(c) == 0 || !stringSliceContains([]string{"size", "md5"}, c) {
		return errors.Wrap(&errInvalidComparison{Comparison: c}, fmt.Sprintf("%s is invalid", "comparison"))
	}
	return nil
}

func checkConfig(v *viper.Viper) error {

	region, err := cli.CheckAWSRegion(v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid", cli.AWSRegionFlag))
	}

	if err := cli.CheckAWSRegionForService(region, s3.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, s3.ServiceName))
	}

	if err := cli.CheckVault(v); err != nil {
		return err
	}

	if err := checkComparison(v); err != nil {
		return err
	}

	return nil
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("Usage of compare-secure-migrations:")
	if flag != nil {
		flag.PrintDefaults()
	}
	os.Exit(1)
}

func main() {
	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		quit(logger, flag, err)
	}

	v := viper.New()
	pflagsErr := v.BindPFlags(flag)
	if pflagsErr != nil {
		quit(logger, flag, err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	verbose := v.GetBool(cli.VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	awsConfig, err := cli.GetAWSConfig(v, verbose)
	if err != nil {
		quit(logger, nil, err)
	}

	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "failed to create AWS session"))
	}

	serviceS3 := s3.New(sess)
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
		resp, err := serviceS3.ListObjects(&s3.ListObjectsInput{
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
