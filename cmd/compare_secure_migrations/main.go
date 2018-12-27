package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

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

func hashFileMd5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func download(downloader *s3manager.Downloader, bucket string, item string) {
	file, err := os.Create("deleteme.tmp")
	if err != nil {
		exitError("Unable to open file %q, %v", err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitError("Unable to download item %v", err)
	}

	hash, err := hashFileMd5("deleteme.tmp")
	if err != nil {
		exitError("Unable to hash file", err)
	}
	fmt.Println("Downloaded", item, numBytes, "bytes", hash, "hash")
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

	for _, bucket := range bucketNames {
		resp, err := s3Service.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Marker: aws.String("secure-migrations"),
		})
		if err != nil {
			exitError("Failed to List Objects", err)
		}

		for _, item := range resp.Contents {
			download(downloader, bucket, *item.Key)
		}
	}
}
