package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/ghcdieselfuelprice"
)

// Command: go run github.com/transcom/mymove/cmd/save_ghc_diesel_fuel_price_data
func saveGHCDieselFuelPriceData(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkSaveGHCDieselFuelPriceConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	eiaURL := v.GetString(cli.EIAURLFlag)
	eiaKey := v.GetString(cli.EIAKeyFlag)
	dieselFuelPriceStorer := ghcdieselfuelprice.NewDieselFuelPriceStorer(eiaURL, eiaKey, ghcdieselfuelprice.FetchEIAData)

	dieselFuelPriceStorer.Run()

	fmt.Println(dieselFuelPriceStorer.EIAURL)

	return nil
}

func checkSaveGHCDieselFuelPriceConfig(v *viper.Viper, logger logger) error {

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

func initSaveGHCDieselFuelPriceFlags(flag *pflag.FlagSet) {

	//DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}
