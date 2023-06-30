package mtoshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentRouter struct {
}

// NewShipmentRouter creates a new shipmentRouter service
func NewShipmentRouter() services.ShipmentRouter {
	return &shipmentRouter{}
}

// Submit is used to submit a shipment either at creation time, or when the customer submits their move. It's up to the
// caller to save the shipment.
func (router shipmentRouter) Submit(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusDraft && shipment.Status != "" {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusSubmitted,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusDraft},
		}
	}
	shipment.Status = models.MTOShipmentStatusSubmitted

	return nil
}

// Approve checks if a shipment can be approved, and if so, sets the appropriate status and date. It's up to the caller
// to save the shipment.
func (router shipmentRouter) Approve(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	// When a shipment is approved, service items automatically get created, but
	// service items can only be created if a Move's status is either Approved
	// or Approvals Requested, so check and fail early.
	move := shipment.MoveTaskOrder
	if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED && !(shipment.ShipmentType == models.MTOShipmentTypePPM && move.Status == models.MoveStatusNeedsServiceCounseling) {
		return apperror.NewConflictError(
			move.ID,
			fmt.Sprintf(
				"Cannot approve a shipment if the move status isn't %s or %s, or if it isn't a PPM shipment with a move status of %s. The current status for the move with ID %s is %s",
				models.MoveStatusAPPROVED,
				models.MoveStatusAPPROVALSREQUESTED,
				models.MoveStatusNeedsServiceCounseling,
				move.ID,
				move.Status,
			),
		)
	}

	if shipment.UsesExternalVendor {
		return apperror.NewConflictError(
			shipment.ID,
			fmt.Sprintf("cannot approve a shipment if it uses an external vendor. The current status for the shipment with ID %s is %s", shipment.ID, shipment.Status),
		)
	}

	if router.approvable(shipment) {
		shipment.Status = models.MTOShipmentStatusApproved
		approvedDate := time.Now()
		shipment.ApprovedDate = &approvedDate

		return nil
	}

	return ConflictStatusError{
		id:                        shipment.ID,
		transitionFromStatus:      shipment.Status,
		transitionToStatus:        models.MTOShipmentStatusApproved,
		transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusSubmitted, models.MTOShipmentStatusDiversionRequested},
	}
}

// RequestCancellation is called when the TOO has requested that the Prime cancel the shipment.
func (router shipmentRouter) RequestCancellation(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusApproved {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusCancellationRequested,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusApproved},
		}
	}
	shipment.Status = models.MTOShipmentStatusCancellationRequested

	return nil
}

// Cancel cancels the shipment
func (router shipmentRouter) Cancel(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusCancellationRequested {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusCanceled,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusCancellationRequested},
		}
	}

	shipment.Status = models.MTOShipmentStatusCanceled

	return nil
}

// Reject rejects the shipment
func (router shipmentRouter) Reject(_ appcontext.AppContext, shipment *models.MTOShipment, reason *string) error {
	if shipment.Status != models.MTOShipmentStatusSubmitted {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusRejected,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusSubmitted},
		}
	}

	shipment.Status = models.MTOShipmentStatusRejected
	shipment.RejectionReason = reason

	return nil
}

// RequestDiversion is called when the TOO has requested that the Prime divert the shipment.
func (router shipmentRouter) RequestDiversion(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusApproved {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusDiversionRequested,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusApproved},
		}
	}
	shipment.Status = models.MTOShipmentStatusDiversionRequested

	return nil
}

// ApproveDiversion is called when the TOO is approving a shipment that the Prime has marked as being diverted.
func (router shipmentRouter) ApproveDiversion(_ appcontext.AppContext, shipment *models.MTOShipment) error {
	if !shipment.Diversion {
		return apperror.NewConflictError(
			shipment.ID,
			fmt.Sprintf("Cannot approve the diversion because the shipment with id %s has the Diversion field set to false.", shipment.ID),
		)
	}

	if shipment.UsesExternalVendor {
		return apperror.NewConflictError(
			shipment.ID,
			fmt.Sprintf("shipmentRouter: cannot approve the diversion because the shipment with id %s has the UsesExternalVendor field set to true.", shipment.ID),
		)
	}

	if shipment.Status != models.MTOShipmentStatusSubmitted {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusApproved,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusSubmitted},
		}
	}
	shipment.Status = models.MTOShipmentStatusApproved

	return nil
}

func (router shipmentRouter) approvable(shipment *models.MTOShipment) bool {
	// first check if the status is in list
	isApprovable := statusSliceContains(validStatusesBeforeApproval, shipment.Status)

	// then check special case for diversion requested status
	if shipment.Status == models.MTOShipmentStatusDiversionRequested {
		// a shipment is considered diverted or part of a diversion if
		// the prime sets the Diversion field to true
		isApprovable = shipment.Diversion
	}

	return isApprovable
}

func statusSliceContains(statusSlice []models.MTOShipmentStatus, status models.MTOShipmentStatus) bool {
	for _, validStatus := range statusSlice {
		if status == validStatus {
			return true
		}
	}
	return false
}

var validStatusesBeforeApproval = []models.MTOShipmentStatus{
	models.MTOShipmentStatusSubmitted,
	models.MTOShipmentStatusDiversionRequested,
}
