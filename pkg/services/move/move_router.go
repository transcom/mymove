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
	transportationOfficesFetcher services.TransportationOfficesFetcher
}

// NewMoveRouter creates a new moveRouter service
func NewMoveRouter(transportationOfficeFetcher services.TransportationOfficesFetcher) services.MoveRouter {
	return &moveRouter{
		transportationOfficesFetcher: transportationOfficeFetcher,
	}
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

	if move.IsPPMOnly() {
		return true, nil
	}

	return originDutyLocation.ProvidesServicesCounseling, nil
}

// sendToServiceCounselor makes the move available for a Service Counselor to review
func (router moveRouter) sendToServiceCounselor(appCtx appcontext.AppContext, move *models.Move) error {
	var orders models.Order
	var originDutyLocation models.DutyLocation
	err := appCtx.DB().Q().
		Where("orders.id = ?", move.OrdersID).
		First(&orders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			appCtx.Logger().Error("failure finding move", zap.Error(err))
			return apperror.NewNotFoundError(move.OrdersID, "looking for move.OrdersID")
		default:
			appCtx.Logger().Error("failure encountered querying for orders associated with the move", zap.Error(err))
			return apperror.NewQueryError("Order", err, fmt.Sprintf("failure encountered querying for orders associated with the move, %s, id: %s", err.Error(), move.ID))
		}
	}

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

	isCivilian := orders.Grade != nil && *orders.Grade == models.ServiceMemberGradeCIVILIANEMPLOYEE
	move.Status = models.MoveStatusNeedsServiceCounseling
	now := time.Now()
	move.SubmittedAt = &now
	if orders.OriginDutyLocationID == nil || *orders.OriginDutyLocationID == uuid.Nil {
		return apperror.NewInvalidInputError(orders.ID, err, nil, "orders missing OriginDutyLocation")
	}

	originDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), *orders.OriginDutyLocationID)
	if err != nil {
		appCtx.Logger().Error("failure finding the origin duty location", zap.Error(err))
		return apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
	}
	orders.OriginDutyLocation = &originDutyLocation
	for i := range move.MTOShipments {
		// if it's a PPMShipment update both the mto and ppm shipment level statuses
		if move.MTOShipments[i].ShipmentType == models.MTOShipmentTypePPM {
			move.MTOShipments[i].Status = models.MTOShipmentStatusSubmitted
			move.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusSubmitted
			// actual expense reimbursement is always true for civilian moves
			move.MTOShipments[i].PPMShipment.IsActualExpenseReimbursement = models.BoolPointer(isCivilian)
			if move.IsPPMOnly() && !orders.OriginDutyLocation.ProvidesServicesCounseling {
				closestCounselingOffice, err := router.transportationOfficesFetcher.FindCounselingOfficeForPrimeCounseled(appCtx, *move.Orders.OriginDutyLocationID, move.Orders.ServiceMemberID)
				if err != nil {
					msg := "Failure setting PPM counseling office to closest service counseling office"
					appCtx.Logger().Error(msg, zap.Error(err))
					return apperror.NewQueryError("Closest Counseling Office", err, "Failed to find counseling office that provides counseling")
				}
				move.CounselingOfficeID = &closestCounselingOffice.ID
			}

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

		// update status for boat or mobile home shipment
		if move.MTOShipments[i].ShipmentType == models.MTOShipmentTypeBoatHaulAway ||
			move.MTOShipments[i].ShipmentType == models.MTOShipmentTypeBoatTowAway ||
			move.MTOShipments[i].ShipmentType == models.MTOShipmentTypeMobileHome {
			move.MTOShipments[i].Status = models.MTOShipmentStatusSubmitted

			if verrs, err := appCtx.DB().ValidateAndUpdate(&move.MTOShipments[i]); verrs.HasAny() || err != nil {
				msg := "failure saving parent MTO shipment object for boat/mobile home shipment when routing move submission"
				appCtx.Logger().Error(msg, zap.Error(err))
				return apperror.NewInvalidInputError(move.MTOShipments[i].ID, err, verrs, msg)
			}

			if move.MTOShipments[i].BoatShipment != nil {
				if verrs, err := appCtx.DB().ValidateAndUpdate(move.MTOShipments[i].BoatShipment); verrs.HasAny() || err != nil {
					msg := "failure saving boat shipment when routing move submission"
					appCtx.Logger().Error(msg, zap.Error(err))
					return apperror.NewInvalidInputError(move.MTOShipments[i].ID, err, verrs, msg)
				}
			}

			if move.MTOShipments[i].MobileHome != nil {
				if verrs, err := appCtx.DB().ValidateAndUpdate(move.MTOShipments[i].MobileHome); verrs.HasAny() || err != nil {
					msg := "failure saving mobile home shipment when routing move submission"
					appCtx.Logger().Error(msg, zap.Error(err))
					return apperror.NewInvalidInputError(move.MTOShipments[i].ID, err, verrs, msg)
				}
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
	if move == nil {
		return errors.New("cannot approve nil move")
	}

	router.logMove(appCtx, move)
	if router.alreadyApproved(move) && router.noAssignedTOOs(move) {
		return nil
	}

	if currentStatusApprovable(*move) {
		move.Status = models.MoveStatusAPPROVED
		now := time.Now()
		move.ApprovedAt = &now
		appCtx.Logger().Info("SUCCESS: Move approved")
		// if a move is approvable, we can clear any assigned office users, if any
		move.TOOTaskOrderAssignedID = nil
		move.TOODestinationAssignedID = nil
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

func (router moveRouter) noAssignedTOOs(move *models.Move) bool {
	return move.TOOTaskOrderAssignedID == nil && move.TOODestinationAssignedID == nil
}

func currentStatusApprovable(move models.Move) bool {
	return statusSliceContains(validStatusesBeforeApproval, move.Status)
}

func approvable(move models.Move) bool {
	return moveHasReviewedServiceItems(move) &&
		moveHasAcknowledgedOrdersAmendment(move.Orders) &&
		moveHasAcknowledgedExcessWeightRisk(move) &&
		moveHasAcknowledgedUBExcessWeightRisk(move) &&
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

func moveHasAcknowledgedUBExcessWeightRisk(move models.Move) bool {
	if move.ExcessUnaccompaniedBaggageWeightQualifiedAt == nil {
		return true
	}
	return move.ExcessUnaccompaniedBaggageWeightAcknowledgedAt != nil
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
		// ignores deleted shipments
		if (shipment.Status == models.MTOShipmentStatusSubmitted || shipment.Status == models.MTOShipmentStatusApprovalsRequested) && shipment.DeletedAt == nil {
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

func (router moveRouter) UpdateShipmentStatusToApprovalsRequested(appCtx appcontext.AppContext, shipment models.MTOShipment) (*models.MTOShipment, error) {
	if shipment.Status == models.MTOShipmentStatusApprovalsRequested {
		return nil, nil
	}
	if shipment.Status == models.MTOShipmentStatusCanceled || shipment.Status == models.MTOShipmentStatusTerminatedForCause {
		errorMessage := fmt.Sprintf("The status for the shipment with ID %s can not be sent to 'Approvals Requested' if the status is %s.", shipment.ID, shipment.Status)
		appCtx.Logger().Warn(errorMessage)

		return nil, errors.Wrap(models.ErrInvalidTransition, errorMessage)
	}
	shipment.Status = models.MTOShipmentStatusApprovalsRequested
	if verrs, err := appCtx.DB().ValidateAndSave(&shipment); verrs.HasAny() || err != nil {
		msg := "failure saving shipment"
		appCtx.Logger().Error(msg, zap.Error(err))
		return nil, apperror.NewInvalidInputError(shipment.ID, err, verrs, msg)
	}
	appCtx.Logger().Info("SUCCESS: Shipment status updated to Approvals Requested")

	return &shipment, nil
}

// Cancel cancels the Move and its associated PPMs
func (router moveRouter) Cancel(appCtx appcontext.AppContext, move *models.Move) error {
	moveDelta := move
	moveDelta.Status = models.MoveStatusCANCELED

	// get all shipments in move for cancellation
	var shipments []models.MTOShipment
	err := appCtx.DB().EagerPreload("Status", "PPMShipment", "PPMShipment.Status").Where("mto_shipments.move_id = $1", move.ID).All(&shipments)
	if err != nil {
		return apperror.NewNotFoundError(move.ID, "while looking for shipments")
	}

	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, shipment := range shipments {
			shipmentDelta := shipment
			shipmentDelta.Status = models.MTOShipmentStatusCanceled

			if shipment.PPMShipment != nil {
				var ppmshipment models.PPMShipment
				qerr := appCtx.DB().Where("id = ?", shipment.PPMShipment.ID).First(&ppmshipment)
				if qerr != nil {
					return apperror.NewNotFoundError(ppmshipment.ID, "while looking for ppm shipment")
				}

				ppmshipment.Status = models.PPMShipmentStatusCanceled

				verrs, err := txnAppCtx.DB().ValidateAndUpdate(&ppmshipment)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(shipment.ID, err, verrs, "Validation errors found while setting shipment status")
				} else if err != nil {
					return apperror.NewQueryError("PPM Shipment", err, "Failed to update status for ppm shipment")
				}
			}

			verrs, err := txnAppCtx.DB().ValidateAndUpdate(&shipmentDelta)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(shipment.ID, err, verrs, "Validation errors found while setting shipment status")
			} else if err != nil {
				return apperror.NewQueryError("Shipment", err, "Failed to update status for shipment")
			}
		}

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(moveDelta)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while setting move status")
		} else if err != nil {
			return apperror.NewQueryError("Move", err, "Failed to update status for move")
		}

		return nil
	})

	if txnErr != nil {
		return txnErr
	}

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
		if s.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS && s.StorageFacilityID == nil {
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
	err := appCtx.DB().Q().
		EagerPreload(
			"MTOServiceItems.ReService",
			"MTOShipments.SITDurationUpdates",
			"MTOShipments.DeliveryAddressUpdate",
			"Orders.ServiceMember",
			"Orders.NewDutyLocation.Address",
			"Orders.UploadedAmendedOrders",
		).
		Find(&move, move.ID)
	if err != nil {
		appCtx.Logger().Error("failed to preload data prior when routing move in ApproveOrRequestApproval", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(move.ID, "looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	// if a TOO is assigned to the move, check if we should clear it
	// this returns the same move with the TOO fields updated (or not)
	// !IMPORTANT - if any TOO actions are added, please also update this function
	if move.TOOTaskOrderAssignedID != nil || move.TOODestinationAssignedID != nil {
		updatedMove, err := models.ClearTOOAssignments(&move)
		if err != nil {
			return nil, err
		}
		move = *updatedMove
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
