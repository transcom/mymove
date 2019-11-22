package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentRequestValidation() {
	suite.T().Run("test valid PaymentRequest", func(t *testing.T) {
		validPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentRequest, expErrors)
	})

	suite.T().Run("test empty PaymentRequest", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{}

		expErrors := map[string][]string{
			"move_task_order_id": {"MoveTaskOrderID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})
}