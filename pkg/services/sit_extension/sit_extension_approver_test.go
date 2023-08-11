package sitextension

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
)

func (suite *SitExtensionServiceSuite) TestApproveSITExtension() {
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionApprover := NewSITExtensionApprover(moveRouter)

	suite.Run("Returns an error when shipment is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		approvedDays := int(20)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := ""

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := sitExtensionApprover.ApproveSITExtension(session, nonexistentUUID, nonexistentUUID, approvedDays, requestReason, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when SIT extension is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		approvedDays := int(20)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionApprover.ApproveSITExtension(session, mtoShipment.ID, nonexistentUUID, approvedDays, requestReason, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when etag does not match", func() {
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), nil, nil)
		mtoShipment := sitExtension.MTOShipment
		approvedDays := int(20)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionApprover.ApproveSITExtension(session, mtoShipment.ID, sitExtension.ID, approvedDays, requestReason, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})

	suite.Run("Returns an error when shipment ID from SIT extension and shipment ID found do not match", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		otherMtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		approvedDays := int(20)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionApprover.ApproveSITExtension(session, otherMtoShipment.ID, sitExtension.ID, approvedDays, requestReason, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), otherMtoShipment.ID.String())
	})

	suite.Run("Returns an error when SIT duration update reduces the SIT allowance to < 1 day", func() {
		sitDaysAllowance := 20
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: &sitDaysAllowance,
				},
			}}, nil)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		approvedDays := int(-30)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionApprover.ApproveSITExtension(session, mtoShipment.ID, sitExtension.ID, approvedDays, requestReason, &officeRemarks, eTag)

		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("can't reduce a SIT duration to less than one day", err.Error())
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
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		approvedDays := int(20)
		// existing SITDaysAllowance plus new approved days
		newSITDaysAllowance := int(40)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		updatedShipment, err := sitExtensionApprover.ApproveSITExtension(session, mtoShipment.ID, sitExtension.ID, approvedDays, requestReason, &officeRemarks, eTag)
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
		suite.Equal(sitExtension.RequestReason, sitExtensionInDB.RequestReason)
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
		sitExtensionToBeApproved := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		// Pending SIT Extension that won't be approved or denied
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		approvedDays := int(20)
		requestReason := models.SITExtensionRequestReasonAwaitingCompletionOfResidence
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionApprover.ApproveSITExtension(session, mtoShipment.ID, sitExtensionToBeApproved.ID, approvedDays, requestReason, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, shipmentInDB.MoveTaskOrder.Status)
	})
}
