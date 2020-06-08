package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestProofOfServiceDocValidation() {
	suite.T().Run("test valid ProofOfServiceDoc", func(t *testing.T) {
		validProofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validProofOfServiceDoc, expErrors)
	})

	suite.T().Run("test empty ProofOfServiceDoc", func(t *testing.T) {
		invalidProofOfServiceDoc := models.ProofOfServiceDoc{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidProofOfServiceDoc, expErrors)
	})
}
