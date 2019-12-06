package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentRequestValidation() {
	suite.T().Run("test valid PaymentRequest", func(t *testing.T) {
		validPaymentRequest := models.PaymentRequest{
			Status: "PENDING",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentRequest, expErrors)
	})

	suite.T().Run("test invalid status for PaymentRequest", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			Status: "Sleeping",
		}
		expErrors := map[string][]string{
			"status": {"Status is not in the list [PENDING, REVIEWED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID]."},
		}
		suite.verifyValidationErrors(&invalidPaymentRequest, expErrors)
	})
}