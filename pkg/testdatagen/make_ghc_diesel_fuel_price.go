package testdatagen

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofrs/uuid"

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

func FetchOrMakeGHCDieselFuelPrice(db *pop.Connection, assertions Assertions) models.GHCDieselFuelPrice {
	var existingGHCDieselFuelPrice models.GHCDieselFuelPrice
	if !assertions.GHCDieselFuelPrice.PublicationDate.IsZero() {
		err := db.Where("publication_date = ?", assertions.GHCDieselFuelPrice.PublicationDate).First(&existingGHCDieselFuelPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for existing GHCDieselFuelPrice by publication date", err)
		}

		if existingGHCDieselFuelPrice.ID != uuid.Nil {
			return existingGHCDieselFuelPrice
		}
	}

	return MakeGHCDieselFuelPrice(db, assertions)
}
