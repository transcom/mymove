package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_EffectiveDateValidation() {
	now := time.Now()

	validPackRate := Tariff400ngFullPackRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}

func (suite *ModelSuite) Test_WeightValidation() {
	validPackRate := Tariff400ngFullPackRate{
		WeightLbsLower: 100,
		WeightLbsUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		WeightLbsLower: 200,
		WeightLbsUpper: 100,
	}

	expErrors = map[string][]string{
		"weight_lbs_lower": []string{"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}

func (suite *ModelSuite) Test_RateValidation() {
	validPackRate := Tariff400ngFullPackRate{
		RateCents:      100,
		WeightLbsLower: 100,
		WeightLbsUpper: 200,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		RateCents:      -1,
		WeightLbsLower: 100,
		WeightLbsUpper: 200,
	}

	expErrors = map[string][]string{
		"rate_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}

func (suite *ModelSuite) Test_FetchFullPackRateCents() {
	t := suite.T()

	rateExpected := unit.Cents(100)
	weight := unit.Pound(1500)
	weightLower := unit.Pound(1000)
	weightUpper := unit.Pound(2000)
	schedule := 1

	fpr := Tariff400ngFullPackRate{
		Schedule:           schedule,
		WeightLbsLower:     weightLower,
		WeightLbsUpper:     weightUpper,
		RateCents:          rateExpected,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fpr)

	rate, err := FetchTariff400ngFullPackRateCents(suite.DB(), weight, schedule, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("Unable to query full pack rate: %v", err)
	}
	if rate != rate {
		t.Errorf("Incorrect full pack rate received. Got: %d. Expected: %d.", rate, rateExpected)
	}

	// Test inclusivity of effective_date_lower
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weight, schedule, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateUpper is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of effective_date_upper
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weight, schedule, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == rateExpected {
		t.Errorf("EffectiveDateUpper is incorrectly inclusive.")
	}

	// Test inclusivity of weight_lbs_lower
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weightLower, schedule, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("WeightLbsLower is incorrectly exclusive: %s", err)
	}

	// Test exclusivity of weight_lbs_upper
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weightUpper, schedule, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == rateExpected {
		t.Errorf("WeightLbsUpper is incorrectly inclusive.")
	}
}
