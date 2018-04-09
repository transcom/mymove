package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_DiscountRateEffectiveDateValidation() {
	now := time.Now()

	validDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: "ABCD",
		EffectiveDateLower:       now,
		EffectiveDateUpper:       now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validDiscountRate, expErrors)

	invalidDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: "ABCD",
		EffectiveDateLower:       now,
		EffectiveDateUpper:       now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidDiscountRate, expErrors)
}

func (suite *ModelSuite) Test_DiscountRateSCACValidation() {
	now := time.Now()

	validDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: "ABCD",
		EffectiveDateLower:       now,
		EffectiveDateUpper:       now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validDiscountRate, expErrors)

	invalidDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: "",
		EffectiveDateLower:       now,
		EffectiveDateUpper:       now.AddDate(1, 0, 0),
	}

	expErrors = map[string][]string{
		"standard_carrier_alpha_code": []string{"StandardCarrierAlphaCode can not be blank."},
	}
	suite.verifyValidationErrors(&invalidDiscountRate, expErrors)
}
