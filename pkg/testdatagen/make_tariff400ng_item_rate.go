package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeTariff400ngItemRate creates a single Tariff400ngItemRate record
func MakeTariff400ngItemRate(db *pop.Connection, assertions Assertions) models.Tariff400ngItemRate {
	rate := models.Tariff400ngItemRate{
		Code:               "105B",
		Schedule:           nil,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(2147483647),
		RateCents:          unit.Cents(1000),
		EffectiveDateLower: PeakRateCycleStart,
		EffectiveDateUpper: NonPeakRateCycleEnd,
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
