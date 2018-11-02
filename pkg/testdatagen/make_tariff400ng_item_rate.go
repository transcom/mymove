package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeTariff400ngItemRate creates a single Tariff400ngItemRate record
func MakeTariff400ngItemRate(db *pop.Connection, assertions Assertions) models.Tariff400ngItemRate {
	rate := models.Tariff400ngItemRate{
		Code:               "105B",
		ServicesSchedule:   nil,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(2147483647),
		RateCents:          1000,
		EffectiveDateLower: time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		EffectiveDateUpper: time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC),
	}

	// Overwrite values with those from assertions
	mergeModels(&rate, assertions.Tariff400ngItemRate)

	mustCreate(db, &rate)

	return rate
}

// MakeDefaultTariff400ngItemRate makes a 400ng item rate with default values
func MakeDefaultTariff400ngItemRate(db *pop.Connection) models.Tariff400ngItemRate {
	return MakeTariff400ngItemRate(db, Assertions{})
}
