package move

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestMoveApproval() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})
	moveRouter := NewMoveRouter(suite.DB(), suite.logger)

	suite.Run("from valid statuses", func() {
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

			err := moveRouter.Approve(&move)

			suite.NoError(err)
			suite.Equal(models.MoveStatusAPPROVED, move.Status)
		}
	})

	suite.Run("from invalid statuses", func() {
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

			err := moveRouter.Approve(&move)

			suite.Error(err)
			suite.Contains(err.Error(), "A move can only be approved if it's in one of these states")
			suite.Contains(err.Error(), fmt.Sprintf("However, its current status is: %s", invalidStatus.status))
		}
	})
}

func (suite *MoveServiceSuite) TestMoveSubmission() {
	moveRouter := NewMoveRouter(suite.DB(), suite.logger)

	suite.Run("returns error when needsServicesCounseling cannot find move", func() {
		var move models.Move
		err := moveRouter.Submit(&move)
		suite.Error(err)
		suite.Contains(err.Error(), "not found looking for move.OrdersID")
	})

	suite.Run("returns error when OriginDutyStation is missing", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		order := move.Orders
		order.OriginDutyStation = nil
		order.OriginDutyStationID = nil
		suite.NoError(suite.DB().Update(&order))

		err := moveRouter.Submit(&move)
		suite.Error(err)
		suite.Contains(err.Error(), "orders missing OriginDutyStation")
	})

	suite.Run("moves with amended orders are set to APPROVALSREQUESTED status", func() {
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

		err := moveRouter.Submit(&move)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
	})

	suite.Run("moves with amended orders return an error if in CANCELLED status", func() {
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

		err := moveRouter.Submit(&move)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("The status for the move with ID %s can not be sent to 'Approvals Requested' if the status is cancelled.", move.ID))
	})

	suite.Run("moves going to the TOO return errors if the move doesn't have DRAFT status", func() {
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
			suite.Run(tt.desc, func() {
				move.Status = tt.status

				err := moveRouter.Submit(&move)
				suite.Error(err)
				suite.Contains(err.Error(), "Cannot move to Submitted state for TOO review when the Move is not in Draft status")
				suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", tt.status))
			})
		}
	})

	suite.Run("moves going to the services counselor return errors if the move doesn't have DRAFT/NEEDS SERVICE COUNSELING status", func() {
		dutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
			DutyStation: models.DutyStation{
				ProvidesServicesCounseling: true,
			},
		})
		assertions := testdatagen.Assertions{
			Order: models.Order{
				OriginDutyStation: &dutyStation,
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
			suite.Run(tt.desc, func() {
				move.Status = tt.status

				err := moveRouter.Submit(&move)
				suite.Error(err)
				suite.Contains(err.Error(), "Cannot move to NeedsServiceCounseling state when the Move is not in Draft status")
				suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", tt.status))
			})
		}
	})

	suite.Run("PPM status changes to Submitted", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())

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

		err := moveRouter.Submit(&move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusSUBMITTED, move.Status, "expected Submitted")
		suite.Equal(models.PPMStatusSUBMITTED, move.PersonallyProcuredMoves[0].Status, "expected Submitted")
	})
}

func (suite *MoveServiceSuite) TestApproveAmendedOrders() {
	moveRouter := NewMoveRouter(suite.DB(), suite.logger)

	suite.Run("approves move with no service items in requested status", func() {
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		approvedMove, approveErr := moveRouter.ApproveAmendedOrders(move.ID, move.Orders.ID)

		suite.NoError(approveErr)
		suite.Equal(models.MoveStatusAPPROVED, approvedMove.Status)
	})

	suite.Run("leaves move in APPROVALS REQUESTED status if service items are awaiting approval", func() {
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:        move,
			MTOShipment: shipment,
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusSubmitted,
			},
		})
		approvedMove, approveErr := moveRouter.ApproveAmendedOrders(move.ID, move.Orders.ID)

		suite.NoError(approveErr)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, approvedMove.Status)
	})
}

func (suite *MoveServiceSuite) TestMoveCancellation() {
	moveRouter := NewMoveRouter(suite.DB(), suite.logger)

	suite.Run("defaults to nil reason if empty string provided", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})
		err := moveRouter.Cancel("", &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusCANCELED, move.Status, "expected Canceled")
		suite.Nil(move.CancelReason)
	})

	suite.Run("adds reason if provided", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

		reason := "SM's orders revoked"
		err := moveRouter.Cancel(reason, &move)

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

		err := moveRouter.Cancel("", &move)

		suite.NoError(err)
		suite.Equal(models.MoveStatusCANCELED, move.Status, "expected Canceled")
		suite.Equal(models.PPMStatusCANCELED, move.PersonallyProcuredMoves[0].Status, "expected Canceled")
		suite.Equal(models.OrderStatusCANCELED, move.Orders.Status, "expected Canceled")
	})

	suite.Run("from valid statuses", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})

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
				move.Status = tt.status
				move.Orders.Status = models.OrderStatusSUBMITTED

				err := moveRouter.Cancel("", &move)

				suite.NoError(err)
				suite.Equal(models.MoveStatusCANCELED, move.Status)
			})
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
			suite.Run(tt.desc, func() {
				move.Status = tt.status

				err := moveRouter.Cancel("", &move)

				suite.Error(err)
				suite.Contains(err.Error(), "Cannot cancel a move that is already canceled.")
			})
		}
	})
}

func (suite *MoveServiceSuite) TestSendToOfficeUser() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Stub: true})
	moveRouter := NewMoveRouter(suite.DB(), suite.logger)

	suite.Run("from valid statuses", func() {
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
			suite.Run(tt.desc, func() {
				move.Status = tt.status

				err := moveRouter.SendToOfficeUser(&move)

				suite.NoError(err)
				suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
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
				move.Status = tt.status

				err := moveRouter.SendToOfficeUser(&move)

				suite.Error(err)
				suite.Contains(err.Error(), fmt.Sprintf("The status for the move with ID %s", move.ID))
				suite.Contains(err.Error(), "can not be sent to 'Approvals Requested' if the status is cancelled.")
			})
		}
	})
}
