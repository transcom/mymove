package models_test

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"testing"
)

func (suite *ModelSuite) TestPaymentRequestValidation() {
	suite.T().Run("test valid PaymentRequest", func(t *testing.T) {
		validPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID:       uuid.Must(uuid.NewV4()),
			ServiceItemIDs:        []uuid.UUID{uuid.Must(uuid.NewV4())},
			RejectionReason:		"Not enough documentation",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentRequest, expErrors)
	})

	suite.T().Run("test empty PaymentRequest", func(t *testing.T){
		invalidPaymentRequest := models.PaymentRequest{}

		expErrors := map[string][]string{
			"move_task_order_id":       {"MoveTaskOrderID can not be blank."},
			"service_item_id_s":        {"ServiceItemIDs can not be empty."},
			"rejection_reason":			{"RejectionReason can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})
}