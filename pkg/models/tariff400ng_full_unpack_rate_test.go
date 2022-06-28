package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_Tariff400ngFullUnpackRateValidation() {
	suite.Run("test valid Tariff400ngFullUnpackRate", func() {
		now := time.Now()
		validTariff400ngFullUnpackRate := Tariff400ngFullUnpackRate{
			Schedule:           1,
			RateMillicents:     100000,
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(1, 0, 0),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validTariff400ngFullUnpackRate, expErrors)
	})

	suite.Run("test invalid Tariff400ngFullUnpackRate", func() {
		invalidTariff400ngFullUnpackRate := Tariff400ngFullUnpackRate{}
		expErrors := map[string][]string{
			"schedule":             {"Schedule can not be blank."},
			"effective_date_lower": {"EffectiveDateLower can not be blank."},
			"effective_date_upper": {"EffectiveDateUpper can not be blank."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngFullUnpackRate, expErrors)
	})

	suite.Run("test negative RateMillicents, badly ordered dates for Tariff400ngFullUnpackRate", func() {
		now := time.Now()
		invalidTariff400ngFullUnpackRate := Tariff400ngFullUnpackRate{
			Schedule:           1,
			RateMillicents:     -200000,
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(-1, 0, 0),
		}
		expErrors := map[string][]string{
			"rate_millicents":      {"-200000 is not greater than -1."},
			"effective_date_upper": {"EffectiveDateUpper must be after EffectiveDateLower."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngFullUnpackRate, expErrors)
	})
}

func (suite *ModelSuite) Test_UnpackCreateAndSave() {
	now := time.Now()

	validUnpackRate := Tariff400ngFullUnpackRate{
		RateMillicents:     100,
		Schedule:           1,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(0, 1, 0),
	}

	suite.MustSave(&validUnpackRate)
}

func (suite *ModelSuite) Test_FetchFullUnPackRateCents() {
	t := suite.T()

	rateExpected := 100
	schedule := 1

	fupr := Tariff400ngFullUnpackRate{
		Schedule:           schedule,
		RateMillicents:     rateExpected,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fupr)

	rate, err := FetchTariff400ngFullUnpackRateMillicents(suite.DB(), schedule, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Errorf("Unable to query full unpack rate: %v", err)
	}
	if rate != rateExpected {
		t.Errorf("Incorrect full unpack rate received. Got: %d. Expected: %d.", rate, rateExpected)
	}

	// Test exclusivity of effective_date_upper
	rate, err = FetchTariff400ngFullUnpackRateMillicents(suite.DB(), schedule, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == rateExpected {
		t.Errorf("EffectiveDateUpper is inclusive of upper bound.")
	}

	// Test inclusivity of weight_lbs_upper
	_, err = FetchTariff400ngFullUnpackRateMillicents(suite.DB(), schedule, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateLower is exclusive of lower bound: %v", err)
	}
}
