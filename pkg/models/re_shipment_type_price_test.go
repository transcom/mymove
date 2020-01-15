package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReShipmentTypePriceValidation() {
	suite.T().Run("test valid ReShipmentTypePrice", func(t *testing.T) {
		validReShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "C",
			Factor:     1.20,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReShipmentTypePrice, expErrors)
	})

	suite.T().Run("test invalid ReShipmentTypePrice", func(t *testing.T) {
		invalidReShipmentTypePrice := models.ReShipmentTypePrice{}
		expErrors := map[string][]string{
			"contract_id": {"ContractID can not be blank."},
			"service_id":  {"ServiceID can not be blank."},
			"market":      {"Market can not be blank.", "Market is not in the list [C, O]."},
			"factor":      {"0.000000 is not greater than 0.000000.", "Factor can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReShipmentTypePrice, expErrors)
	})

	suite.T().Run("test invalid market for ReShipmentTypePrice", func(t *testing.T) {
		invalidShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "R",
			Factor:     1.20,
		}
		expErrors := map[string][]string{
			"market": {"Market is not in the list [C, O]."},
		}
		suite.verifyValidationErrors(&invalidShipmentTypePrice, expErrors)
	})

	suite.T().Run("test factor hundredths less than 1 for ReShipmentTypePrice", func(t *testing.T) {
		invalidShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "C",
			Factor:     -3,
		}
		expErrors := map[string][]string{
			"factor": {"-3.000000 is not greater than 0.000000."},
		}
		suite.verifyValidationErrors(&invalidShipmentTypePrice, expErrors)
	})
}
