package move

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveRouter struct {
}

// NewMoveRouter creates a new moveRouter service
func NewMoveRouter() services.MoveRouter {
	return &moveRouter{}
}

// Submit is called when the customer submits their move. It determines whether
// to send the move to Service Counseling or directly to the TOO. If it goes to
// Service Counseling, its status becomes "Needs Service Counseling", otherwise,
// "Submitted".
func (router moveRouter) Submit(appCfg appconfig.AppConfig, move *models.Move) error {
	router.logMove(appCfg, move)

	needsServicesCounseling, err := router.needsServiceCounseling(appCfg, move)
	if err != nil {
		appCfg.Logger().Error("failure determining if a move needs services counseling", zap.Error(err))
		return err
	}
	appCfg.Logger().Info("SUCCESS: Determining if move needs services counseling or not")

	if needsServicesCounseling {
		err = router.sendToServiceCounselor(appCfg, move)
		if err != nil {
			appCfg.Logger().Error("failure routing move to services counseling", zap.Error(err))
			return err
		}
		appCfg.Logger().Info("SUCCESS: Move sent to services counseling")
	} else if move.Orders.UploadedAmendedOrders != nil {
		appCfg.Logger().Info("Move has amended orders")
		transactionError := appCfg.NewTransaction(func(txnAppCfg appconfig.AppConfig) error {
			err = router.SendToOfficeUser(txnAppCfg, move)
			if err != nil {
				txnAppCfg.Logger().Error("failure routing move submission with amended orders", zap.Error(err))
				return err
			}
			// Let's get the orders for this move so we can wipe out the acknowledgement if it exists already (from a prior orders amendment process)
			var ordersForMove models.Order
			err = txnAppCfg.DB().Find(&ordersForMove, move.OrdersID)
			if err != nil {
				return err
			}
			// Here we'll nil out the value (if it's set already) so that on the client-side we'll see view this change
			// in status as 'new orders' that need acknowledging by the TOO.
			// We shouldn't need more complicated logic here since we only hit this point from calling Submit().
			// Other circumstances like new MTOServiceItems will be calling SendToOfficeUser() directly.
			txnAppCfg.Logger().Info("Determining whether there is a preexisting orders acknowledgement")
			if ordersForMove.AmendedOrdersAcknowledgedAt != nil {
				txnAppCfg.Logger().Info("Move has a preexisting acknowledgement")
				ordersForMove.AmendedOrdersAcknowledgedAt = nil
				_, err = txnAppCfg.DB().ValidateAndSave(&ordersForMove)
				if err != nil {
					txnAppCfg.Logger().Error("failure resetting orders AmendedOrdersAcknowledgeAt field when routing move submission with amended orders ", zap.Error(err))
					return err
				}
				txnAppCfg.Logger().Info("Successfully reset orders acknowledgement")
			}
			return nil
		})
		if transactionError != nil {
			return transactionError
		}
		appCfg.Logger().Info("SUCCESS: Move with amended orders sent to office user / TOO queue")
	} else {
		err = router.sendNewMoveToOfficeUser(appCfg, move)
		if err != nil {
			appCfg.Logger().Error("failure routing move to office user / TOO queue", zap.Error(err))
			return err
		}
		appCfg.Logger().Info("SUCCESS: Move sent to office user / TOO queue")
	}

	appCfg.Logger().Info("SUCCESS: Move submitted and routed to the appropriate queue")
	return nil
}

func (router moveRouter) needsServiceCounseling(appCfg appconfig.AppConfig, move *models.Move) (bool, error) {
	var orders models.Order
	err := appCfg.DB().Q().
		Where("orders.id = ?", move.OrdersID).
		First(&orders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			appCfg.Logger().Error("failure finding move", zap.Error(err))
			return false, services.NewNotFoundError(move.OrdersID, "looking for move.OrdersID")
		default:
			appCfg.Logger().Error("failure encountered querying for orders associated with the move", zap.Error(err))
			return false, fmt.Errorf("failure encountered querying for orders associated with the move, %s, id: %s", err.Error(), move.ID)
		}
	}

	var originDutyStation models.DutyStation

	if orders.OriginDutyStationID == nil || *orders.OriginDutyStationID == uuid.Nil {
		return false, services.NewInvalidInputError(orders.ID, err, nil, "orders missing OriginDutyStation")
	}

	originDutyStation, err = models.FetchDutyStation(appCfg.DB(), *orders.OriginDutyStationID)
	if err != nil {
		appCfg.Logger().Error("failure finding the origin duty station", zap.Error(err))
		return false, services.NewInvalidInputError(*orders.OriginDutyStationID, err, nil, "unable to find origin duty station")
	}

	if move.ServiceCounselingCompletedAt != nil {
		return false, nil
	}

	return originDutyStation.ProvidesServicesCounseling, nil
}

// sendToServiceCounselor makes the move available for a Service Counselor to review
func (router moveRouter) sendToServiceCounselor(appCfg appconfig.AppConfig, move *models.Move) error {
	if move.Status == models.MoveStatusNeedsServiceCounseling {
		return nil
	}

	if move.Status != models.MoveStatusDRAFT {
		appCfg.Logger().Warn(fmt.Sprintf(
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
func (router moveRouter) sendNewMoveToOfficeUser(appCfg appconfig.AppConfig, move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		appCfg.Logger().Warn(fmt.Sprintf(
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
			appCfg.Logger().Error("Failure submitting ppm", zap.Error(err))
			return err
		}
	}

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				appCfg.Logger().Error("Failure requesting reimbursement for ppm", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

// Approve makes the Move available to the Prime. The Prime cannot create
// Service Items unless the Move is approved.
func (router moveRouter) Approve(appCfg appconfig.AppConfig, move *models.Move) error {
	router.logMove(appCfg, move)
	if router.approvable(move) {
		move.Status = models.MoveStatusAPPROVED
		appCfg.Logger().Info("SUCCESS: Move approved")
		return nil
	}
	if router.alreadyApproved(move) {
		return nil
	}

	appCfg.Logger().Warn(fmt.Sprintf(
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

// SendToOfficeUser sets the moves status to
// "Approvals Requested", which indicates to the TOO that they have new
// service items to review.
func (router moveRouter) SendToOfficeUser(appCfg appconfig.AppConfig, move *models.Move) error {
	router.logMove(appCfg, move)
	// Do nothing if it's already in the desired state
	if move.Status == models.MoveStatusAPPROVALSREQUESTED {
		return nil
	}
	if move.Status == models.MoveStatusCANCELED {
		errorMessage := fmt.Sprintf("The status for the move with ID %s can not be sent to 'Approvals Requested' if the status is cancelled.", move.ID)
		appCfg.Logger().Warn(errorMessage)

		return errors.Wrap(models.ErrInvalidTransition, errorMessage)
	}
	move.Status = models.MoveStatusAPPROVALSREQUESTED
	appCfg.Logger().Info("SUCCESS: Move sent to TOO to request approval")

	return nil
}

// Cancel cancels the Move and its associated PPMs
func (router moveRouter) Cancel(appCfg appconfig.AppConfig, reason string, move *models.Move) error {
	router.logMove(appCfg, move)
	// We can cancel any move that isn't already complete.
	// TODO: What does complete mean? How do we determine when a move is complete?
	if move.Status == models.MoveStatusCANCELED {
		return errors.Wrap(models.ErrInvalidTransition, "Cannot cancel a move that is already canceled.")
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

	appCfg.Logger().Info("SUCCESS: Move Canceled")
	return nil

}

// CompleteServiceCounseling sets the move status to "Service Counseling Completed",
// which makes the move available to the TOO. This gets called when the Service
// Counselor is done reviewing the move and submits it.
func (router moveRouter) CompleteServiceCounseling(appCfg appconfig.AppConfig, move *models.Move) error {
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

// ApproveAmendedOrders sets the move status to APPROVED if its status was set to
// APPROVALS REQUESTED because of the customer amending their orders.  If there are accessorial
// service items needing review from the TOO the status should remain in APPROVALS REQUESTED
func (router moveRouter) ApproveAmendedOrders(appCfg appconfig.AppConfig, moveID uuid.UUID, ordersID uuid.UUID) (models.Move, error) {

	var move models.Move
	err := appCfg.DB().EagerPreload("MTOServiceItems").
		Where("moves.id = ?", moveID).
		First(&move)

	if err != nil {
		appCfg.Logger().Error("failure encountered querying for move associated with orders", zap.Error(err))
		return models.Move{}, fmt.Errorf("failure encountered querying for move associated with orders, %s, id: %s", err.Error(), ordersID)
	}

	if move.Status != models.MoveStatusAPPROVALSREQUESTED {
		return models.Move{}, errors.Wrap(
			models.ErrInvalidTransition,
			"Cannot approve move with amended orders because the move status is not APPROVALS REQUESTED",
		)
	}

	var hasRequestedServiceItems bool
	for _, serviceItem := range move.MTOServiceItems {
		if serviceItem.Status == models.MTOServiceItemStatusSubmitted {
			hasRequestedServiceItems = true
			break
		}
	}

	if !hasRequestedServiceItems {
		approveErr := router.Approve(appCfg, &move)
		if approveErr != nil {
			return models.Move{}, approveErr
		}
	}

	return move, nil
}

func (router moveRouter) logMove(appCfg appconfig.AppConfig, move *models.Move) {
	appCfg.Logger().Info("Move log",
		zap.String("Move.ID", move.ID.String()),
		zap.String("Move.Locator", move.Locator),
		zap.String("Move.Status", string(move.Status)),
		zap.String("Move.OrdersID", move.OrdersID.String()),
	)
}
