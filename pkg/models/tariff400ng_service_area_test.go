package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ServiceAreaEffectiveDateValidation() {
	now := time.Now()

	validServiceArea := Tariff400ngServiceArea{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaServiceChargeValidation() {
	validServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: 100,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: -1,
	}

	expErrors = map[string][]string{
		"service_charge_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaCreateAndSave() {
	now := time.Now()

	validServiceArea := Tariff400ngServiceArea{
		Name:               "Test",
		ServiceChargeCents: 100,
		ServiceArea:        1,
		LinehaulFactor:     1,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(0, 1, 0),
	}

	suite.mustSave(&validServiceArea)
}
