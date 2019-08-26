package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func getIamCreds(v *viper.Viper, logger *zap.Logger) (*credentials.Credentials, error) {
	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) || (v.GetString(cli.EmailBackendFlag) == "ses") || (v.GetString(cli.StorageBackendFlag) == "s3") {
		c, errorConfig := cli.GetAWSConfig(v, v.GetBool(cli.VerboseFlag))
		if errorConfig != nil {
			logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		// We want to get the credentials from the logged in AWS session rather than create directly,
		// because the session conflates the environment, shared, and container metdata config
		// within NewSession.  With stscreds, we use the Secure Token Service,
		// to assume the given role (that has rds db connect permissions).
		dbIamRole := v.GetString(cli.DbIamRoleFlag)
		logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
		dbCreds = stscreds.NewCredentials(session, dbIamRole)
	}

	return dbCreds, nil
}

func makeConnection(v *viper.Viper, dbCreds *credentials.Credentials, logger *zap.Logger) (*pop.Connection, error) {
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		if dbConnection == nil {
			// No connection object means that the configuraton failed to validate and we should kill server startup
			logger.Fatal("Invalid DB Configuration", zap.Error(err))
		} else {
			// A valid connection object that still has an error indicates that the DB is not up and
			// thus is not ready for migrations
			logger.Fatal("DB is not ready for connections", zap.Error(err))
		}
	}

	logger.Info("DB connection successful")
	return dbConnection, nil

}

// evalActiveConnection calls a DB query every minute for 20 minutes
func evalActiveConnection(v *viper.Viper, logger *zap.Logger) {

	creds, err := getIamCreds(v, logger)
	if err != nil {
		log.Fatalf("Failed to get IAM creds")
	}

	// Lets create a conneciton to evaluate with
	dbConn, err := makeConnection(v, creds, logger)
	if err != nil {
		log.Fatalf("Failed to properly connect to database")
	}

	runtime := 18

	logger.Info(fmt.Sprintf("Starting evalActiveConnection test, one query a minute for %d minutes", runtime))
	ticker := time.NewTicker(time.Minute)
	count := 0
	go func() {
		for range ticker.C {
			err = dbConn.RawQuery("SELECT datname from pg_database").Exec()
			if err == nil {
				count = count + 1
				logger.Info(fmt.Sprintf("Ran query successfully at %d minute", count))
				dbConn.Close()
			} else {
				log.Fatalf(err.Error())
				return
			}
		}
	}()

	time.Sleep(time.Duration(runtime) * time.Minute)
	ticker.Stop()

	logger.Info("evalActiveConnection is complete")

}

// evalActiveCredentials this tests opening a connection with 15m credentials
func evalActiveCredentials(v *viper.Viper, logger *zap.Logger) {
	creds, err := getIamCreds(v, logger)
	if err != nil {
		log.Fatalf("Failed to get IAM creds")
	}

	runtime := time.Duration(1080) // 30 hours
	period := time.Duration(1)
	count := 0

	logger.Info(fmt.Sprintf("Starting evalActiveCredentials test, use existing credential to open connection every %dm for %dm", period, runtime))

	ticker := time.NewTicker(period * time.Second)
	go func() {
		for range ticker.C {

			// Lets create a conneciton to evaluate with
			dbConn, err := makeConnection(v, creds, logger)
			if err != nil {
				log.Fatalf("Failed to properly connect to database")
				return
			}

			err = dbConn.RawQuery("SELECT datname from pg_database").Exec()
			if err == nil {
				count = count + 1
				logger.Info(fmt.Sprintf("Ran query successfully at %d minute", count))
				err = dbConn.Close()
				if err != nil {
					logger.Fatal(fmt.Sprintf("Error closing db connection: %s", err))
				}
			} else {
				log.Fatalf(err.Error())
				return
			}
		}
	}()

	time.Sleep(runtime * time.Minute)
	ticker.Stop()

	logger.Info("evalActiveCredentials is complete")

}

func main() {
	logger, err := logging.Config("test", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	flag := pflag.CommandLine
	cli.InitDatabaseFlags(flag)
	err = flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatalf("Failed parsing args")
	}

	evalActiveCredentials(v, logger)
	evalActiveConnection(v, logger)

}
