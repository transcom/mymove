package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeGHCDieselFuelPrice(db *pop.Connection, assertions Assertions) models.GHCDieselFuelPrice {

	ghcDieselFuelPrice := models.GHCDieselFuelPrice{
		FuelPriceInMillicents: unit.Millicents(243300),
		PublicationDate:       time.Date(GHCTestYear, time.July, 20, 0, 0, 0, 0, time.UTC),
	}

	mergeModels(&ghcDieselFuelPrice, assertions.GHCDieselFuelPrice)

	mustCreate(db, &ghcDieselFuelPrice, assertions.Stub)

	return ghcDieselFuelPrice
}
