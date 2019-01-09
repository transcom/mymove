package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_LinehaulEffectiveDateValidation() {
	now := time.Now()

	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulWeightValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     200,
		WeightLbsUpper:     100,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors = map[string][]string{
		"weight_lbs_lower": []string{"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulRateValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		RateCents:          100,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		RateCents:          -1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors = map[string][]string{
		"rate_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
}

func (suite *ModelSuite) Test_LinehaulDistanceValidation() {
	validLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 100,
		DistanceMilesUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validLinehaulRate, expErrors)

	invalidLinehaulRate := Tariff400ngLinehaulRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		DistanceMilesLower: 200,
		DistanceMilesUpper: 100,
	}

	expErrors = map[string][]string{
		"distance_miles_lower": []string{"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidLinehaulRate, expErrors)
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
		t.Fatalf("Something went wrong with saving the test object: %s\n", err)
	}
	if rate != mySpecificRate {
		t.Errorf("The record object didn't save!")
	}

	// Test inclusivity of EffectiveDateLower
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, goodWeight, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of EffectiveDateUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, goodWeight, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == mySpecificRate {
		t.Errorf("EffectiveDateUpper is incorrectly inclusive.")
	}

	// Test inclusivity of DistanceMilesLower
	rate, err = FetchBaseLinehaulRate(suite.DB(), distanceLower, goodWeight, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("DistanceMilesLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of DistanceMilesUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), distanceUpper, goodWeight, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == mySpecificRate {
		t.Errorf("DistanceMilesUpper is incorrectly inclusive.")
	}

	// Test inclusivity of WeightLbsLower
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, weightLbsLower, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("WeightLbsLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of WeightLbsUpper
	rate, err = FetchBaseLinehaulRate(suite.DB(), goodDistance, weightLbsUpper, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == mySpecificRate {
		t.Errorf("DistanceMilesUpper is incorrectly inclusive.")
	}

}
