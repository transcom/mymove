package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCustomerValidation() {
	suite.T().Run("test valid Customer", func(t *testing.T) {
		validCustomer := models.Customer{
			DODID:     "1234567890",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validCustomer, expErrors)
	})
}
