package move

import (
	"fmt"

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

func (suite *MoveServiceSuite) TestSubmitted() {
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
		for _, invalidStatus := range invalidStatuses {
			move.Status = invalidStatus.status

			err := moveRouter.Submit(&move)
			suite.Error(err)
			suite.Contains(err.Error(), "Cannot move to Submitted state for TOO review when the Move is not in Draft status")
			suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", invalidStatus.status))
		}
	})

	suite.Run("moves going to the services counselor return errors if the move doesn't have DRAFT status", func() {
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
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, invalidStatus := range invalidStatuses {
			move.Status = invalidStatus.status

			err := moveRouter.Submit(&move)
			suite.Error(err)
			suite.Contains(err.Error(), "Cannot move to NeedsServiceCounseling state when the Move is not in Draft status")
			suite.Contains(err.Error(), fmt.Sprintf("Its current status is: %s", invalidStatus.status))
		}
	})
}
