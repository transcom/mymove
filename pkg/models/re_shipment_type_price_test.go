package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReShipmentTypePriceValidation() {
	suite.T().Run("test valid ReShipmentTypePrice", func(t *testing.T) {
		validReShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID:       uuid.Must(uuid.NewV4()),
			ShipmentTypeID:   uuid.Must(uuid.NewV4()),
			Market:           "C",
			FactorHundredths: 1,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReShipmentTypePrice, expErrors)
	})

	suite.T().Run("test invalid ReShipmentTypePrice", func(t *testing.T) {
		invalidReShipmentTypePrice := models.ReShipmentTypePrice{}
		expErrors := map[string][]string{
			"contract_id":       {"ContractID can not be blank."},
			"shipment_type_id":  {"ShipmentTypeID can not be blank."},
			"market":            {"Market can not be blank.", "Market is not in the list [C, O]."},
			"factor_hundredths": {"FactorHundredths can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReShipmentTypePrice, expErrors)
	})

	suite.T().Run("test invalid market for ReShipmentTypePrice", func(t *testing.T) {
		invalidShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID:       uuid.Must(uuid.NewV4()),
			ShipmentTypeID:   uuid.Must(uuid.NewV4()),
			Market:           "R",
			FactorHundredths: 1,
		}
		expErrors := map[string][]string{
			"market": {"Market is not in the list [C, O]."},
		}
		suite.verifyValidationErrors(&invalidShipmentTypePrice, expErrors)
	})
}