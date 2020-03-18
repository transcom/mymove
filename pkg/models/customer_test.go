package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCustomerValidation() {
	suite.T().Run("test valid Customer", func(t *testing.T) {
		validCustomer := models.Customer{
			UserID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validCustomer, expErrors)
	})
}
