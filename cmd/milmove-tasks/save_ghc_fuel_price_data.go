package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
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

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil, nil)

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
	return nil
}
