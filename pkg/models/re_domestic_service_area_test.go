package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReDomesticServiceAreaValidation() {
	suite.Run("test valid ReDomesticServiceArea", func() {
		validReDomesticServiceArea := models.ReDomesticServiceArea{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceArea:      "009",
			ServicesSchedule: 2,
			SITPDSchedule:    2,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticServiceArea, expErrors)
	})

	suite.Run("test invalid ReDomesticServiceArea", func() {
		emptyReDomesticServiceArea := models.ReDomesticServiceArea{}
		expErrors := map[string][]string{
			"contract_id":       {"ContractID can not be blank."},
			"service_area":      {"ServiceArea can not be blank."},
			"services_schedule": {"0 is not greater than 0."},
			"sitpdschedule":     {"0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&emptyReDomesticServiceArea, expErrors)
	})

	suite.Run("test schedules over 3 for ReDomesticServiceArea", func() {
		invalidReDomesticServiceArea := models.ReDomesticServiceArea{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceArea:      "009",
			ServicesSchedule: 4,
			SITPDSchedule:    5,
		}
		expErrors := map[string][]string{
			"services_schedule": {"4 is not less than 4."},
			"sitpdschedule":     {"5 is not less than 4."},
		}
		suite.verifyValidationErrors(&invalidReDomesticServiceArea, expErrors)
	})

	suite.Run("test schedules less than 1 for ReDomesticServiceArea", func() {
		invalidReDomesticServiceArea := models.ReDomesticServiceArea{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceArea:      "009",
			ServicesSchedule: -3,
			SITPDSchedule:    -1,
		}
		expErrors := map[string][]string{
			"services_schedule": {"-3 is not greater than 0."},
			"sitpdschedule":     {"-1 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReDomesticServiceArea, expErrors)
	})
}
