package mtoshipment

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApproveSITExtension() {
	suite.T().Run("Returns an error when shipment is not found", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		nonexistentUUID := uuid.Must(uuid.NewV4())
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), nonexistentUUID, nonexistentUUID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when SIT extension is not found", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), mtoShipment.ID, nonexistentUUID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when etag does not match", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})

	suite.T().Run("Returns an error when shipment ID from SIT extension and shipment ID found do not match", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		otherMtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), otherMtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.Contains(err.Error(), "SITExtension's shipment ID does not match shipment ID provided")
	})

	suite.T().Run("Updates the shipment's SIT days allowance and the SIT extension's status and approved days if all fields are valid", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: swag.Int(20),
			},
		})
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		// existing SITDaysAllowance plus new approved days
		newSITDaysAllowance := int(40)
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		updatedShipment, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITExtension
		err = suite.DB().Find(&sitExtensionInDB, sitExtension.ID)
		suite.NoError(err)

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(newSITDaysAllowance, *updatedShipment.SITDaysAllowance)
		suite.Equal(approvedDays, *sitExtensionInDB.ApprovedDays)
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(models.SITExtensionStatusApproved, sitExtensionInDB.Status)
	})
}
