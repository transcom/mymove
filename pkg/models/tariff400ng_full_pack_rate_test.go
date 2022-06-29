package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestTariff400ngFullPackRateValidations() {
	now := time.Now()
	suite.Run("test valid Tariff400ngFullPackRate", func() {
		validTariff400ngFullPackRate := Tariff400ngFullPackRate{
			Schedule:           1,
			WeightLbsLower:     100,
			WeightLbsUpper:     200,
			RateCents:          100,
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(1, 0, 0),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validTariff400ngFullPackRate, expErrors)
	})

	suite.Run("test empty Tariff400ngFullPackRate", func() {
		emptyTariff400ngFullPackRate := Tariff400ngFullPackRate{}
		expErrors := map[string][]string{
			"schedule":             {"Schedule can not be blank."},
			"weight_lbs_lower":     {"0 is not less than 0."},
			"weight_lbs_upper":     {"WeightLbsUpper can not be blank."},
			"effective_date_lower": {"EffectiveDateLower can not be blank."},
			"effective_date_upper": {"EffectiveDateUpper can not be blank."},
		}
		suite.verifyValidationErrors(&emptyTariff400ngFullPackRate, expErrors)
	})
}

func (suite *ModelSuite) Test_EffectiveDateValidation() {
	now := time.Now()

	validPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		RateCents:          100,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		RateCents:          100,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": {"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}

func (suite *ModelSuite) Test_WeightValidation() {
	now := time.Now()

	validPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		RateCents:          100,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     200,
		WeightLbsUpper:     100,
		RateCents:          100,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors = map[string][]string{
		"weight_lbs_lower": {"200 is not less than 100."},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}

func (suite *ModelSuite) Test_RateValidation() {
	now := time.Now()

	validPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		RateCents:          100,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		Schedule:           1,
		RateCents:          -1,
		WeightLbsLower:     100,
		WeightLbsUpper:     200,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors = map[string][]string{
		"rate_cents": {"-1 is not greater than -1."},
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
		t.Errorf("Unable to query full pack rate: %v", err)
	}
	if rate != rateExpected {
		t.Errorf("Incorrect full pack rate received. Got: %d. Expected: %d.", rate, rateExpected)
	}

	// Test inclusivity of effective_date_lower
	_, err = FetchTariff400ngFullPackRateCents(suite.DB(), weight, schedule, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateUpper is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of effective_date_upper
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weight, schedule, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == rateExpected {
		t.Errorf("EffectiveDateUpper is incorrectly inclusive.")
	}

	// Test inclusivity of weight_lbs_lower
	_, err = FetchTariff400ngFullPackRateCents(suite.DB(), weightLower, schedule, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("WeightLbsLower is incorrectly exclusive: %s", err)
	}

	// Test exclusivity of weight_lbs_upper
	rate, err = FetchTariff400ngFullPackRateCents(suite.DB(), weightUpper, schedule, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == rateExpected {
		t.Errorf("WeightLbsUpper is incorrectly inclusive.")
	}
}
