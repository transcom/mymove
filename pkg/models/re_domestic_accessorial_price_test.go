package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReDomesticAccessorialPriceValidation() {
	suite.T().Run("test valid ReDomesticAccessorialPrice", func(t *testing.T) {
		validReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
			ContractID:      uuid.Must(uuid.NewV4()),
			ServiceID:       uuid.Must(uuid.NewV4()),
			ServiceSchedule: 2,
			PerUnitCents:    100,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticAccessorialPrice, expErrors)
	})

	suite.T().Run("test invalid ReDomesticAccessorialPrice", func(t *testing.T) {
		invalidReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{}
		expErrors := map[string][]string{
			"contract_id":      {"ContractID can not be blank."},
			"service_id":       {"ServiceID can not be blank."},
			"service_schedule": {"0 is not greater than 0."},
			"per_unit_cents":   {"PerUnitCents can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReDomesticAccessorialPrice, expErrors)
	})

	suite.T().Run("test service schedule over 3 for ReDomesticAccessorialPrice", func(t *testing.T) {
		invalidReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
			ContractID:      uuid.Must(uuid.NewV4()),
			ServiceID:       uuid.Must(uuid.NewV4()),
			ServiceSchedule: 4,
			PerUnitCents:    100,
		}
		expErrors := map[string][]string{
			"service_schedule": {"4 is not less than 4."},
		}
		suite.verifyValidationErrors(&invalidReDomesticAccessorialPrice, expErrors)
	})
}