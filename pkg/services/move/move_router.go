package move

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveRouter struct {
}

// NewMoveRouter creates a new moveRouter service
func NewMoveRouter() services.MoveRouter {
	return &moveRouter{}
}

// Submit is called when the customer submits amended orders or submits their move. It determines whether
// to send the move to Service Counseling or directly to the TOO. If it goes to
// Service Counseling, its status becomes "Needs Service Counseling", otherwise,
// "Submitted". A signed certification should be passed in when submitting a move, but not when submitting
// amended orders.
func (router moveRouter) Submit(appCtx appcontext.AppContext, move *models.Move, newSignedCertification *models.SignedCertification) error {
	router.logMove(appCtx, move)
	var verrs *validate.Errors
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		// when a move is submitted and needs to be routed to a service counselor, the move status is updated,
		// ppm shipment statuses are updated (both the ppm shipment and the parent mtoShipment),
		// and a signedCertification is created. sendToServiceCounselor handles the move and shipment status changes
		needsServicesCounseling, err := router.needsServiceCounseling(appCtx, move)
		if err != nil {
			appCtx.Logger().Error("failure determining if a move needs services counseling", zap.Error(err))
			return err
		}
		if needsServicesCounseling {
			err = router.sendToServiceCounselor(txnAppCtx, move)
			if err != nil {
				appCtx.Logger().Error("failure routing move to services counseling", zap.Error(err))
				return err
			}

			appCtx.Logger().Info("SUCCESS: Move sent to services counselor")
		} else {
			err = router.sendNewMoveToOfficeUser(txnAppCtx, move)
			if err != nil {
				txnAppCtx.Logger().Error("failure routing move to office user", zap.Error(err))
				return err
			}
			appCtx.Logger().Info("SUCCESS: Move sent to office user")
		}

		if newSignedCertification == nil {
			msg := "signedCertification is required"
			appCtx.Logger().Error(msg, zap.Error(err))
			return apperror.NewInvalidInputError(move.ID, err, verrs, msg)
		}
		verrs, err = txnAppCtx.DB().ValidateAndCreate(newSignedCertification)
		if err != nil || verrs.HasAny() {
			txnAppCtx.Logger().Error("error saving signed certification: %w", zap.Error(err))
		}
		return err
	})

	if transactionError != nil {
		return transactionError
	}
	appCtx.Logger().Info("SUCCESS: Move submitted and routed to the appropriate queue")
	return nil
}

func (router moveRouter) RouteAfterAmendingOrders(appCtx appcontext.AppContext, move *models.Move) error {
	return appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		needsServicesCounseling, err := router.needsServiceCounseling(appCtx, move)
		if err != nil {
			appCtx.Logger().Error("failure determining if a move needs services counseling", zap.Error(err))
			return err
		}
		if needsServicesCounseling {
			err = router.sendToServiceCounselor(txnAppCtx, move)
			if err != nil {
				appCtx.Logger().Error("failure routing move to services counseling", zap.Error(err))
				return err
			}
			appCtx.Logger().Info("SUCCESS: Move sent to services counselor")
		} else {
			err := router.SendToOfficeUser(txnAppCtx, move)
			if err != nil {
				txnAppCtx.Logger().Error("failure routing move to office user", zap.Error(err))
				return err
			}
			appCtx.Logger().Info("SUCCESS: Move sent to office user")

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
		}
		return nil
	})
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
		appCtx.Logger().Error("failure finding the origin duty location", zap.Error(err))
		return false, apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
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

	// if it's a PPMShipment update both the mto and ppm shipment level statuses
	for i := range move.MTOShipments {
		if move.MTOShipments[i].ShipmentType == models.MTOShipmentTypePPM {
			move.MTOShipments[i].Status = models.MTOShipmentStatusSubmitted
			move.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusSubmitted

			if verrs, err := appCtx.DB().ValidateAndUpdate(&move.MTOShipments[i]); verrs.HasAny() || err != nil {
				msg := "failure saving shipment when routing move submission"
				appCtx.Logger().Error(msg, zap.Error(err))
				return apperror.NewInvalidInputError(move.MTOShipments[i].ID, err, verrs, msg)
			}
			if verrs, err := appCtx.DB().ValidateAndUpdate(move.MTOShipments[i].PPMShipment); verrs.HasAny() || err != nil {
				msg := "failure saving PPM shipment when routing move submission"
				appCtx.Logger().Error(msg, zap.Error(err))
				return apperror.NewInvalidInputError(move.MTOShipments[i].PPMShipment.ID, err, verrs, msg)
			}
		}
	}

	if verrs, err := appCtx.DB().ValidateAndSave(move); verrs.HasAny() || err != nil {
		msg := "failure saving move when routing move submission"
		appCtx.Logger().Error(msg, zap.Error(err))
		return apperror.NewInvalidInputError(move.ID, err, verrs, msg)
	}

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

	// if it's a PPMShipment update both the mto and ppm shipment level statuses
	for i := range move.MTOShipments {
		if move.MTOShipments[i].ShipmentType == models.MTOShipmentTypePPM {
			move.MTOShipments[i].Status = models.MTOShipmentStatusSubmitted
			move.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusSubmitted

			if verrs, err := appCtx.DB().ValidateAndUpdate(&move.MTOShipments[i]); verrs.HasAny() || err != nil {
				msg := "failure saving shipment when routing move submission"
				appCtx.Logger().Error(msg, zap.Error(err))
				return apperror.NewInvalidInputError(move.MTOShipments[i].ID, err, verrs, msg)
			}

			if verrs, err := appCtx.DB().ValidateAndUpdate(move.MTOShipments[i].PPMShipment); verrs.HasAny() || err != nil {
				msg := "failure saving PPM shipment when routing move submission"
				appCtx.Logger().Error(msg, zap.Error(err))
				return apperror.NewInvalidInputError(move.MTOShipments[i].PPMShipment.ID, err, verrs, msg)
			}
		}
	}

	if verrs, err := appCtx.DB().ValidateAndSave(move); verrs.HasAny() || err != nil {
		msg := "failure saving move when routing move submission"
		appCtx.Logger().Error(msg, zap.Error(err))
		return apperror.NewInvalidInputError(move.ID, err, verrs, msg)
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
		allSITExtensionsAreReviewed(move) &&
		allShipmentAddressUpdatesAreReviewed(move) &&
		allShipmentsAreApproved(move)
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
		for _, sitDurationUpdate := range shipment.SITDurationUpdates {
			if sitDurationUpdate.Status == models.SITExtensionStatusPending {
				return false
			}
		}
	}

	return true
}

func allShipmentsAreApproved(move models.Move) bool {
	for _, shipment := range move.MTOShipments {
		if shipment.Status == models.MTOShipmentStatusSubmitted {
			return false
		}
	}

	return true
}

func allShipmentAddressUpdatesAreReviewed(move models.Move) bool {
	for _, shipment := range move.MTOShipments {
		if shipment.DeliveryAddressUpdate != nil && shipment.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusRequested {
			return false
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
	now := time.Now()
	move.ApprovalsRequestedAt = &now

	if verrs, err := appCtx.DB().ValidateAndSave(move); verrs.HasAny() || err != nil {
		msg := "failure saving move when routing move submission"
		appCtx.Logger().Error(msg, zap.Error(err))
		return apperror.NewInvalidInputError(move.ID, err, verrs, msg)
	}
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
	for i := range move.MTOShipments {
		err := move.MTOShipments[i].PPMShipment.CancelShipment()
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

// CompleteServiceCounseling sets the move status to:
//   - "Service Counseling Completed" if a non-PPM move
//   - "Approved" if a PPM-only move
//
// This makes the move available to the TOO.  This gets called when the Service Counselor is done
// reviewing the move and submits it.
func (router moveRouter) CompleteServiceCounseling(_ appcontext.AppContext, move *models.Move) error {
	// Verify shipments are present.
	if len(move.MTOShipments) == 0 {
		return apperror.NewConflictError(move.ID, "No shipments associated with move")
	}

	// Verify the shipment's existing status.
	if move.Status != models.MoveStatusNeedsServiceCounseling {
		return apperror.NewConflictError(move.ID, fmt.Sprintf("The status for the Move with ID %s can only be set to 'Service Counseling Completed' from the 'Needs Service Counseling' status, but its current status is %s.", move.ID, move.Status))
	}

	// Examine shipments for valid state and how to transition.
	ppmOnlyMove := true
	for _, s := range move.MTOShipments {
		if s.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && s.StorageFacilityID == nil {
			return apperror.NewConflictError(s.ID, "NTS-release shipment must include facility info")
		}
		if s.ShipmentType != models.MTOShipmentTypePPM {
			ppmOnlyMove = false
		}
	}

	// Set target state based on associated shipments.
	targetState := models.MoveStatusServiceCounselingCompleted
	if ppmOnlyMove {
		targetState = models.MoveStatusAPPROVED
	}

	now := time.Now()
	move.ServiceCounselingCompletedAt = &now
	move.Status = targetState

	return nil
}

// ApproveOrRequestApproval routes the move appropriately based on whether or
// not the TOO has any tasks requiring their attention.
func (router moveRouter) ApproveOrRequestApproval(appCtx appcontext.AppContext, move models.Move) (*models.Move, error) {
	err := appCtx.DB().Q().EagerPreload("MTOServiceItems", "Orders.ServiceMember", "Orders.NewDutyLocation.Address", "MTOShipments.SITDurationUpdates", "MTOShipments.DeliveryAddressUpdate").Find(&move, move.ID)
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
