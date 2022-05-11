package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReDomesticAccessorialPriceValidation() {
	suite.T().Run("test valid ReDomesticAccessorialPrice", func(t *testing.T) {
		validReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceID:        uuid.Must(uuid.NewV4()),
			ServicesSchedule: 2,
			PerUnitCents:     unit.Cents(99),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticAccessorialPrice, expErrors)
	})

	suite.T().Run("test invalid ReDomesticAccessorialPrice", func(t *testing.T) {
		invalidReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{}
		expErrors := map[string][]string{
			"contract_id":       {"ContractID can not be blank."},
			"service_id":        {"ServiceID can not be blank."},
			"services_schedule": {"0 is not greater than 0."},
			"per_unit_cents":    {"PerUnitCents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReDomesticAccessorialPrice, expErrors)
	})

	suite.T().Run("test service schedule over 3 for ReDomesticAccessorialPrice", func(t *testing.T) {
		invalidReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceID:        uuid.Must(uuid.NewV4()),
			ServicesSchedule: 4,
			PerUnitCents:     unit.Cents(99),
		}
		expErrors := map[string][]string{
			"services_schedule": {"4 is not less than 4."},
		}
		suite.verifyValidationErrors(&invalidReDomesticAccessorialPrice, expErrors)
	})

	suite.T().Run("test per unit cents is not negative ReDomesticAccessorialPrice", func(t *testing.T) {
		invalidReDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
			ContractID:       uuid.Must(uuid.NewV4()),
			ServiceID:        uuid.Must(uuid.NewV4()),
			ServicesSchedule: 2,
			PerUnitCents:     unit.Cents(-10),
		}
		expErrors := map[string][]string{
			"per_unit_cents": {"-10 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReDomesticAccessorialPrice, expErrors)
	})
}
