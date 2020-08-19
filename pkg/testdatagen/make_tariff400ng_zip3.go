package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTariff400ngZip3 finds or makes a single Tariff400ngZip3 record
func MakeTariff400ngZip3(db *pop.Connection, assertions Assertions) models.Tariff400ngZip3 {
	zip3 := models.Tariff400ngZip3{
		Zip3:          DefaultZip3,
		BasepointCity: "Beverly Hills",
		State:         "CA",
		ServiceArea:   "56",
		RateArea:      "US88",
		Region:        "2",
	}

	mergeModels(&zip3, assertions.Tariff400ngZip3)

	mustCreate(db, &zip3)

	return zip3
}

// FetchOrMakeTariff400ngZip3 Tries fetching an existing zip3 first, then falls back to creating one
func FetchOrMakeTariff400ngZip3(db *pop.Connection, assertions Assertions) models.Tariff400ngZip3 {
	var existingZip3s models.Tariff400ngZip3s
	zip3 := DefaultZip3
	if assertions.Tariff400ngZip3.Zip3 != "" {
		zip3 = assertions.Tariff400ngZip3.Zip3
	}
	err := db.Where("zip3 = ?", zip3).All(&existingZip3s)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingZip3s) == 0 {
		return MakeTariff400ngZip3(db, assertions)
	}

	return existingZip3s[0]
}

// FetchOrMakeDefaultTariff400ngZip3 makes a Tariff400ngZip3 record with default values
func FetchOrMakeDefaultTariff400ngZip3(db *pop.Connection) models.Tariff400ngZip3 {
	return FetchOrMakeTariff400ngZip3(db, Assertions{})
}
