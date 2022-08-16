package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

const (
	//BucketNameFlag is the flag for the name of the S3 bucket to which the DB dump will be sent
	BucketNameFlag string = "bucket-name"
	//TmpDirFlag is the flag for the temp directory to save the dump file; must have write access
	TmpDirFlag string = "tmp-dir"
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
	flag.String(BucketNameFlag, "mybucket", "Name of S3 bucket to store db dump")
	flag.String("aws-access-key-id", "1234", "access key id of AWS user")
	flag.String("aws-secret-access-key", "1234", "secret access key of AWS user")
	flag.String(TmpDirFlag, "/tmp", "absolute path of temporary directory to save dump file; must have write access")
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

	dbConnectionDetails, err := getDBConnectionDetails(v, logger)
	if err != nil {
		logger.Error("Error getting DB connection details")
		return err
	}

	filePath := getDBDumpFilePath(v)

	doLoggingForSadDevs(logger)

	err = createDBDump(dbConnectionDetails, filePath, logger)
	if err != nil {
		logger.Error("Error creating database dump")
		return err
	}
	logger.Info("Dump success created file " + filePath)

	err = exportToS3Bucket(filePath, v.GetString(BucketNameFlag))
	if err != nil {
		logger.Error("Error in upload to S3 bucket")
		return err
	}
	fmt.Println("Upload to S3 bucket successful")
	return nil
}

func doLoggingForSadDevs(logger *zap.Logger) {
	cmd1 := exec.Command("ls", "-fal", "/var")
	cmd2 := exec.Command("ls", "-fal", "/var/db-export")
	cmd3 := exec.Command("which", "pg_dump")
	cmd4 := exec.Command("type", "-p", "pg_dump")
	cmd5 := exec.Command("echo", "$PATH")

	err := cmd1.Start()
	if err != nil {
		logger.Error("Error with ls for /var dir")
	}
	err = cmd2.Start()
	if err != nil {
		logger.Error("Error with ls for /var/db-export")
	}
	err = cmd3.Start()
	if err != nil {
		logger.Error("Error with which for pg_dump")
	}
	err = cmd4.Start()
	if err != nil {
		logger.Error("Error with type for pg_dump")
	}
	err = cmd5.Start()
	if err != nil {
		logger.Error("Error with echo for $PATH")
	}
}

func getDBConnectionDetails(v *viper.Viper, logger *zap.Logger) (*pop.ConnectionDetails, error) {
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
	}
	return dbConnectionDetails, err
}

func getDBDumpFilePath(v *viper.Viper) string {
	fileName := fmt.Sprintf("milmove-database-export_%s.sql", time.Now().Format(time.RFC3339))
	dir := v.GetString(TmpDirFlag)
	return filepath.Join(dir, fileName)
}

func createDBDump(dbConn *pop.ConnectionDetails, filePath string, logger *zap.Logger) error {
	dbname := fmt.Sprintf("--dbname=postgres://%s:%s@%s:%s/%s",
		dbConn.User,
		dbConn.Password,
		dbConn.Host,
		dbConn.Port,
		dbConn.Database)
	cmd := exec.Command("pg_dump", dbname)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating file at path %s", filePath))
		return err
	}
	cmd.Stdout = file
	err = cmd.Start()
	if err != nil {
		logger.Error("Error executing pg_dump command")
		return err
	}
	return err
}

func exportToS3Bucket(filePath string, bucketName string) error {
	sess, err := awssession.NewSession()
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	return err
}
