package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReTaskOrderFeeValidation() {
	suite.Run("test valid ReTaskOrderFee", func() {
		validReTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: uuid.Must(uuid.NewV4()),
			ServiceID:      uuid.Must(uuid.NewV4()),
			PriceCents:     9000,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReTaskOrderFee, expErrors)
	})

	suite.Run("test invalid ReTaskOrderFee", func() {
		invalidReTaskOrderFee := models.ReTaskOrderFee{}
		expErrors := map[string][]string{
			"contract_year_id": {"ContractYearID can not be blank."},
			"service_id":       {"ServiceID can not be blank."},
			"price_cents":      {"PriceCents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReTaskOrderFee, expErrors)
	})

	suite.Run("test price cents less than 1 for ReDomesticServiceArea", func() {
		invalidReTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: uuid.Must(uuid.NewV4()),
			ServiceID:      uuid.Must(uuid.NewV4()),
			PriceCents:     -3,
		}
		expErrors := map[string][]string{
			"price_cents": {"-3 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReTaskOrderFee, expErrors)
	})
}
