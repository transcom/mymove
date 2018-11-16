package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngItemRateDefaultValidDate provides a date for which default item rates will be valid
var Tariff400ngItemRateDefaultValidDate = time.Date(2018, time.August, 15, 0, 0, 0, 0, time.UTC)

// Tariff400ngItemRateEffectiveDateLower provides a standard lower date
var Tariff400ngItemRateEffectiveDateLower = time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC)

// Tariff400ngItemRateEffectiveDateUpper provides a standard upper date
var Tariff400ngItemRateEffectiveDateUpper = time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC)

// MakeTariff400ngItemRate creates a single Tariff400ngItemRate record
func MakeTariff400ngItemRate(db *pop.Connection, assertions Assertions) models.Tariff400ngItemRate {
	rate := models.Tariff400ngItemRate{
		Code:               "105B",
		Schedule:           nil,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(2147483647),
		RateCents:          unit.Cents(1000),
		EffectiveDateLower: Tariff400ngItemRateEffectiveDateLower,
		EffectiveDateUpper: Tariff400ngItemRateEffectiveDateUpper,
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
