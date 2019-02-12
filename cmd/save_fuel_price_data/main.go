package main

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/service/fuelprice"
)

func main() {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}

	addFuelPrices := fuelprice.AddFuelDieselPrices{
		DB: db,
	}

	verrs, err := addFuelPrices.Call()
	if err != nil || verrs != nil {
		log.Fatal(err, verrs)
	}
}
