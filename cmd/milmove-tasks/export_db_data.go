package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobuffalo/pop/v6"
	pg "github.com/habx/pg-commands"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func checkExportDBDataConfig(v *viper.Viper, logger *zap.Logger) error {
	logger.Debug("checking config")
	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}
	return nil
}

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

	err = checkExportDBDataConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metadata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
		}
	}
	dbConnectionDetails, err := cli.BuildConnectionDetails(v, dbCreds, logger)
	if err != nil {
		logger.Error("Error building the db connection details")
		return err
	}

	dumpExec := createDBDump(dbConnectionDetails, logger)
	if dumpExec.Error != nil {
		logger.Error("Error in pg_dump")
		return dumpExec.Error.Err
	}
	logger.Info("Dump success created file " + dumpExec.File)

	wd, _ := os.Getwd()
	err = exportToS3Bucket(dumpExec.File, wd, v.GetString("bucket-name"))
	if err != nil {
		logger.Error("Error in upload to S3 bucket")
		return err
	}
	fmt.Println("Upload to S3 bucket successful")
	return nil
}

func getPGConfig(dbConnectionDetails *pop.ConnectionDetails, logger *zap.Logger) (pg.Postgres, error) {
	port, err := strconv.Atoi(dbConnectionDetails.Port)
	if err != nil {
		logger.Error("Port must be a valid integer")
	}
	return pg.Postgres{
		Host:     dbConnectionDetails.Host,
		Port:     port,
		DB:       dbConnectionDetails.Database,
		Username: dbConnectionDetails.User,
		Password: dbConnectionDetails.Password,
	}, err
}

func createDBDump(dbConnectionDetails *pop.ConnectionDetails, logger *zap.Logger) pg.Result {
	pgConfig, err := getPGConfig(dbConnectionDetails, logger)
	if err != nil {
		logger.Error("Error building pg config from db connection details")
	}
	dump := pg.NewDump(&pgConfig)
	return dump.Exec(pg.ExecOptions{StreamPrint: false})
}

func exportToS3Bucket(fileName string, dir string, bucketName string) error {
	sess, err := awssession.NewSession()
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
