package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func intPointer(i int) *int {
	return &i
}

func (suite *ModelSuite) TestFetchTariff400ngItemRateBySchedule() {
	testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			ServicesSchedule: intPointer(1),
			RateCents:        unit.Cents(1001),
		},
	})
	rate2 := testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			ServicesSchedule: intPointer(2),
			RateCents:        unit.Cents(1002),
		},
	})
	testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			ServicesSchedule: intPointer(3),
			RateCents:        unit.Cents(1003),
		},
	})

	rate, err := models.FetchTariff400ngItemRate(suite.db, rate2.Code, *rate2.ServicesSchedule, 1000, time.Date(2018, time.August, 15, 0, 0, 0, 0, time.UTC))

	// Ensure we get back rate2's rate and not one for a different schedule
	if suite.NoError(err) {
		suite.Equal(rate.RateCents, rate2.RateCents)
	}
}

func (suite *ModelSuite) TestFetchTariff400ngItemRateNullSchedule() {
	rate1 := testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			ServicesSchedule: nil,
			RateCents:        unit.Cents(1001),
		},
	})

	rate, err := models.FetchTariff400ngItemRate(suite.db, rate1.Code, 3, 1000, time.Date(2018, time.August, 15, 0, 0, 0, 0, time.UTC))

	// Ensure we get back rate1's rate
	if suite.NoError(err) {
		suite.Equal(rate.RateCents, rate1.RateCents)
	}
}

func (suite *ModelSuite) TestFetchTariff400ngItemRateByWeight() {
	testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			WeightLbsLower: unit.Pound(1000),
			WeightLbsUpper: unit.Pound(1099),
			RateCents:      unit.Cents(1001),
		},
	})
	rate2 := testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			WeightLbsLower: unit.Pound(1100),
			WeightLbsUpper: unit.Pound(1199),
			RateCents:      unit.Cents(1002),
		},
	})
	testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			WeightLbsLower: unit.Pound(1200),
			WeightLbsUpper: unit.Pound(1299),
			RateCents:      unit.Cents(1003),
		},
	})
	testdatagen.MakeTariff400ngItemRate(suite.db, testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			Code:           "other code",
			WeightLbsLower: unit.Pound(1100),
			WeightLbsUpper: unit.Pound(1199),
			RateCents:      unit.Cents(1003),
		},
	})

	rate, err := models.FetchTariff400ngItemRate(suite.db, rate2.Code, 2, 1150, time.Date(2018, time.August, 15, 0, 0, 0, 0, time.UTC))

	// Ensure we get back rate2's rate and not one for a different weight range
	if suite.NoError(err) {
		suite.Equal(rate.RateCents, rate2.RateCents)
	}
}
