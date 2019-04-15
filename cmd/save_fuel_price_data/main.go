package main

import (
	"log"
	"os"
	"strings"

	"github.com/facebookgo/clock"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/fuelprice"
)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Info("checking config")

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

// Command: go run cmd/save_fuel_price_data/main.go
func main() {

	flag := pflag.CommandLine

	flag.String("env", "development", "The environment to run in, which configures the database.")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	flag.SortFlags = false
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	env := v.GetString("env")

	logger, err := logging.Config(env, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	clock := clock.New()
	fuelPrices := fuelprice.NewDieselFuelPriceStorer(
		dbConnection,
		logger,
		clock,
		fuelprice.FetchFuelPriceData,
		v.GetString("eia-key"),
		v.GetString("eia-url"),
	)

	verrs, err := fuelPrices.StoreFuelPrices(12)
	if err != nil || verrs.HasAny() {
		log.Fatal(err, verrs)
	}
}
