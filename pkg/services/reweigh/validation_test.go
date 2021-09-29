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

func (suite *ReweighSuite) TestReweighChanged() {
	id, _ := uuid.NewV4()
	shipmentID, _ := uuid.NewV4()
	requestedAt := time.Now()

	reweigh1 := models.Reweigh{
		ID:          id,
		ShipmentID:  shipmentID,
		RequestedAt: requestedAt,
		RequestedBy: models.ReweighRequesterTOO,
	}

	weight1 := int64(2000)
	verificationProvidedAt := time.Now()
	verificationReason := "Because I said so"

	reweigh2 := models.Reweigh{
		Weight:                 handlers.PoundPtrFromInt64Ptr(&weight1),
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}

	// weights went from nil to some value
	weightChanged := reweighChanged(reweigh2, reweigh1)
	suite.Equal(true, weightChanged)
	suite.NotEqual(reweigh2.Weight, reweigh1.Weight)

	weight2 := int64(2000)
	reweigh3 := models.Reweigh{
		Weight:                 handlers.PoundPtrFromInt64Ptr(&weight2),
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}

	// weights are the same
	weightChanged = reweighChanged(reweigh3, reweigh2)
	suite.Equal(false, weightChanged)
	suite.Equal(*reweigh3.Weight, *reweigh2.Weight)

	// weights go from some value to nil
	reweigh4 := models.Reweigh{
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}
	weightChanged = reweighChanged(reweigh4, reweigh3)
	suite.Equal(true, weightChanged)
	suite.NotEqual(reweigh4.Weight, reweigh3.Weight)

	// weights both nil
	reweigh5 := models.Reweigh{
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}
	weightChanged = reweighChanged(reweigh5, reweigh4)
	suite.Equal(false, weightChanged)
	suite.Equal(reweigh5.Weight, reweigh4.Weight)

	// weights are different
	weight3 := int64(5000)
	reweigh6 := models.Reweigh{
		Weight:                 handlers.PoundPtrFromInt64Ptr(&weight3),
		VerificationProvidedAt: &verificationProvidedAt,
		VerificationReason:     &verificationReason,
	}
	weightChanged = reweighChanged(reweigh6, reweigh3)
	suite.Equal(true, weightChanged)
	suite.NotEqual(*reweigh6.Weight, *reweigh3.Weight)
}
