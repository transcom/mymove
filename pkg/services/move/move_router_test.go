package move

import (
	"fmt"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MoveServiceSuite) TestMoveApproval() {
	moveRouter := NewMoveRouter()

	suite.Run("from valid statuses", func() {
		orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{Stub: true})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: orders,
			Stub:  true})
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
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})
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
		newSignedCertification := testdatagen.MakeDefaultSignedCertification(suite.DB())
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.Error(err)
		suite.Contains(err.Error(), "Not found looking for move.OrdersID")
	})

	suite.Run("returns error when OriginDutyLocation is missing", func() {
		// Under test: MoveRouter.Submit
		// Set up: Submit a move without an originDutyLocation
		// Expected outcome: Error on ordersID
		move := testdatagen.MakeDefaultMove(suite.DB())
		order := move.Orders
		order.OriginDutyLocation = nil
		order.OriginDutyLocationID = nil
		suite.NoError(suite.DB().Update(&order))
		newSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
			SignedCertification: models.SignedCertification{
				MoveID: move.ID,
			},
			Stub: true,
		})
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
		suite.Error(err)
		suite.Contains(err.Error(), "orders missing OriginDutyLocation")
	})

	suite.Run("moves with amended orders are set to APPROVALSREQUESTED status", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Submit an approved move with an orders record
		// Expected outcome: move status updated to APPROVALSREQUESTED
		document := testdatagen.MakeDefaultDocument(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ID:                    uuid.Must(uuid.NewV4()),
				UploadedAmendedOrders: &document,
			},
		})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
			Order: order,
		})

		err := moveRouter.RouteAfterAmendingOrders(suite.AppContextForTest(), &move)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
	})

	suite.Run("moves with amended orders return an error if in CANCELLED status", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Create a CANCELLED move without an OrdersID
		// Expected outcome: Error on ordersID
		document := testdatagen.MakeDefaultDocument(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ID:                    uuid.Must(uuid.NewV4()),
				UploadedAmendedOrders: &document,
			},
		})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusCANCELED,
			},
			Order: order,
		})

		err := moveRouter.RouteAfterAmendingOrders(suite.AppContextForTest(), &move)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("The status for the move with ID %s can not be sent to 'Approvals Requested' if the status is cancelled.", move.ID))
	})

	suite.Run("moves with amended orders that already had amended orders go into the 'Approvals Requested' status and have a nil value for 'AmendedOrdersAcknowledgedAt", func() {
		// Under test: MoveRouter.RouteAfterAmendingOrders
		// Set up: Create a move amended orders acknowledged, then submit with amended orders
		// Expected outcome: Status goes to APPROVALSREQUESTED and timestamp is cleared
		document := testdatagen.MakeDefaultDocument(suite.DB())
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ID:                    uuid.Must(uuid.NewV4()),
				UploadedAmendedOrders: &document,
				// we need a time here that's non-nil
				AmendedOrdersAcknowledgedAt: swag.Time(testdatagen.DateInsidePerformancePeriod),
			},
		})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
			Order: order,
		})
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
		move := testdatagen.MakeDefaultMove(suite.DB())

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
			newSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
				SignedCertification: models.SignedCertification{
					MoveID: move.ID,
				},
				Stub: true,
			})
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
		dutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})
		assertions := testdatagen.Assertions{
			Order: models.Order{
				OriginDutyLocation: &dutyLocation,
			},
		}
		move := testdatagen.MakeMove(suite.DB(), assertions)

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
			newSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
				Move: move,
				Stub: true,
			})
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
				dutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
					DutyLocation: models.DutyLocation{
						ProvidesServicesCounseling: tt.ProvidesServicesCounseling,
					},
				})

				move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
					Order: models.Order{
						OriginDutyLocation: &dutyLocation,
					},
					Move: models.Move{
						Status: models.MoveStatusDRAFT,
					},
				})

				shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
					MTOShipment: models.MTOShipment{
						Status:          models.MTOShipmentStatusDraft,
						ShipmentType:    models.MTOShipmentTypePPM,
						MoveTaskOrder:   move,
						MoveTaskOrderID: move.ID,
					},
					Stub: true,
				})
				ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
					PPMShipment: models.PPMShipment{
						Status: models.PPMShipmentStatusDraft,
					},
					Stub: true,
				})

				move.MTOShipments = models.MTOShipments{shipment}
				move.MTOShipments[0].PPMShipment = &ppmShipment

				newSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
					SignedCertification: models.SignedCertification{
						MoveID: move.ID,
					},
					Stub: true,
				})
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
				dutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
					DutyLocation: models.DutyLocation{
						ProvidesServicesCounseling: tt.ProvidesServicesCounseling,
					},
				})

				move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
					Order: models.Order{
						OriginDutyLocation: &dutyLocation,
					},
					Move: models.Move{
						Status: models.MoveStatusDRAFT,
					},
				})

				err := moveRouter.Submit(suite.AppContextForTest(), &move, nil)
				suite.Error(err)
				suite.Contains(err.Error(), "signedCertification is required")
			})
		}
	})

	suite.Run("PPM status changes to Submitted", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())

		hhgShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusDraft,
				ShipmentType:    models.MTOShipmentTypePPM,
				MoveTaskOrder:   move,
				MoveTaskOrderID: move.ID,
			},
			Stub: true,
		})
		ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusDraft,
			},
			Stub: true,
		})

		move.MTOShipments = models.MTOShipments{hhgShipment}
		move.MTOShipments[0].PPMShipment = &ppmShipment

		move.MTOShipments = models.MTOShipments{hhgShipment}
		move.MTOShipments[0].PPMShipment = &ppmShipment
		newSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
			SignedCertification: models.SignedCertification{
				MoveID: move.ID,
			},
			Stub: true,
		})
		err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)

		suite.NoError(err)
		suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
		suite.Equal(models.MTOShipmentStatusSubmitted, move.MTOShipments[0].Status, "expected Submitted")
		suite.Equal(models.PPMShipmentStatusSubmitted, move.MTOShipments[0].PPMShipment.Status, "expected Submitted")
	})
}

func (suite *MoveServiceSuite) TestMoveCancellation() {
	moveRouter := NewMoveRouter()

	suite.Run("defaults to nil reason if empty string provided", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

		err := moveRouter.Cancel(suite.AppContextForTest(), "", &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusCANCELED, move.Status, "expected Canceled")
		suite.Nil(move.CancelReason)
	})

	suite.Run("adds reason if provided", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

		reason := "SM's orders revoked"
		err := moveRouter.Cancel(suite.AppContextForTest(), reason, &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusCANCELED, move.Status, "expected Canceled")
		suite.Equal(&reason, move.CancelReason, "expected 'SM's orders revoked'")
	})

	suite.Run("cancels PPM and Order when move is canceled", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

		// Create PPM on this move
		advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
		ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				Move:      move,
				MoveID:    move.ID,
				Status:    models.PPMStatusDRAFT,
				Advance:   &advance,
				AdvanceID: &advance.ID,
			},
			Stub: true,
		})
		move.PersonallyProcuredMoves = append(move.PersonallyProcuredMoves, ppm)

		err := moveRouter.Cancel(suite.AppContextForTest(), "", &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusCANCELED, move.Status, "expected Canceled")
		suite.Equal(models.PPMStatusCANCELED, move.PersonallyProcuredMoves[0].Status, "expected Canceled")
		suite.Equal(models.OrderStatusCANCELED, move.Orders.Status, "expected Canceled")
	})

	suite.Run("from valid statuses", func() {
		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Submitted", models.MoveStatusSUBMITTED},
			{"Approvals Requested", models.MoveStatusAPPROVALSREQUESTED},
			{"Service Counseling Completed", models.MoveStatusServiceCounselingCompleted},
			{"Approved", models.MoveStatusAPPROVED},
			{"Draft", models.MoveStatusDRAFT},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, tt := range validStatuses {
			suite.Run(tt.desc, func() {
				move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

				move.Status = tt.status
				move.Orders.Status = models.OrderStatusSUBMITTED

				err := moveRouter.Cancel(suite.AppContextForTest(), "", &move)

				suite.NoError(err)
				suite.Equal(models.MoveStatusCANCELED, move.Status)
			})
		}
	})

	suite.Run("from invalid statuses", func() {
		invalidStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Canceled", models.MoveStatusCANCELED},
		}
		for _, tt := range invalidStatuses {
			suite.Run(tt.desc, func() {
				move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

				move.Status = tt.status

				err := moveRouter.Cancel(suite.AppContextForTest(), "", &move)

				suite.Error(err)
				suite.Contains(err.Error(), "cannot cancel a move that is already canceled")
			})
		}
	})
}

func (suite *MoveServiceSuite) TestSendToOfficeUser() {
	moveRouter := NewMoveRouter()

	suite.Run("from valid statuses", func() {
		orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: orders})
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
		}
	})

	suite.Run("from invalid statuses", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})
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
}

func (suite *MoveServiceSuite) TestApproveOrRequestApproval() {
	moveRouter := NewMoveRouter()

	suite.Run("approves the move if TOO no longer has actions to perform", func() {
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
	})

	suite.Run("does not approve the move if excess weight risk exists and has not been acknowledged", func() {
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ExcessWeightQualifiedAt: &now,
			},
		})

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
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
	})

	suite.Run("does not approve the move if unacknowledged amended orders exist", func() {
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &amendedDocument.ID,
				Document:   amendedDocument,
				UploaderID: amendedDocument.ServiceMember.UserID,
			},
			UserUploader: userUploader,
		})

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &amendedDocument,
				UploadedAmendedOrdersID: &amendedDocument.ID,
				ServiceMember:           amendedDocument.ServiceMember,
				ServiceMemberID:         amendedDocument.ServiceMemberID,
			},
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
	})

	suite.Run("does not approve the move if unreviewed SIT extensions exist", func() {
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		updatedMove, err := moveRouter.ApproveOrRequestApproval(suite.AppContextForTest(), move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
	})
}

func (suite *MoveServiceSuite) TestCompleteServiceCounseling() {
	moveRouter := NewMoveRouter()

	suite.Run("status changed to service counseling completed", func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusNeedsServiceCounseling)
		hhgShipment := testdatagen.MakeStubbedShipment(suite.DB())
		move.MTOShipments = models.MTOShipments{hhgShipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusServiceCounselingCompleted, move.Status)
	})

	suite.Run("status changed to approved", func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusNeedsServiceCounseling)
		ppmShipment := testdatagen.MakeStubbedPPMShipment(suite.DB())
		move.MTOShipments = models.MTOShipments{ppmShipment.Shipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)
	})

	suite.Run("no shipments present", func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusNeedsServiceCounseling)

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "No shipments associated with move")
	})

	suite.Run("move has unexpected existing status", func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusDRAFT)
		ppmShipment := testdatagen.MakeStubbedPPMShipment(suite.DB())
		move.MTOShipments = models.MTOShipments{ppmShipment.Shipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "The status for the Move")
	})

	suite.Run("NTS-release with no facility info", func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusNeedsServiceCounseling)
		ntsrShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ID:           uuid.Must(uuid.NewV4()),
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
			Move: move,
			Stub: true,
		})
		move.MTOShipments = models.MTOShipments{ntsrShipment}

		err := moveRouter.CompleteServiceCounseling(suite.AppContextForTest(), &move)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "NTS-release shipment must include facility info")
	})
}

func (suite *MoveServiceSuite) createServiceItem() (models.MTOServiceItem, models.Move) {
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})

	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	return serviceItem, move
}
