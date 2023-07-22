package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from the command line with go run ./cmd/milmove-tasks connect-to-gex-via-sftp

func initConnectToGEXViaSFTPFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// GEX SFTP
	cli.InitGEXSFTPFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func connectToGEXViaSFTP(_ *cobra.Command, _ []string) error {
	v := viper.New()

	logger, _, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}
	zap.ReplaceGlobals(logger)

	flag := pflag.CommandLine
	initConnectToGEXViaSFTPFlags(flag)
	err = flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("failed to parse flags", zap.Error(err))
	}

	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatal("failed to bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

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

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	// SSH and SFTP Connection Setup
	sshClient, err := cli.InitGEXSSH(appCtx, v)
	if err != nil {
		logger.Error("couldn't initialize SSH client", zap.Error(err))
		return err
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	sftpClient, err := cli.InitGEXSFTP(appCtx, sshClient)
	if err != nil {
		logger.Error("couldn't initialize SFTP client", zap.Error(err))
		return err
	}
	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	pwd, err := sftpClient.Getwd()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully connected via SFTP. The present working directory is %v", pwd)

	return nil
}
