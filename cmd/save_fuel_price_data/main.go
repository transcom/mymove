package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/facebookgo/clock"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/fuelprice"
)

func checkConfig(v *viper.Viper, logger logger) error {

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

func initFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Command: go run github.com/transcom/mymove/cmd/save_fuel_price_data
func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
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
	if session != nil {
		// We want to get the credentials from the logged in AWS session rather than create directly,
		// because the session conflates the environment, shared, and container metdata config
		// within NewSession.  With stscreds, we use the Secure Token Service,
		// to assume the given role (that has rds db connect permissions).
		dbCreds = stscreds.NewCredentials(session, v.GetString(cli.DbIamRoleFlag))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	clock := clock.New()
	eiaKey := v.GetString(cli.EIAKeyFlag)
	eiaURL := v.GetString(cli.EIAURLFlag)
	fuelPrices := fuelprice.NewDieselFuelPriceStorer(
		dbConnection,
		logger,
		clock,
		fuelprice.FetchFuelPriceData,
		eiaKey,
		eiaURL,
	)

	verrs, err := fuelPrices.StoreFuelPrices(12)
	if err != nil || verrs.HasAny() {
		log.Fatal(err, verrs)
	}
}
