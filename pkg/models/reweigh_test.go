package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReweighValidation() {
	suite.Run("test valid Reweigh", func() {
		validReweigh := models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterCustomer,
			ShipmentID:  uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReweigh, expErrors)
	})

	suite.Run("test empty reweigh", func() {
		expErrors := map[string][]string{
			"requested_at": {"RequestedAt can not be blank."},
			"requested_by": {"RequestedBy is not in the list [CUSTOMER, PRIME, SYSTEM, TOO]."},
			"shipment_id":  {"ShipmentID can not be blank."},
		}
		suite.verifyValidationErrors(&models.Reweigh{}, expErrors)
	})

	suite.Run("test validation failures", func() {
		var verificationReason string
		weight := unit.Pound(-1)
		invalidReweigh := models.Reweigh{
			RequestedAt:        time.Now(),
			RequestedBy:        models.ReweighRequesterCustomer,
			ShipmentID:         uuid.Must(uuid.NewV4()),
			VerificationReason: &verificationReason,
			Weight:             &weight,
		}
		expErrors := map[string][]string{
			"weight":              {"-1 is less than or equal to zero"},
			"verification_reason": {"VerificationReason can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReweigh, expErrors)
	})
}
