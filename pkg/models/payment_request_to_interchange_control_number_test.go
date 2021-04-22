package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestPaymentRequestToInterchangeControlNumber() {
	suite.T().Run("test PaymentRequest association", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		validPR2ICN := testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{PaymentRequestID: paymentRequest.ID}})
		suite.Equal(paymentRequest.ID, validPR2ICN.PaymentRequestID)
		err := suite.DB().Load(&validPR2ICN, "PaymentRequest")
		suite.NoError(err)
		suite.Equal(paymentRequest.ID, validPR2ICN.PaymentRequest.ID)
	})

	suite.T().Run("test valid PaymentRequestToInterchangeControlNumber", func(t *testing.T) {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Must(uuid.NewV4()),
			InterchangeControlNumber: 1,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})

	suite.T().Run("test invalid PaymentRequestToInterchangeControlNumber", func(t *testing.T) {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Nil,
			InterchangeControlNumber: 0,
		}
		expErrors := map[string][]string{
			"payment_request_id":         {"PaymentRequestID can not be blank."},
			"interchange_control_number": {"0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})

	suite.T().Run("test invalid InterchangeControlNumber max", func(t *testing.T) {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Must(uuid.NewV4()),
			InterchangeControlNumber: 1000000000,
		}
		expErrors := map[string][]string{
			"interchange_control_number": {"1000000000 is not less than 1000000000."},
		}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})
}
