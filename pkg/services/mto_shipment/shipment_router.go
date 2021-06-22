package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentRouter struct {
	db *pop.Connection
}

// NewShipmentRouter creates a new shipmentRouter service
func NewShipmentRouter(db *pop.Connection) services.ShipmentRouter {
	return &shipmentRouter{db}
}

// Submit is used to submit a shipment at the time the customer submits
// their move.
func (router shipmentRouter) Submit(shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusDraft {
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

// Approve is called when the TOO approves the shipment.
func (router shipmentRouter) Approve(shipment *models.MTOShipment) error {
	// When a shipment is approved, service items automatically get created, but
	// service items can only be created if a Move's status is either Approved
	// or Approvals Requested, so check and fail early.
	move := shipment.MoveTaskOrder
	if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED {
		return services.NewConflictError(
			move.ID,
			fmt.Sprintf("Cannot approve a shipment if the move isn't approved. The current status for the move with ID %s is %s", move.ID, move.Status),
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
func (router shipmentRouter) RequestCancellation(shipment *models.MTOShipment) error {
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
func (router shipmentRouter) Cancel(shipment *models.MTOShipment) error {
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
func (router shipmentRouter) Reject(shipment *models.MTOShipment, reason *string) error {
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
func (router shipmentRouter) RequestDiversion(shipment *models.MTOShipment) error {
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
func (router shipmentRouter) ApproveDiversion(shipment *models.MTOShipment) error {
	if shipment.Status != models.MTOShipmentStatusDiversionRequested {
		return ConflictStatusError{
			id:                        shipment.ID,
			transitionFromStatus:      shipment.Status,
			transitionToStatus:        models.MTOShipmentStatusApproved,
			transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusDiversionRequested},
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
