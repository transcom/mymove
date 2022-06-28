package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentRequestToInterchangeControlNumber() {
	suite.Run("test valid PaymentRequestToInterchangeControlNumber", func() {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Must(uuid.NewV4()),
			InterchangeControlNumber: 1,
			EDIType:                  models.EDIType997,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})

	suite.Run("test invalid PaymentRequestToInterchangeControlNumber", func() {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Nil,
			InterchangeControlNumber: 0,
			EDIType:                  "models.EDIType997",
		}
		expErrors := map[string][]string{
			"payment_request_id":         {"PaymentRequestID can not be blank."},
			"interchange_control_number": {"0 is not greater than 0."},
			"editype":                    {"EDIType is not in the list [810, 824, 858, 997]."},
		}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})

	suite.Run("test invalid InterchangeControlNumber max", func() {
		validPR2ICN := models.PaymentRequestToInterchangeControlNumber{
			PaymentRequestID:         uuid.Must(uuid.NewV4()),
			InterchangeControlNumber: 1000000000,
			EDIType:                  models.EDIType997,
		}
		expErrors := map[string][]string{
			"interchange_control_number": {"1000000000 is not less than 1000000000."},
		}
		suite.verifyValidationErrors(&validPR2ICN, expErrors)
	})
}
