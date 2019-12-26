package models_test

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentServiceItemValidation() {
	suite.T().Run("test valid PaymentServiceItem", func(t *testing.T) {
		validPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			ServiceItemID:    uuid.Must(uuid.NewV4()),
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
			PriceCents:       unit.Cents(1000),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentServiceItem, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
			"service_item_id":    {"ServiceItemID can not be blank."},
			"status":             {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID]."},
			"requested_at":       {"RequestedAt can not be blank."},
			// TODO: Removing this until we have pricing to populate
			// "price_cents":        {"PriceCents can not be blank.", "0 is not greater than 0."},
		}

		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})

	suite.T().Run("test invalid status for PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			ServiceItemID:    uuid.Must(uuid.NewV4()),
			Status:           "Sleeping",
			RequestedAt:      time.Now(),
			PriceCents:       unit.Cents(1000),
		}
		expErrors := map[string][]string{
			"status": {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID]."},
		}
		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})
}
