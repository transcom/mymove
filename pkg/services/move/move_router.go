package move

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
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
func (router moveRouter) Submit(appCtx appcontext.AppContext, move *models.Move) error {
	router.logMove(appCtx, move)

	needsServicesCounseling, err := router.needsServiceCounseling(appCtx, move)
	if err != nil {
		appCtx.Logger().Error("failure determining if a move needs services counseling", zap.Error(err))
		return err
	}
	appCtx.Logger().Info("SUCCESS: Determining if move needs services counseling or not")

	if needsServicesCounseling {
		err = router.sendToServiceCounselor(appCtx, move)
		if err != nil {
			appCtx.Logger().Error("failure routing move to services counseling", zap.Error(err))
			return err
		}
		appCtx.Logger().Info("SUCCESS: Move sent to services counseling")
	} else if move.Orders.UploadedAmendedOrders != nil {
		appCtx.Logger().Info("Move has amended orders")
		transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			err = router.SendToOfficeUser(txnAppCtx, move)
			if err != nil {
				txnAppCtx.Logger().Error("failure routing move submission with amended orders", zap.Error(err))
				return err
			}
			// Let's get the orders for this move so we can wipe out the acknowledgement if it exists already (from a prior orders amendment process)
			var ordersForMove models.Order
			err = txnAppCtx.DB().Find(&ordersForMove, move.OrdersID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return apperror.NewNotFoundError(move.OrdersID, "looking for Order")
				default:
					return apperror.NewQueryError("Order", err, "")
				}
			}
			// Here we'll nil out the value (if it's set already) so that on the client-side we'll see view this change
			// in status as 'new orders' that need acknowledging by the TOO.
			// We shouldn't need more complicated logic here since we only hit this point from calling Submit().
			// Other circumstances like new MTOServiceItems will be calling SendToOfficeUser() directly.
			txnAppCtx.Logger().Info("Determining whether there is a preexisting orders acknowledgement")
			if ordersForMove.AmendedOrdersAcknowledgedAt != nil {
				txnAppCtx.Logger().Info("Move has a preexisting acknowledgement")
				ordersForMove.AmendedOrdersAcknowledgedAt = nil
				_, err = txnAppCtx.DB().ValidateAndSave(&ordersForMove)
				if err != nil {
					txnAppCtx.Logger().Error("failure resetting orders AmendedOrdersAcknowledgeAt field when routing move submission with amended orders ", zap.Error(err))
					return err
				}
				txnAppCtx.Logger().Info("Successfully reset orders acknowledgement")
			}
			return nil
		})
		if transactionError != nil {
			return transactionError
		}
		appCtx.Logger().Info("SUCCESS: Move with amended orders sent to office user / TOO queue")
	} else {
		err = router.sendNewMoveToOfficeUser(appCtx, move)
		if err != nil {
			appCtx.Logger().Error("failure routing move to office user / TOO queue", zap.Error(err))
			return err
		}
		appCtx.Logger().Info("SUCCESS: Move sent to office user / TOO queue")
	}

	appCtx.Logger().Info("SUCCESS: Move submitted and routed to the appropriate queue")
	return nil
}

func (router moveRouter) needsServiceCounseling(appCtx appcontext.AppContext, move *models.Move) (bool, error) {
	var orders models.Order
	err := appCtx.DB().Q().
		Where("orders.id = ?", move.OrdersID).
		First(&orders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			appCtx.Logger().Error("failure finding move", zap.Error(err))
			return false, apperror.NewNotFoundError(move.OrdersID, "looking for move.OrdersID")
		default:
			appCtx.Logger().Error("failure encountered querying for orders associated with the move", zap.Error(err))
			return false, apperror.NewQueryError("Order", err, fmt.Sprintf("failure encountered querying for orders associated with the move, %s, id: %s", err.Error(), move.ID))
		}
	}

	var originDutyLocation models.DutyLocation

	if orders.OriginDutyLocationID == nil || *orders.OriginDutyLocationID == uuid.Nil {
		return false, apperror.NewInvalidInputError(orders.ID, err, nil, "orders missing OriginDutyLocation")
	}

	originDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), *orders.OriginDutyLocationID)
	if err != nil {
		appCtx.Logger().Error("failure finding the origin duty station", zap.Error(err))
		return false, apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty station")
	}

	if move.ServiceCounselingCompletedAt != nil {
		return false, nil
	}

	return originDutyLocation.ProvidesServicesCounseling, nil
}

// sendToServiceCounselor makes the move available for a Service Counselor to review
func (router moveRouter) sendToServiceCounselor(appCtx appcontext.AppContext, move *models.Move) error {
	if move.Status == models.MoveStatusNeedsServiceCounseling {
		return nil
	}

	if move.Status != models.MoveStatusDRAFT {
		appCtx.Logger().Warn(fmt.Sprintf(
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
func (router moveRouter) sendNewMoveToOfficeUser(appCtx appcontext.AppContext, move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		appCtx.Logger().Warn(fmt.Sprintf(
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
			appCtx.Logger().Error("Failure submitting ppm", zap.Error(err))
			return err
		}
	}

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				appCtx.Logger().Error("Failure requesting reimbursement for ppm", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

// Approve makes the Move available to the Prime. The Prime cannot create
// Service Items unless the Move is approved.
func (router moveRouter) Approve(appCtx appcontext.AppContext, move *models.Move) error {
	router.logMove(appCtx, move)
	if router.alreadyApproved(move) {
		return nil
	}

	if currentStatusApprovable(*move) {
		move.Status = models.MoveStatusAPPROVED
		appCtx.Logger().Info("SUCCESS: Move approved")
		return nil
	}

	appCtx.Logger().Warn(fmt.Sprintf(
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

func currentStatusApprovable(move models.Move) bool {
	return statusSliceContains(validStatusesBeforeApproval, move.Status)
}

func approvable(move models.Move) bool {
	return moveHasReviewedServiceItems(move) &&
		moveHasAcknowledgedOrdersAmendment(move.Orders) &&
		moveHasAcknowledgedExcessWeightRisk(move) &&
		allSITExtensionsAreReviewed(move)
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

func moveHasAcknowledgedOrdersAmendment(order models.Order) bool {
	if order.UploadedAmendedOrdersID != nil && order.AmendedOrdersAcknowledgedAt == nil {
		return false
	}
	return true
}

func moveHasReviewedServiceItems(move models.Move) bool {
	for _, mtoServiceItem := range move.MTOServiceItems {
		if mtoServiceItem.Status == models.MTOServiceItemStatusSubmitted {
			return false
		}
	}

	return true
}

func moveHasAcknowledgedExcessWeightRisk(move models.Move) bool {
	// If the move hasn't been flagged for being at risk of excess weight, then
	// we don't need to check if the risk has been acknowledged.
	if move.ExcessWeightQualifiedAt == nil {
		return true
	}
	return move.ExcessWeightAcknowledgedAt != nil
}

func allSITExtensionsAreReviewed(move models.Move) bool {
	for _, shipment := range move.MTOShipments {
		for _, sitExtension := range shipment.SITExtensions {
			if sitExtension.Status == models.SITExtensionStatusPending {
				return false
			}
		}
	}

	return true
}

// SendToOfficeUser sets the moves status to
// "Approvals Requested", which indicates to the TOO that they have new
// service items to review.
func (router moveRouter) SendToOfficeUser(appCtx appcontext.AppContext, move *models.Move) error {
	router.logMove(appCtx, move)
	// Do nothing if it's already in the desired state
	if move.Status == models.MoveStatusAPPROVALSREQUESTED {
		return nil
	}
	if move.Status == models.MoveStatusCANCELED {
		errorMessage := fmt.Sprintf("The status for the move with ID %s can not be sent to 'Approvals Requested' if the status is cancelled.", move.ID)
		appCtx.Logger().Warn(errorMessage)

		return errors.Wrap(models.ErrInvalidTransition, errorMessage)
	}
	move.Status = models.MoveStatusAPPROVALSREQUESTED
	appCtx.Logger().Info("SUCCESS: Move sent to TOO to request approval")

	return nil
}

// Cancel cancels the Move and its associated PPMs
func (router moveRouter) Cancel(appCtx appcontext.AppContext, reason string, move *models.Move) error {
	router.logMove(appCtx, move)
	// We can cancel any move that isn't already complete.
	// TODO: What does complete mean? How do we determine when a move is complete?
	if move.Status == models.MoveStatusCANCELED {
		return errors.Wrap(models.ErrInvalidTransition, "cannot cancel a move that is already canceled")
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

	appCtx.Logger().Info("SUCCESS: Move Canceled")
	return nil

}

// CompleteServiceCounseling sets the move status to "Service Counseling Completed",
// which makes the move available to the TOO. This gets called when the Service
// Counselor is done reviewing the move and submits it.
func (router moveRouter) CompleteServiceCounseling(appCtx appcontext.AppContext, move *models.Move) error {
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

// ApproveOrRequestApproval routes the move appropriately based on whether or
// not the TOO has any tasks requiring their attention.
func (router moveRouter) ApproveOrRequestApproval(appCtx appcontext.AppContext, move models.Move) (*models.Move, error) {
	err := appCtx.DB().Q().EagerPreload("MTOServiceItems", "Orders", "MTOShipments.SITExtensions").Find(&move, move.ID)
	if err != nil {
		appCtx.Logger().Error("Failed to preload MTOServiceItems and Orders for Move", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(move.ID, "looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	if approvable(move) {
		err = router.Approve(appCtx, &move)
	} else {
		err = router.SendToOfficeUser(appCtx, &move)
	}

	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if e := handleError(move.ID, verrs, err); e != nil {
		return nil, e
	}

	return &move, nil
}

func handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}

func (router moveRouter) logMove(appCtx appcontext.AppContext, move *models.Move) {
	appCtx.Logger().Info("Move log",
		zap.String("Move.ID", move.ID.String()),
		zap.String("Move.Locator", move.Locator),
		zap.String("Move.Status", string(move.Status)),
		zap.String("Move.OrdersID", move.OrdersID.String()),
	)
}
