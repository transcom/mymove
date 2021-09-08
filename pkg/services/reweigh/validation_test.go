package reweigh

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ReweighSuite) TestMergeReweigh() {
	id, _ := uuid.NewV4()
	shipmentID, _ := uuid.NewV4()
	requestedAt := time.Now()

	oldReweigh := models.Reweigh{
		ID:          id,
		ShipmentID:  shipmentID,
		RequestedAt: requestedAt,
		RequestedBy: models.ReweighRequesterTOO,
	}

	weight := int64(2000)
	verificationProvidedAt := time.Now()
	verificationReason := "Because I said so"

	newReweigh := models.Reweigh{
		Weight:                 handlers.PoundPtrFromInt64Ptr(&weight),
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}

	mergedReweigh := mergeReweigh(newReweigh, &oldReweigh)

	// check the new fields were populates
	suite.Equal(newReweigh.Weight, mergedReweigh.Weight)
	suite.Equal(newReweigh.VerificationReason, mergedReweigh.VerificationReason)
	suite.Equal(newReweigh.VerificationProvidedAt, mergedReweigh.VerificationProvidedAt)

	// Checking that the old reweigh instances weren't changed:
	suite.Equal(oldReweigh.ID, mergedReweigh.ID)
	suite.Equal(oldReweigh.ShipmentID, mergedReweigh.ShipmentID)
	suite.Equal(oldReweigh.RequestedAt, mergedReweigh.RequestedAt)
	suite.Equal(oldReweigh.RequestedBy, mergedReweigh.RequestedBy)
}
