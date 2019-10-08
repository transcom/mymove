package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReShipmentTypeValidation() {
	suite.T().Run("test valid ReShipmentType", func(t *testing.T) {
		validReShipmentType := models.ReShipmentType{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReShipmentType, expErrors)
	})

	suite.T().Run("test empty ReShipmentType", func(t *testing.T) {
		invalidReShipmentType := models.ReShipmentType{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReShipmentType, expErrors)
	})
}
