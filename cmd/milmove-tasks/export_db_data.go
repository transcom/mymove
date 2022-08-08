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
	dump := pg.NewDump(&pgConfig)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		return dumpExec.Error.Err
	}

	logger.Info("Dump success created file " + dumpExec.File)

	// export to S3 bucket
	sess, err := session.NewSession()
	if err != nil {
		logger.Error("could not create AWS session")
		return err
	}
	uploader := s3manager.NewUploader(sess)
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("could not get working directory")
		return err
	}
	filePath := filepath.Join(wd, dumpExec.File)
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("could not open file at path " + filePath)
		return err
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(v.GetString("bucket-name")),

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
		logger.Error("error in upload to S3 bucket")
		return err
	}
	fmt.Println("Upload to S3 bucket successful")
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
