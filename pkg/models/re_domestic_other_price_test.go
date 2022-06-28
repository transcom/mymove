package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReDomesticOtherPriceValidations() {
	suite.Run("test valid ReDomesticOtherPrice", func() {
		validReDomesticOtherPrice := models.ReDomesticOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			IsPeakPeriod: true,
			Schedule:     2,
			PriceCents:   unit.Cents(431),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticOtherPrice, expErrors)
	})

	suite.Run("test empty ReDomesticOtherPrice", func() {
		emptyReDomesticOtherPrice := models.ReDomesticOtherPrice{}
		expErrors := map[string][]string{
			"contract_id": {"ContractID can not be blank."},
			"service_id":  {"ServiceID can not be blank."},
			"schedule":    {"0 is not greater than 0."},
			"price_cents": {"PriceCents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&emptyReDomesticOtherPrice, expErrors)
	})

	suite.Run("test ReDomesticOtherPrice with schedule about limit", func() {
		badScheduleReDomesticOtherPrice := models.ReDomesticOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			IsPeakPeriod: false,
			Schedule:     4,
			PriceCents:   unit.Cents(123),
		}
		expErrors := map[string][]string{
			"schedule": {"4 is not less than 4."},
		}
		suite.verifyValidationErrors(&badScheduleReDomesticOtherPrice, expErrors)
	})
}
