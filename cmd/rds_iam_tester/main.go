package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"

	// "github.com/pkg/errors"
	// "github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func main() {
	logger, err := logging.Config("test", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	flag := pflag.CommandLine
	cli.InitDatabaseFlags(flag)
	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatalf("Failed parsing arg blags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	var dbCreds *credentials.Credentials
	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
		// We want to get the credentials from the logged in AWS session rather than create directly,
		// because the session conflates the environment, shared, and container metdata config
		// within NewSession.  With stscreds, we use the Secure Token Service,
		// to assume the given role (that has rds db connect permissions).
		dbIamRole := v.GetString(cli.DbIamRoleFlag)
		logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
		dbCreds = stscreds.NewCredentials(session, dbIamRole)
	}

	// dbConnection, errDbConnection := cli.InitDatabase(v, dbCreds, logger)
	cli.InitDatabase(v, dbCreds, logger)

}
