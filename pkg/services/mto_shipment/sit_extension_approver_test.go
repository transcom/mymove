package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApproveSITExtension() {
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionApprover := NewSITExtensionApprover(moveRouter)

	suite.Run("Returns an error when shipment is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), nonexistentUUID, nonexistentUUID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when SIT extension is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), mtoShipment.ID, nonexistentUUID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when etag does not match", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		sitExtension := testdatagen.MakePendingSITDurationUpdate(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), mtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})

	suite.Run("Returns an error when shipment ID from SIT extension and shipment ID found do not match", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		otherMtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		sitExtension := testdatagen.MakePendingSITDurationUpdate(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), otherMtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), otherMtoShipment.ID.String())
	})

	suite.Run("Updates the shipment's SIT days allowance and the SIT extension's status and approved days if all fields are valid", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: models.IntPointer(20),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		sitExtension := testdatagen.MakePendingSITDurationUpdate(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		// existing SITDaysAllowance plus new approved days
		newSITDaysAllowance := int(40)
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		updatedShipment, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), mtoShipment.ID, sitExtension.ID, approvedDays, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITDurationUpdate
		err = suite.DB().Find(&sitExtensionInDB, sitExtension.ID)
		suite.NoError(err)

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(newSITDaysAllowance, *updatedShipment.SITDaysAllowance)
		suite.Equal(approvedDays, *sitExtensionInDB.ApprovedDays)
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(models.SITExtensionStatusApproved, sitExtensionInDB.Status)
		suite.Equal(models.MoveStatusAPPROVED, shipmentInDB.MoveTaskOrder.Status)
	})

	suite.Run("Sets move to approvals requested if there are remaining pending SIT extensions", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: models.IntPointer(20),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		sitExtensionToBeApproved := testdatagen.MakePendingSITDurationUpdate(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		// Pending SIT Extension that won't be approved or denied
		testdatagen.MakePendingSITDurationUpdate(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		_, err := sitExtensionApprover.ApproveSITExtension(suite.AppContextForTest(), mtoShipment.ID, sitExtensionToBeApproved.ID, approvedDays, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, shipmentInDB.MoveTaskOrder.Status)
	})
}
