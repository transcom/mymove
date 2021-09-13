package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestCreateApprovedSITExtension() {
	suite.T().Run("Returns an error when shipment is not found", func(t *testing.T) {
		sitExtensionCreator := NewApprovedSITExtensionCreator()
		nonexistentUUID := uuid.Must(uuid.NewV4())
		requestedDays := 45
		officeRemarks := "office remarks"
		sitExtensionToSave := models.SITExtension{
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			RequestedDays: requestedDays,
			ApprovedDays:  &requestedDays,
			OfficeRemarks: &officeRemarks,
			Status:        models.SITExtensionStatusApproved,
		}
		eTag := ""

		_, err := sitExtensionCreator.CreateApprovedSITExtension(suite.TestAppContext(), &sitExtensionToSave, nonexistentUUID, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when etag does not match", func(t *testing.T) {
		sitExtensionCreator := NewApprovedSITExtensionCreator()
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		requestedDays := 45
		officeRemarks := "office remarks"
		sitExtensionToSave := models.SITExtension{
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			RequestedDays: requestedDays,
			ApprovedDays:  &requestedDays,
			OfficeRemarks: &officeRemarks,
			Status:        models.SITExtensionStatusApproved,
		}
		eTag := ""

		_, err := sitExtensionCreator.CreateApprovedSITExtension(suite.TestAppContext(), &sitExtensionToSave, mtoShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})

	suite.T().Run("Creates one approved SIT extension when all fields are valid and updates the shipment's SIT days allowance", func(t *testing.T) {
		sitExtensionCreator := NewApprovedSITExtensionCreator()
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		requestedDays := 45
		officeRemarks := "office remarks"
		sitExtensionToSave := models.SITExtension{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			RequestedDays: requestedDays,
			ApprovedDays:  &requestedDays,
			OfficeRemarks: &officeRemarks,
			Status:        models.SITExtensionStatusApproved,
		}

		updatedShipment, err := sitExtensionCreator.CreateApprovedSITExtension(suite.TestAppContext(), &sitExtensionToSave, mtoShipment.ID, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITExtension
		err = suite.DB().First(&sitExtensionInDB)
		suite.NoError(err)

		var allSITExtensions []models.SITExtension
		err = suite.DB().All(&allSITExtensions)
		suite.NoError(err)
		suite.Equal(1, len(allSITExtensions))

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(requestedDays, *updatedShipment.SITDaysAllowance)
		suite.Equal(requestedDays, *sitExtensionInDB.ApprovedDays)
		suite.Equal(requestedDays, sitExtensionInDB.RequestedDays)
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(models.SITExtensionStatusApproved, sitExtensionInDB.Status)
	})
}
