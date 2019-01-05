package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_UnpackEffectiveDateValidation() {
	now := time.Now()

	validUnpackRate := Tariff400ngFullUnpackRate{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validUnpackRate, expErrors)

	invalidUnpackRate := Tariff400ngFullUnpackRate{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidUnpackRate, expErrors)
}

func (suite *ModelSuite) Test_UnpackRateValidation() {
	validUnpackRate := Tariff400ngFullUnpackRate{
		RateMillicents: 100,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validUnpackRate, expErrors)

	invalidUnpackRate := Tariff400ngFullUnpackRate{
		RateMillicents: -1,
	}

	expErrors = map[string][]string{
		"rate_millicents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidUnpackRate, expErrors)
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
		t.Fatalf("Unable to query full unpack rate: %v", err)
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
	rate, err = FetchTariff400ngFullUnpackRateMillicents(suite.DB(), schedule, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateLower is exclusive of lower bound: %v", err)
	}
}
