package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
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
