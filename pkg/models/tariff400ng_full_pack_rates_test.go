package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_EffectiveDateValidation() {
	now := time.Now()

	validPackRate := Tariff400ngFullPackRate{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validPackRate, expErrors)

	invalidPackRate := Tariff400ngFullPackRate{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"EffectiveDateLower must be after EffectiveDateUpper": {""},
	}
	suite.verifyValidationErrors(&invalidPackRate, expErrors)
}
