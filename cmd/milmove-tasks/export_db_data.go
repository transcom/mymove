package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pg "github.com/habx/pg-commands"
	"github.com/spf13/cobra"
)

func exportDBData(cmd *cobra.Command, args []string) error {

	pgConfig := getPGConfig()
	dump := pg.NewDump(&pgConfig)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		return dumpExec.Error.Err
	}

	fmt.Println("Dump success")
	fmt.Println(dumpExec.Output)

	// export to S3 bucket
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	file, err := os.Open(dumpExec.File)
	if err != nil {
		return err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(getEnvOrPanic("BUCKET_NAME")),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(dumpExec.File),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err != nil {
		return err
	}
	fmt.Println("Upload to S3 bucket successful")
	return nil
}

func getPGConfig() pg.Postgres {
	host := getEnvOrPanic("PGHOST")
	port, err := strconv.Atoi(getEnvOrPanic("PGPORT"))
	if err != nil {
		panic("PGPORT must be an integer")
	}
	db := getEnvOrPanic("PGDB")
	user := getEnvOrPanic("PGUSER")
	password := getEnvOrPanic("PGPASSWORD")

	return pg.Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		Username: user,
		Password: password,
	}
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if len(value) <= 0 {
		panic(fmt.Sprintf("config loading failed; required environment variable %s must be set", key))
	}
	return value
}
