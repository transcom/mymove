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
			Status:          "PENDING",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentRequest, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItem", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
		}

		expErrors := map[string][]string{
			"status": {"Status is not in the list [PENDING, REVIEWED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID]."},
		}

		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})

	suite.T().Run("test invalid status for PaymentRequest", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			Status:          "Sleeping",
		}
		expErrors := map[string][]string{
			"status": {"Status is not in the list [PENDING, REVIEWED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID]."},
		}
		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})
}
