package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
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

	suite.mustSave(&validShorthaulRate)
}
