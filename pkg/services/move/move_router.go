package move

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveRouter struct {
	db     *pop.Connection
	logger Logger
}

// NewMoveRouter creates a new moveRouter service
func NewMoveRouter(db *pop.Connection, logger Logger) services.MoveRouter {
	return &moveRouter{db, logger}
}

// InitNewMoveRouter initialize NewMoveRouter
func InitNewMoveRouter(db *pop.Connection, logger Logger) (services.MoveRouter, error) {
	return NewMoveRouter(
		db,
		logger), nil
}

// Submit is called when the customer submits their move. It determines whether
// to send the move to Service Counseling or directly to the TOO. If it goes to
// Service Counseling, its status becomes "Needs Service Counseling", otherwise,
// "Submitted".
func (router moveRouter) Submit(move *models.Move) error {
	var err error
	router.logMove(move)

	needsServicesCounseling, err := router.needsServiceCounseling(move)
	if err != nil {
		router.logger.Error("failure determining if a move needs services counseling", zap.Error(err))
		return err
	}
	router.logger.Info("SUCCESS: Determining if move needs services counseling or not")

	if needsServicesCounseling {
		err = router.sendToServiceCounselor(move)
		if err != nil {
			router.logger.Error("failure routing move to services counseling", zap.Error(err))
			return err
		}
		router.logger.Info("SUCCESS: Move sent to services counseling")
	} else {
		err = router.sendNewMoveToOfficeUser(move)
		if err != nil {
			router.logger.Error("failure routing move to office user / TOO queue", zap.Error(err))
			return err
		}
		router.logger.Info("SUCCESS: Move sent to office user / TOO queue")
	}

	router.logger.Info("SUCCESS: Move submitted and routed to the appropriate queue")
	return nil
}

func (router moveRouter) needsServiceCounseling(move *models.Move) (bool, error) {
	var orders models.Order
	err := router.db.Q().
		Where("orders.id = ?", move.OrdersID).
		First(&orders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			router.logger.Error("failure finding move", zap.Error(err))
			return false, services.NewNotFoundError(move.OrdersID, "looking for move.OrdersID")
		default:
			router.logger.Error("failure encountered querying for orders associated with the move", zap.Error(err))
			return false, fmt.Errorf("failure encountered querying for orders associated with the move, %s, id: %s", err.Error(), move.ID)
		}
	}

	var originDutyStation models.DutyStation

	if orders.OriginDutyStationID == nil || *orders.OriginDutyStationID == uuid.Nil {
		return false, services.NewInvalidInputError(orders.ID, err, nil, "orders missing OriginDutyStation")
	}

	originDutyStation, err = models.FetchDutyStation(router.db, *orders.OriginDutyStationID)
	if err != nil {
		router.logger.Error("failure finding the origin duty station", zap.Error(err))
		return false, services.NewInvalidInputError(*orders.OriginDutyStationID, err, nil, "unable to find origin duty station")
	}
	return originDutyStation.ProvidesServicesCounseling, nil
}

// sendToServiceCounselor makes the move available for a Service Counselor to review
func (router moveRouter) sendToServiceCounselor(move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		router.logger.Warn(fmt.Sprintf(
			"Cannot move to NeedsServiceCounseling state when the Move is not in Draft status. Its current status is: %s",
			move.Status,
		))

		return errors.Wrap(
			models.ErrInvalidTransition, fmt.Sprintf(
				"Cannot move to NeedsServiceCounseling state when the Move is not in Draft status. Its current status is: %s",
				move.Status,
			),
		)
	}

	move.Status = models.MoveStatusNeedsServiceCounseling
	now := time.Now()
	move.SubmittedAt = &now

	return nil
}

// sendNewMoveToOfficeUser makes the move available for a TOO to review
// The Submitted status indicates to the TOO that this is a new move.
func (router moveRouter) sendNewMoveToOfficeUser(move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		router.logger.Warn(fmt.Sprintf(
			"Cannot move to Submitted state for TOO review when the Move is not in Draft status. Its current status is: %s",
			move.Status))

		return errors.Wrap(models.ErrInvalidTransition, fmt.Sprintf(
			"Cannot move to Submitted state for TOO review when the Move is not in Draft status. Its current status is: %s",
			move.Status))
	}
	move.Status = models.MoveStatusSUBMITTED
	now := time.Now()
	move.SubmittedAt = &now

	// Update PPM status too
	for i := range move.PersonallyProcuredMoves {
		ppm := &move.PersonallyProcuredMoves[i]
		err := ppm.Submit(now)
		if err != nil {
			router.logger.Error("Failure submitting ppm", zap.Error(err))
			return err
		}
	}

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				router.logger.Error("Failure requesting reimbursement for ppm", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

// Approve makes the Move available to the Prime. The Prime cannot create
// Service Items unless the Move is approved.
func (router moveRouter) Approve(move *models.Move) error {
	if router.approvable(move) {
		move.Status = models.MoveStatusAPPROVED
		return nil
	}
	if router.alreadyApproved(move) {
		return nil
	}

	router.logger.Warn(fmt.Sprintf(
		"A move can only be approved if it's in one of these states: %q. However, its current status is: %s",
		validStatusesBeforeApproval, move.Status,
	))

	return errors.Wrap(
		models.ErrInvalidTransition, fmt.Sprintf(
			"A move can only be approved if it's in one of these states: %q. However, its current status is: %s",
			validStatusesBeforeApproval, move.Status,
		),
	)
}

func (router moveRouter) alreadyApproved(move *models.Move) bool {
	return move.Status == models.MoveStatusAPPROVED
}

func (router moveRouter) approvable(move *models.Move) bool {
	return statusSliceContains(validStatusesBeforeApproval, move.Status)
}

func statusSliceContains(statusSlice []models.MoveStatus, status models.MoveStatus) bool {
	for _, validStatus := range statusSlice {
		if status == validStatus {
			return true
		}
	}
	return false
}

var validStatusesBeforeApproval = []models.MoveStatus{
	models.MoveStatusSUBMITTED,
	models.MoveStatusAPPROVALSREQUESTED,
	models.MoveStatusServiceCounselingCompleted,
}

// SendToOfficeUserToReviewNewServiceItems sets the moves status to
// "Approvals Requested", which indicates to the TOO that they have new
// service items to review.
func (router moveRouter) SendToOfficeUserToReviewNewServiceItems(move *models.Move) error {
	// Do nothing if it's already in the desired state
	if move.Status == models.MoveStatusAPPROVALSREQUESTED {
		return nil
	}
	if move.Status != models.MoveStatusAPPROVED {
		return errors.Wrap(models.ErrInvalidTransition, fmt.Sprintf("The status for the Move with ID %s can only be set to 'Approvals Requested' from the 'Approved' status, but its current status is %s.", move.ID, move.Status))
	}
	move.Status = models.MoveStatusAPPROVALSREQUESTED
	return nil
}

// Cancel cancels the Move and its associated PPMs
func (router moveRouter) Cancel(reason string, move *models.Move) error {
	// We can cancel any move that isn't already complete.
	if move.Status == models.MoveStatusCANCELED {
		return errors.Wrap(models.ErrInvalidTransition, "Cancel")
	}

	move.Status = models.MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if reason != "" {
		move.CancelReason = &reason
	}

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range move.PersonallyProcuredMoves {
		err := move.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			return err
		}
	}

	// TODO: Orders can exist after related moves are canceled
	err := move.Orders.Cancel()
	if err != nil {
		return err
	}

	return nil

}

// CompleteServiceCounseling sets the move status to "Service Counseling Completed",
// which makes the move available to the TOO. This gets called when the Service
// Counselor is done reviewing the move and submits it.
func (router moveRouter) CompleteServiceCounseling(move *models.Move) error {
	if move.Status != models.MoveStatusNeedsServiceCounseling {
		return errors.Wrap(
			models.ErrInvalidTransition,
			fmt.Sprintf("The status for the Move with ID %s can only be set to 'Service Counseling Completed' from the 'Needs Service Counseling' status, but its current status is %s.",
				move.ID, move.Status,
			),
		)
	}

	now := time.Now()
	move.ServiceCounselingCompletedAt = &now
	move.Status = models.MoveStatusServiceCounselingCompleted

	return nil
}

func (router moveRouter) logMove(move *models.Move) {
	router.logger.Info("Move log",
		zap.String("Move.ID", move.ID.String()),
		zap.String("Move.Locator", move.Locator),
		zap.String("Move.Status", string(move.Status)),
		zap.String("Move.OrdersID", move.OrdersID.String()),
	)
}
