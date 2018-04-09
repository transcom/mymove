package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_DiscountRateEffectiveDateValidation() {
	now := time.Now()

	validDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: testdatagen.RandomSCAC(),
		EffectiveDateLower:       now,
		EffectiveDateUpper:       now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validDiscountRate, expErrors)

	invalidDiscountRate := DiscountRate{
		StandardCarrierAlphaCode: testdatagen.RandomSCAC(),
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
		StandardCarrierAlphaCode: testdatagen.RandomSCAC(),
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
