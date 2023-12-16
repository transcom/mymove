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

func (suite *SitExtensionServiceSuite) TestDenySITExtension() {
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionDenier := NewSITExtensionDenier(moveRouter)

	suite.Run("Returns an error when shipment is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionDenier.DenySITExtension(session, nonexistentUUID, nonexistentUUID, &officeRemarks, convertToMembersExpense, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when SIT extension is not found", func() {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := sitExtensionDenier.DenySITExtension(session, mtoShipment.ID, nonexistentUUID, &officeRemarks, convertToMembersExpense, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when etag does not match", func() {
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), nil, nil)
		mtoShipment := sitExtension.MTOShipment
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := sitExtensionDenier.DenySITExtension(session, mtoShipment.ID, sitExtension.ID, &officeRemarks, convertToMembersExpense, eTag)

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
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := sitExtensionDenier.DenySITExtension(session, otherMtoShipment.ID, sitExtension.ID, &officeRemarks, convertToMembersExpense, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), otherMtoShipment.ID.String())
	})

	suite.Run("Updates the SIT extension's status to DENIED and approves move when all fields are valid", func() {
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
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		updatedShipment, err := sitExtensionDenier.DenySITExtension(session, mtoShipment.ID, sitExtension.ID, &officeRemarks, convertToMembersExpense, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITDurationUpdate
		err = suite.DB().Find(&sitExtensionInDB, sitExtension.ID)
		suite.NoError(err)

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(convertToMembersExpense, sitExtensionInDB.MembersExpense)
		suite.Equal(models.SITExtensionStatusDenied, sitExtensionInDB.Status)
		suite.Equal(models.MoveStatusAPPROVED, shipmentInDB.MoveTaskOrder.Status)
	})

	suite.Run("Updates the SIT extension's status to DENIED and updates members_expense to TRUE when 'Convert to Member's Expense' is chosen.", func() {
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
		officeRemarks := "office remarks"
		convertToMembersExpense := true
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		updatedShipment, err := sitExtensionDenier.DenySITExtension(session, mtoShipment.ID, sitExtension.ID, &officeRemarks, convertToMembersExpense, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITDurationUpdate
		err = suite.DB().Find(&sitExtensionInDB, sitExtension.ID)
		suite.NoError(err)

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(convertToMembersExpense, sitExtensionInDB.MembersExpense)
		suite.Equal(models.SITExtensionStatusDenied, sitExtensionInDB.Status)
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
		sitExtensionToBeDenied := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
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
		officeRemarks := "office remarks"
		convertToMembersExpense := false
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := sitExtensionDenier.DenySITExtension(session, mtoShipment.ID, sitExtensionToBeDenied.ID, &officeRemarks, convertToMembersExpense, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, shipmentInDB.MoveTaskOrder.Status)
	})
}
