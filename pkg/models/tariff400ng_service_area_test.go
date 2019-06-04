package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ServiceAreaEffectiveDateValidation() {
	now := time.Now()

	validServiceArea := Tariff400ngServiceArea{
		ServiceArea:        testdatagen.DefaultServiceArea,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		ServiceArea:        testdatagen.DefaultServiceArea,
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors = map[string][]string{
		"effective_date_upper": {"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaServiceChargeValidation() {
	validServiceArea := Tariff400ngServiceArea{
		ServiceArea:        testdatagen.DefaultServiceArea,
		ServiceChargeCents: 100,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		ServiceArea:        testdatagen.DefaultServiceArea,
		ServiceChargeCents: -1,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors = map[string][]string{
		"service_charge_cents": {"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaSITRatesValidation() {
	invalidServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: 1,
	}

	expErrors := map[string][]string{
		"service_area": {"ServiceArea can not be blank.",
			"ServiceArea does not match the expected format."},
		"s_i_t185_b_rate_cents": {"SIT185BRateCents can not be blank."},
		"s_i_t_p_d_schedule":    {"SITPDSchedule can not be blank."},
		"s_i_t185_a_rate_cents": {"SIT185ARateCents can not be blank."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}
