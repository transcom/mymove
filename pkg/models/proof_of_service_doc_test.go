package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestProofOfServiceDocValidation() {
	suite.Run("test valid ProofOfServiceDoc", func() {
		validProofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validProofOfServiceDoc, expErrors)
	})

	suite.Run("test empty ProofOfServiceDoc", func() {
		invalidProofOfServiceDoc := models.ProofOfServiceDoc{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidProofOfServiceDoc, expErrors)
	})
}
