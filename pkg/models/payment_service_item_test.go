package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentServiceItemValidation() {
	suite.T().Run("test valid PaymentServiceItem", func(t *testing.T) {
		validPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			Status:           "REQUESTED",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentServiceItem, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
			"status":             {"Status can not be blank.", "Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID]."},
		}

		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})

	suite.T().Run("test invalid status for PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			Status:           "Sleeping",
		}
		expErrors := map[string][]string{
			"status": {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID]."},
		}
		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})
}