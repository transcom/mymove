package main

import (
	"fmt"
	"log"
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
	"github.com/transcom/mymove/pkg/iampostgres"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/ghcdieselfuelprice"
)

func checkSaveGHCFuelPriceConfig(v *viper.Viper, logger *zap.Logger) error {

	logger.Debug("checking config")

	err := cli.CheckEIA(v)
	if err != nil {
		return err
	}

	err = cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initSaveGHCFuelPriceFlags(flag *pflag.FlagSet) {

	//DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Command: go run github.com/transcom/mymove/cmd/milmove-tasks/save_ghc_fuel_price_data
func saveGHCFuelPriceData(cmd *cobra.Command, args []string) error {

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
	zap.ReplaceGlobals(logger)

	err = checkSaveGHCFuelPriceConfig(v, logger)
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

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	eiaURL := v.GetString(cli.EIAURLFlag)
	eiaKey := v.GetString(cli.EIAKeyFlag)
	newDieselFuelPriceInfo := ghcdieselfuelprice.NewDieselFuelPriceInfo(eiaURL, eiaKey, ghcdieselfuelprice.FetchEIAData, logger)

	err = newDieselFuelPriceInfo.RunFetcher(appCtx)
	if err != nil {
		logger.Fatal("error returned by RunFetcher function in ghcdieselfuelprice service", zap.Error(err))
	}

	err = newDieselFuelPriceInfo.RunStorer(appCtx)
	if err != nil {
		logger.Fatal("error returned by RunStorer function in ghcdieselfuelprice service", zap.Error(err))
	}

	if v.GetBool(cli.DbIamFlag) {
		iampostgres.ShutdownIAM()
	}

	return dbConnection.Close()
}
