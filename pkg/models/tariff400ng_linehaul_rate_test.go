package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_Tariff400ngLinehaulRateValidation() {
	suite.Run("test valid Tariff400ngLinehaulRate", func() {
		now := time.Now()
		validTariff400ngLinehaulRate := Tariff400ngLinehaulRate{
			DistanceMilesLower: 100,
			DistanceMilesUpper: 200,
			Type:               "ConusLinehaul",
			WeightLbsLower:     unit.Pound(100),
			WeightLbsUpper:     unit.Pound(200),
			RateCents:          unit.Cents(100),
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(1, 0, 0),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validTariff400ngLinehaulRate, expErrors)
	})

	suite.Run("test invalid Tariff400ngLinehaulRate", func() {
		invalidTariff400ngLinehaulRate := Tariff400ngLinehaulRate{}
		expErrors := map[string][]string{
			"distance_miles_lower": {"DistanceMilesLower can not be blank.", "0 is not less than 0."},
			"distance_miles_upper": {"DistanceMilesUpper can not be blank."},
			"type":                 {"Type can not be blank."},
			"weight_lbs_lower":     {"WeightLbsLower can not be blank.", "0 is not less than 0."},
			"weight_lbs_upper":     {"WeightLbsUpper can not be blank."},
			"effective_date_lower": {"EffectiveDateLower can not be blank."},
			"effective_date_upper": {"EffectiveDateUpper can not be blank."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngLinehaulRate, expErrors)
	})

	suite.Run("test negative RateCents, badly ordered dates for Tariff400ngLinehaulRate", func() {
		now := time.Now()
		invalidTariff400ngLinehaulRate := Tariff400ngLinehaulRate{
			DistanceMilesLower: 100,
			DistanceMilesUpper: 200,
			Type:               "ConusLinehaul",
			WeightLbsLower:     unit.Pound(100),
			WeightLbsUpper:     unit.Pound(200),
			RateCents:          unit.Cents(-200),
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(-1, 0, 0),
		}
		expErrors := map[string][]string{
			"rate_cents":           {"-200 is not greater than -1."},
			"effective_date_upper": {"EffectiveDateUpper must be after EffectiveDateLower."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngLinehaulRate, expErrors)
	})
}

func (suite *ModelSuite) Test_FetchBaseLinehaulRate() {
	t := suite.T()

	mySpecificRate := unit.Cents(474747)
	distanceLower := 3101
	distanceUpper := 3300
	weightLbsLower := unit.Pound(5000)
	weightLbsUpper := unit.Pound(10000)

	newBaseLinehaul := Tariff400ngLinehaulRate{
		DistanceMilesLower: distanceLower,
		DistanceMilesUpper: distanceUpper,
		WeightLbsLower:     weightLbsLower,
		WeightLbsUpper:     weightLbsUpper,
		RateCents:          mySpecificRate,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&newBaseLinehaul)

	goodDistance := 3200
	goodWeight := unit.Pound(6000)

	// Test the best case
	rate, err := FetchBaseLinehaulRate(suite.DB(), goodDistance, goodWeight, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("Something went wrong with saving the test object: %s\n", err)
	}
	if rate != mySpecificRate {
		t.Errorf("The record object didn't save!")
	}

	// Test inclusivity of EffectiveDateLower
	_, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, goodWeight, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of EffectiveDateUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, goodWeight, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == mySpecificRate {
		t.Errorf("EffectiveDateUpper is incorrectly inclusive.")
	}

	// Test inclusivity of DistanceMilesLower
	_, err = FetchBaseLinehaulRate(suite.DB(), distanceLower, goodWeight, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("DistanceMilesLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of DistanceMilesUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), distanceUpper, goodWeight, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == mySpecificRate {
		t.Errorf("DistanceMilesUpper is incorrectly inclusive.")
	}

	// Test inclusivity of WeightLbsLower
	_, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, weightLbsLower, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("WeightLbsLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of WeightLbsUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, weightLbsUpper, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == mySpecificRate {
		t.Errorf("DistanceMilesUpper is incorrectly inclusive.")
	}

}
