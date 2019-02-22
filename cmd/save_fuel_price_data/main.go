package main

import (
	"log"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/services/fuelprice"
)

// Command: go run cmd/save_fuel_price_data/main.go
func main() {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}
	clock := clock.New()
	fuelPrices := fuelprice.DieselFuelPriceStorer{
		DB:            db,
		Clock:         clock,
		FetchFuelData: fuelprice.FetchFuelPriceData,
	}

	verrs, err := fuelPrices.StoreFuelPrices(12)
	if err != nil || verrs != nil {
		log.Fatal(err, verrs)
	}
}
