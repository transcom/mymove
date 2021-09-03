package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentRequestValidation() {
	suite.T().Run("test valid PaymentRequest", func(t *testing.T) {
		repricedPaymentRequestID := uuid.Must(uuid.NewV4())
		validPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID:          uuid.Must(uuid.NewV4()),
			Status:                   models.PaymentRequestStatusPending,
			PaymentRequestNumber:     "1111-2222-1",
			SequenceNumber:           1,
			RepricedPaymentRequestID: &repricedPaymentRequestID,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentRequest, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItem", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
		}

		expErrors := map[string][]string{
			"status":                 {"Status is not in the list [PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR]."},
			"payment_request_number": {"PaymentRequestNumber can not be blank."},
			"sequence_number":        {"0 is not greater than 0."},
		}

		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})

	suite.T().Run("test invalid fields for PaymentRequest", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID:          uuid.Must(uuid.NewV4()),
			Status:                   "Sleeping",
			PaymentRequestNumber:     "1111-2222-1",
			SequenceNumber:           1,
			RepricedPaymentRequestID: &uuid.Nil,
		}
		expErrors := map[string][]string{
			"status":                      {"Status is not in the list [PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR]."},
			"repriced_payment_request_id": {"RepricedPaymentRequestID can not be blank."},
		}
		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})
}
