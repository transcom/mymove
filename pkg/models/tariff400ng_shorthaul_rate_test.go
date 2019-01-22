package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ShorthaulRateEffectiveDateValidation() {
	now := time.Now()

	validShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower:      10,
		CwtMilesUpper:      20,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validShorthaulRate, expErrors)

	invalidShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower:      10,
		CwtMilesUpper:      20,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidShorthaulRate, expErrors)
}

func (suite *ModelSuite) Test_ShorthaulRateServiceChargeValidation() {
	validShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower: 10,
		CwtMilesUpper: 20,
		RateCents:     100,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validShorthaulRate, expErrors)

	invalidShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower: 10,
		CwtMilesUpper: 20,
		RateCents:     -1,
	}

	expErrors = map[string][]string{
		"service_charge_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidShorthaulRate, expErrors)
}

func (suite *ModelSuite) Test_ShorthaulRateCwtMilesValidation() {
	validShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower: 10,
		CwtMilesUpper: 20,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validShorthaulRate, expErrors)

	invalidShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower: 20,
		CwtMilesUpper: 10,
	}

	expErrors = map[string][]string{
		"cwt_miles_upper": []string{"10 is not greater than 20."},
	}
	suite.verifyValidationErrors(&invalidShorthaulRate, expErrors)
}

func (suite *ModelSuite) Test_ShorthaulRateCreateAndSave() {
	now := time.Now()

	validShorthaulRate := Tariff400ngShorthaulRate{
		CwtMilesLower:      10,
		CwtMilesUpper:      20,
		RateCents:          100,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(0, 1, 0),
	}

	suite.MustSave(&validShorthaulRate)
}

func (suite *ModelSuite) Test_FetchShorthaulRateCents() {
	t := suite.T()

	rate1 := unit.Cents(100)
	sh1 := Tariff400ngShorthaulRate{
		CwtMilesLower:      1000,
		CwtMilesUpper:      2000,
		RateCents:          rate1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&sh1)

	// Test inclusivity CwtMilesLower
	rate, err := FetchShorthaulRateCents(suite.DB(), 1000, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("Unable to query shorthaul rate: %s", err)
	}
	if rate != rate1 {
		t.Errorf("Incorrect shorthaul rate. Got: %d, expected %d", rate, rate1)
	}

	// Test exclusivity CwtMilesUpper
	rate, err = FetchShorthaulRateCents(suite.DB(), 2000, testdatagen.DateInsidePeakRateCycle)
	if err == nil && rate == rate1 {
		t.Errorf("CwtMilesLower is incorrectly inlusive")
	}

	// Test inclusivity of EffectiveDateLower
	rate, err = FetchShorthaulRateCents(suite.DB(), 1500, testdatagen.PeakRateCycleStart)
	if err != nil {
		t.Errorf("EffectiveDateLower is incorrectly exlusive: %s", err)
	}

	// Test exclusivity of EffectiveDateUpper
	rate, err = FetchShorthaulRateCents(suite.DB(), 1500, testdatagen.PeakRateCycleEnd)
	if err == nil && rate == rate1 {
		t.Errorf("EffectiveDateUpper is incorrectly inclusive.")
	}

	rate2 := unit.Cents(200)
	sh2 := Tariff400ngShorthaulRate{
		CwtMilesLower:      2000,
		CwtMilesUpper:      3000,
		RateCents:          rate2,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&sh2)

	// Test the upper bound of the CwtMiles
	rate, err = FetchShorthaulRateCents(suite.DB(), 2999, testdatagen.DateInsidePeakRateCycle)
	if err != nil {
		t.Fatalf("Unable to query shorthaul rate: %s", err)
	}
	if rate != rate2 {
		t.Errorf("Incorrect shorthaul rate. Got: %d, expected %d", rate, rate2)
	}

	// And this rate is in a different EffectiveDate band
	rate3 := unit.Cents(300)
	sh3 := Tariff400ngShorthaulRate{
		CwtMilesLower:      1000,
		CwtMilesUpper:      2000,
		RateCents:          rate3,
		EffectiveDateLower: testdatagen.NonPeakRateCycleStart,
		EffectiveDateUpper: testdatagen.NonPeakRateCycleEnd,
	}
	suite.MustSave(&sh3)

	rate, err = FetchShorthaulRateCents(suite.DB(), 1000, testdatagen.DateOutsidePeakRateCycle)
	if err != nil {
		t.Fatalf("Unable to query shorthaul rate: %s", err)
	}
	if rate != rate3 {
		t.Errorf("Incorrect shorthaul rate. Got: %d, expected %d", rate, rate3)
	}

}
