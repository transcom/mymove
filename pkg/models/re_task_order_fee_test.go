package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReTaskOrderFeeValidation() {
	suite.T().Run("test valid ReTaskOrderFee", func(t *testing.T) {
		validReTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: uuid.Must(uuid.NewV4()),
			ServiceID:      uuid.Must(uuid.NewV4()),
			PriceCents:     9000,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReTaskOrderFee, expErrors)
	})

	suite.T().Run("test invalid ReTaskOrderFee", func(t *testing.T) {
		invalidReTaskOrderFee := models.ReTaskOrderFee{}
		expErrors := map[string][]string{
			"contract_year_id": {"ContractYearID can not be blank."},
			"service_id":       {"ServiceID can not be blank."},
			"price_cents":      {"PriceCents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReTaskOrderFee, expErrors)
	})
}