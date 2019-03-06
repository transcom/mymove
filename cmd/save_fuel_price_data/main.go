package main

import (
	"log"
	"os"
	"strings"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/fuelprice"
)

// Command: go run cmd/save_fuel_price_data/main.go
func main() {

	flag := pflag.CommandLine

	flag.String("eia-key", "", "key for Energy Information Administration (EIA) api")
	flag.String("eia-url", "", "url for EIA api")
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	db, err := pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.Config("development", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	clock := clock.New()
	fuelPrices := fuelprice.NewDieselFuelPriceStorer(
		db,
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
