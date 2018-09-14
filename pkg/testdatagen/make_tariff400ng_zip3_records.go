package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeTariff400ngZip3 finds or makes a single Tariff400ngZip3 record
func MakeTariff400ngZip3(db *pop.Connection, assertions Assertions) models.Tariff400ngZip3 {
	zip3 := models.Tariff400ngZip3{
		Zip3:          "902",
		BasepointCity: "Beverly Hills",
		State:         "CA",
		ServiceArea:   "56",
		RateArea:      "US88",
		Region:        "2",
	}

	mergeModels(&zip3, assertions.Tariff400ngZip3)

	var zip3s models.Tariff400ngZip3s
	err := db.Where("zip3 = ?", zip3.Zip3).All(&zip3s)
	if err != nil {
		log.Panic(err)
	}

	if len(zip3s) == 0 {
		mustCreate(db, &zip3)
		return zip3
	}

	return zip3s[0]
}

// MakeDefaultTariff400ngZip3 makes a Tariff400ngZip3 record with default values
func MakeDefaultTariff400ngZip3(db *pop.Connection) models.Tariff400ngZip3 {
	return MakeTariff400ngZip3(db, Assertions{})
}
