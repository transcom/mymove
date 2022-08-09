package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pg "github.com/habx/pg-commands"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func initExportDBDataFlags(flag *pflag.FlagSet) {
	cli.InitDatabaseFlags(flag)
	cli.InitAWSFlags(flag)
	cli.InitStorageFlags(flag)
	cli.InitLoggingFlags(flag)
	flag.String("bucket-name", "mybucket", "Name of S3 bucket to store db dump")
	flag.String("aws-access-key-id", "1234", "access key id of AWS user")
	flag.String("aws-secret-access-key", "1234", "secret access key of AWS user")
	flag.SortFlags = false
}

func exportDBData(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return fmt.Errorf("could not parse args: %w", err)
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return fmt.Errorf("could not bind flags: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, _, err := logging.Config(
		logging.WithEnvironment(dbEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	pgConfig := getPGConfig(*v)
	dumpExec := createDBDump(pgConfig)
	if dumpExec.Error != nil {
		logger.Error("error in pg_dump")
		return dumpExec.Error.Err
	}
	logger.Info("dump success created file " + dumpExec.File)

	wd, _ := os.Getwd()
	err = exportToS3Bucket(dumpExec.File, wd, v.GetString("bucket-name"))
	if err != nil {
		logger.Error("error in upload to S3 bucket")
		return err
	}
	fmt.Println("upload to S3 bucket successful")
	return nil
}

func getPGConfig(v viper.Viper) pg.Postgres {
	return pg.Postgres{
		Host:     v.GetString(cli.DbHostFlag),
		Port:     v.GetInt(cli.DbPortFlag),
		DB:       v.GetString(cli.DbNameFlag),
		Username: v.GetString(cli.DbUserFlag),
		Password: v.GetString(cli.DbPasswordFlag),
	}
}

func createDBDump(pgConfig pg.Postgres) pg.Result {
	dump := pg.NewDump(&pgConfig)
	return dump.Exec(pg.ExecOptions{StreamPrint: false})
}

func exportToS3Bucket(fileName string, dir string, bucketName string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	return err
}
