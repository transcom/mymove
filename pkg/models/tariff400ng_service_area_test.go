package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_Tariff400ngServiceAreaValidation() {
	suite.Run("test valid Tariff400ngServiceArea", func() {
		now := time.Now()
		validTariff400ngServiceArea := Tariff400ngServiceArea{
			Name:               "Birmingham, AL",
			ServiceArea:        testdatagen.DefaultServiceArea,
			ServicesSchedule:   1,
			LinehaulFactor:     unit.Cents(100),
			ServiceChargeCents: unit.Cents(100),
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(1, 0, 0),
			SIT185ARateCents:   unit.Cents(100),
			SIT185BRateCents:   unit.Cents(100),
			SITPDSchedule:      1,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validTariff400ngServiceArea, expErrors)
	})

	suite.Run("test invalid Tariff400ngServiceArea", func() {
		invalidTariff400ngServiceArea := Tariff400ngServiceArea{}
		expErrors := map[string][]string{
			"name":                 {"Name can not be blank."},
			"service_area":         {"ServiceArea can not be blank.", "ServiceArea does not match the expected format."},
			"services_schedule":    {"ServicesSchedule can not be blank."},
			"effective_date_lower": {"EffectiveDateLower can not be blank."},
			"effective_date_upper": {"EffectiveDateUpper can not be blank."},
			"sit185_arate_cents":   {"SIT185ARateCents can not be blank."},
			"sit185_brate_cents":   {"SIT185BRateCents can not be blank."},
			"sitpdschedule":        {"SITPDSchedule can not be blank."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngServiceArea, expErrors)
	})

	suite.Run("test other validations not exercised above for Tariff400ngServiceArea", func() {
		now := time.Now()
		invalidTariff400ngServiceArea := Tariff400ngServiceArea{
			Name:               "Birmingham, AL",
			ServiceArea:        testdatagen.DefaultServiceArea,
			ServicesSchedule:   1,
			LinehaulFactor:     unit.Cents(-100),
			ServiceChargeCents: unit.Cents(-100),
			EffectiveDateLower: now,
			EffectiveDateUpper: now.AddDate(-1, 0, 0),
			SIT185ARateCents:   unit.Cents(100),
			SIT185BRateCents:   unit.Cents(100),
			SITPDSchedule:      1,
		}
		expErrors := map[string][]string{
			"linehaul_factor":      {"-100 is not greater than -1."},
			"service_charge_cents": {"-100 is not greater than -1."},
			"effective_date_upper": {"EffectiveDateUpper must be after EffectiveDateLower."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngServiceArea, expErrors)
	})
}
