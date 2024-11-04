package move

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MoveServiceSuite) TestMoveApproval() {
	moveRouter := NewMoveRouter()

	suite.Run("from valid statuses", func() {
		move := factory.BuildMove(nil, nil, nil)
		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Submitted", models.MoveStatusSUBMITTED},
			{"Approvals Requested", models.MoveStatusAPPROVALSREQUESTED},
			{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
			{"Approved", models.MoveStatusAPPROVED},
		}
		for _, validStatus := range validStatuses {
			move.Status = validStatus.status

			err := moveRouter.Approve(suite.AppContextForTest(), &move)

			suite.NoError(err)
			suite.Equal(models.MoveStatusAPPROVED, move.Status)
		}
	})

	suite.Run("from invalid statuses", func() {
		move := factory.BuildMove(nil, nil, nil)
		invalidStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Draft", models.MoveStatusDRAFT},
			{"Canceled", models.MoveStatusCANCELED},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, invalidStatus := range invalidStatuses {
			move.Status = invalidStatus.status

			err := moveRouter.Approve(suite.AppContextForTest(), &move)

			suite.Error(err)
			suite.Contains(err.Error(), "A move can only be approved if it's in one of these states")
			suite.Contains(err.Error(), fmt.Sprintf("However, its current status is: %s", invalidStatus.status))
		}
	})
}

func (suite *MoveServiceSuite) TestMoveSubmission() {
	moveRouter := NewMoveRouter()

	suite.Run("returns error when needsServicesCounseling cannot find move", func() {
		// Under test: MoveRouter.Submit
		// Mocked: None
		// Set up: Submit a move without an OrdersID
		// Expected outcome: Error on ordersID
		var move models.Move
		newSignedCertification := factory.BuildSignedCertification(suite.DB(), nil, nil)
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.Error(err)
		suite.Contains(err.Error(), "Not found looking for move.OrdersID")
	})

	suite.Run("returns error when OriginDutyLocation is missing", func() {
		// Under test: MoveRouter.Submit
		// Set up: Submit a move without an originDutyLocation
		// Expected outcome: Error on ordersID
		move := factory.BuildMove(suite.DB(), nil, nil)
		order := move.Orders
		order.OriginDutyLocation = nil
		order.OriginDutyLocationID = nil
		suite.NoError(suite.DB().Update(&order))
		newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.Error(err)
		suite.Contains(err.Error(), "orders missing OriginDutyLocation")
	})

	suite.Run("moves with amended orders are set to APPROVALSREQUESTED status", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Submit an approved move with an orders record
		// Expected outcome: move status updated to APPROVALSREQUESTED
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Document{},
				Type:  &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.RouteAfterAmendingOrders(suite.AppContextForTest(), &move)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
	})

	suite.Run("moves with amended orders return an error if in CANCELLED status", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Create a CANCELLED move without an OrdersID
		// Expected outcome: Error on ordersID
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Document{},
				Type:  &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusCANCELED,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.RouteAfterAmendingOrders(suite.AppContextForTest(), &move)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("The status for the move with ID %s can not be sent to 'Approvals Requested' if the status is cancelled.", move.ID))
	})

	suite.Run("moves with amended orders that already had amended orders go into the 'Approvals Requested' status and have a nil value for 'AmendedOrdersAcknowledgedAt", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Create a move amended orders acknowledged, then submit with amended orders
		// Expected outcome: Status goes to APPROVALSREQUESTED and timestamp is cleared
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					// we need a time here that's non-nil
					AmendedOrdersAcknowledgedAt: models.TimePointer(testdatagen.DateInsidePerformancePeriod),
				},
			},
			{
				Model: models.Document{},
				Type:  &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model:    order,
				LinkOnly: true,
			},
		}, nil)
		suite.NotNil(move.Orders.AmendedOrdersAcknowledgedAt)

		err := moveRouter.RouteAfterAmendingOrders(suite.AppContextForTest(), &move)
		suite.NoError(err)
		var updatedOrders models.Order
		err = suite.DB().Find(&updatedOrders, move.OrdersID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Nil(updatedOrders.AmendedOrdersAcknowledgedAt)
	})

	suite.Run("moves going to the TOO return errors if the move doesn't have DRAFT status", func() {
		// Under test: MoveRouter.Submit
		// Set up: Create a move that is not in DRAFT status, submit a move to other statuses
		// Expected outcome: Error
		move := factory.BuildMove(suite.DB(), nil, nil)

		invalidStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Approvals Requested", models.MoveStatusAPPROVALSREQUESTED},
			{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
			{"Submitted", models.MoveStatusSUBMITTED},
			{"Approved", models.MoveStatusAPPROVED},
			{"Canceled", models.MoveStatusCANCELED},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, tt := range invalidStatuses {
			move.Status = tt.status
			newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil)
			err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
			suite.Error(err)
			suite.Contains(err.Error(), "Cannot move to Submitted state for TOO review when the Move is not in Draft status")
			suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", tt.status))
		}
	})

	suite.Run("moves going to the services counselor return errors if the move doesn't have DRAFT/NEEDS SERVICE COUNSELING status", func() {
		// Under test: MoveRouter.Submit
		// Set up: Create a move that should go to services counselor, but doesn't have DRAFT or NEEDS SERVICE COUNSELING STATUS
		// Expected outcome: Error

		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ProvidesServicesCounseling: true,
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		invalidStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Approvals Requested", models.MoveStatusAPPROVALSREQUESTED},
			{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
			{"Submitted", models.MoveStatusSUBMITTED},
			{"Approved", models.MoveStatusAPPROVED},
			{"Canceled", models.MoveStatusCANCELED},
		}
		for _, tt := range invalidStatuses {
			move.Status = tt.status
			newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil)
			err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
			suite.Error(err)
			suite.Contains(err.Error(), "Cannot move to NeedsServiceCounseling state when the Move is not in Draft status")
			suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", tt.status))
		}
	})

	suite.Run("Moves are routed correctly and SignedCertification is created", func() {
		// Under test: MoveRouter.Submit (both routing to services counselor and office user)
		// Set up: Create moves and SignedCertification
		// Expected outcome: signed cert is created and move status is updated
		tests := []struct {
			desc                       string
			ProvidesServicesCounseling bool
			moveStatus                 models.MoveStatus
		}{
			{"Routes to Service Counseling", true, models.MoveStatusNeedsServiceCounseling},
			{"Routes to office user", false, models.MoveStatusSUBMITTED},
		}
		for _, tt := range tests {
			suite.Run(tt.desc, func() {
				move := factory.BuildMove(suite.DB(), []factory.Customization{
					{
						Model: models.DutyLocation{
							ProvidesServicesCounseling: tt.ProvidesServicesCounseling,
						},
						Type: &factory.DutyLocations.OriginDutyLocation,
					},
					{
						Model: models.Move{
							Status: models.MoveStatusDRAFT,
						},
					},
				}, nil)

				shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
					{
						Model: models.MTOShipment{
							Status:       models.MTOShipmentStatusDraft,
							ShipmentType: models.MTOShipmentTypePPM,
						},
					},
					{
						Model:    move,
						LinkOnly: true,
					},
				}, nil)

				ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
					{
						Model: models.PPMShipment{
							Status: models.PPMShipmentStatusDraft,
						},
					},
				}, nil)
				move.MTOShipments = models.MTOShipments{shipment}
				move.MTOShipments[0].PPMShipment = &ppmShipment

				newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
					{
						Model:    move,
						LinkOnly: true,
					},
				}, nil)
				err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
				suite.NoError(err)
				err = suite.DB().Where("move_id = $1", move.ID).First(&newSignedCertification)
				suite.NoError(err)
				suite.NotNil(newSignedCertification)

				err = suite.DB().Find(&move, move.ID)
				suite.NoError(err)
				suite.Equal(tt.moveStatus, move.Status)
			})
		}
	})

	suite.Run("Returns error if signedCertificate is missing", func() {
		// Under test: MoveRouter.Submit (both routing to services counselor and office user)
		// Set up: Create moves and SignedCertification
		// Expected outcome: signed cert is created and move status is updated
		tests := []struct {
			desc                       string
			ProvidesServicesCounseling bool
		}{
			{"Routing to Service Counseling", true},
			{"Routing to office user", false},
		}
		for _, tt := range tests {
			suite.Run(tt.desc, func() {
				move := factory.BuildMove(suite.DB(), []factory.Customization{
					{
						Model: models.DutyLocation{
							ProvidesServicesCounseling: tt.ProvidesServicesCounseling,
						},
						Type: &factory.DutyLocations.OriginDutyLocation,
					},
					{
						Model: models.Move{
							Status: models.MoveStatusDRAFT,
						},
					},
				}, nil)
				err := moveRouter.Submit(suite.AppContextForTest(), &move, nil)
				suite.Error(err)
				suite.Contains(err.Error(), "signedCertification is required")
			})
		}
	})

	suite.Run("PPM status changes to Submitted", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusDraft,
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusDraft,
				},
			},
		}, nil)
		move.MTOShipments = models.MTOShipments{hhgShipment}
		move.MTOShipments[0].PPMShipment = &ppmShipment

		newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)

		suite.NoError(err)
		suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
		suite.Equal(models.MTOShipmentStatusSubmitted, move.MTOShipments[0].Status, "expected Submitted")
		suite.Equal(models.PPMShipmentStatusSubmitted, move.MTOShipments[0].PPMShipment.Status, "expected Submitted")
	})

	suite.Run("PPM Actual Expense Reimbursement is true for Civilian Employee on submit", func() {
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ProvidesServicesCounseling: true,
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
			{
				Model: models.Move{
					Status: models.MoveStatusDRAFT,
				},
			},
			{
				Model: models.Order{
					Grade: models.ServiceMemberGradeCIVILIANEMPLOYEE.Pointer(),
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusDraft,
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					Status:                       models.PPMShipmentStatusDraft,
					IsActualExpenseReimbursement: models.BoolPointer(false),
				},
			},
		}, nil)
		move.MTOShipments = models.MTOShipments{shipment}
		move.MTOShipments[0].PPMShipment = &ppmShipment

		newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.NoError(err)
		suite.NotNil(newSignedCertification)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.True(*move.MTOShipments[0].PPMShipment.IsActualExpenseReimbursement)
	})

	suite.Run("PPM Actual Expense Reimbursement is false for non-civilian on submit", func() {
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ProvidesServicesCounseling: true,
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
			{
				Model: models.Move{
					Status: models.MoveStatusDRAFT,
				},
			},
			{
				Model: models.Order{
					Grade: models.ServiceMemberGradeE1.Pointer(),
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusDraft,
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					Status:                       models.PPMShipmentStatusDraft,
					IsActualExpenseReimbursement: models.BoolPointer(true),
				},
			},
		}, nil)
		move.MTOShipments = models.MTOShipments{shipment}
		move.MTOShipments[0].PPMShipment = &ppmShipment

		newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.NoError(err)
		suite.NotNil(newSignedCertification)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.False(*move.MTOShipments[0].PPMShipment.IsActualExpenseReimbursement)
	})
}

func (suite *MoveServiceSuite) TestMoveCancelation() {
	moveRouter := NewMoveRouter()

	suite.Run("Cancel move with no shipments", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		err := moveRouter.Cancel(suite.AppContextForTest(), &move)
		suite.NoError(err)

		suite.Equal(models.MoveStatusCANCELED, move.Status)
	})

	suite.Run("Cancel move with HHG", func() {
		move := factory.BuildMoveWithShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusDRAFT,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		err := moveRouter.Cancel(suite.AppContextForTest(), &move)
		suite.NoError(err)

		_ = suite.DB().Reload(&move.MTOShipments)
		suite.Equal(models.MoveStatusCANCELED, move.Status)
		suite.Equal(models.MTOShipmentStatusCanceled, move.MTOShipments[0].Status)
	})

	suite.Run("Cancel move with PPM", func() {
		move := factory.BuildMove(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusDRAFT,
				},
			},
		}, nil)

		ppm := factory.BuildPPMShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusSubmitted,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		err := moveRouter.Cancel(suite.AppContextForTest(), &move)
		suite.NoError(err)
		_ = suite.DB().Reload(&move.MTOShipments)
		suite.Equal(models.MoveStatusCANCELED, move.Status)
		ppms, _ := models.FetchPPMShipmentByPPMShipmentID(suite.AppContextForTest().DB(), ppm.ID)
		suite.Equal(models.PPMShipmentStatusCanceled, ppms.Status)
	})
}

func (suite *MoveServiceSuite) TestSendToOfficeUser() {
	moveRouter := NewMoveRouter()

	suite.Run("from valid statuses", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Draft", models.MoveStatusDRAFT},
			{"Submitted", models.MoveStatusSUBMITTED},
			{"Approved", models.MoveStatusAPPROVED},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
			{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
		}
		for _, tt := range validStatuses {
			move.Status = tt.status

			err := moveRouter.SendToOfficeUser(suite.AppContextForTest(), &move)

			suite.NoError(err)
			suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
			suite.NotNil(move.ApprovalsRequestedAt)
		}
	})

	suite.Run("from invalid statuses", func() {
		move := factory.BuildMove(nil, nil, nil)
		invalidStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Canceled", models.MoveStatusCANCELED},
		}
		for _, tt := range invalidStatuses {
			move.Status = tt.status

			err := moveRouter.SendToOfficeUser(suite.AppContextForTest(), &move)

			suite.Error(err)
			suite.Contains(err.Error(), fmt.Sprintf("The status for the move with ID %s", move.ID))
			suite.Contains(err.Error(), "can not be sent to 'Approvals Requested' if the status is cancelled.")
		}
	})

	suite.Run("from APPROVALS REQUESTED status", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		err := moveRouter.SendToOfficeUser(suite.AppContextForTest(), &move)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})
}

func (suite *MoveServiceSuite) TestApproveOrRequestApproval() {
	moveRouter := NewMoveRouter()

	suite.Run("approves the move if TOO no longer has actions to perform", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})

	suite.Run("does not approve the move if excess weight risk exists and has not been acknowledged", func() {
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})

	suite.Run("does not approve the move if unreviewed service items exist", func() {
		_, move := suite.createServiceItem()

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})

	suite.Run("does not approve the move if unacknowledged amended orders exist", func() {
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
		amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    amendedDocument,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
			{
				Model:    amendedDocument,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
			{
				Model:    amendedDocument.ServiceMember,
				LinkOnly: true,
			},
		}, nil)

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})

	suite.Run("does not approve the move if unreviewed SIT extensions exist", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.Equal(move.ApprovalsRequestedAt.Format(time.RFC3339), moveInDB.ApprovalsRequestedAt.Format(time.RFC3339))
	})
}

func (suite *MoveServiceSuite) TestCompleteServiceCounseling() {
	moveRouter := NewMoveRouter()

	suite.Run("status changed to service counseling completed", func() {
		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusNeedsServiceCounseling)
		hhgShipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		move.MTOShipments = models.MTOShipments{hhgShipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusServiceCounselingCompleted, move.Status)
	})

	suite.Run("status changed to approved", func() {
		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusNeedsServiceCounseling)
		ppmShipment := factory.BuildPPMShipment(nil, nil, nil)
		move.MTOShipments = models.MTOShipments{ppmShipment.Shipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)
	})

	suite.Run("no shipments present", func() {
		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusNeedsServiceCounseling)

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "No shipments associated with move")
	})

	suite.Run("move has unexpected existing status", func() {
		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusDRAFT)
		ppmShipment := factory.BuildPPMShipment(nil, nil, nil)
		move.MTOShipments = models.MTOShipments{ppmShipment.Shipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "The status for the Move")
	})

	suite.Run("NTS-release with no facility info", func() {
		move := factory.BuildStubbedMoveWithStatus(models.MoveStatusNeedsServiceCounseling)
		ntsrShipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID:           uuid.Must(uuid.NewV4()),
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		move.MTOShipments = models.MTOShipments{ntsrShipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "NTS-release shipment must include facility info")
	})
}

func (suite *MoveServiceSuite) createServiceItem() (models.MTOServiceItem, models.Move) {
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return serviceItem, move
}
