package ppmshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentRouter is the concrete struct implementing the services.PPMShipmentRouter interface
type ppmShipmentRouter struct {
	services.ShipmentRouter
}

// NewPPMShipmentRouter creates a new ppmShipmentRouter
func NewPPMShipmentRouter(shipmentRouter services.ShipmentRouter) services.PPMShipmentRouter {
	return &ppmShipmentRouter{
		ShipmentRouter: shipmentRouter,
	}
}

// SetToDraft sets the PPM shipment to the DRAFT status
func (p *ppmShipmentRouter) SetToDraft(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	if ppmShipment.Status != "" {
		return apperror.NewConflictError(
			ppmShipment.ID,
			fmt.Sprintf("PPM shipment can't be set to %s because it's not new.", models.PPMShipmentStatusDraft),
		)
	}

	ppmShipment.Status = models.PPMShipmentStatusDraft

	// TODO: this should be done using the shipment router, but it currently doesn't have a way of setting this.
	ppmShipment.Shipment.Status = models.MTOShipmentStatusDraft

	return nil
}

// Submit sets the PPM shipment to the SUBMITTED status
func (p *ppmShipmentRouter) Submit(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	if ppmShipment.Status != "" && ppmShipment.Status != models.PPMShipmentStatusDraft {
		return apperror.NewConflictError(
			ppmShipment.ID,
			fmt.Sprintf(
				"PPM shipment can't be set to %s because it's not new or in the %s status.",
				models.PPMShipmentStatusSubmitted,
				models.PPMShipmentStatusDraft,
			),
		)
	}

	err := p.ShipmentRouter.Submit(appCtx, &ppmShipment.Shipment)

	if err != nil {
		return err
	}

	ppmShipment.Status = models.PPMShipmentStatusSubmitted

	return nil
}

// SendToCustomer sets the PPM shipment to the WAITING_ON_CUSTOMER status
func (p *ppmShipmentRouter) SendToCustomer(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	if ppmShipment.Status != models.PPMShipmentStatusSubmitted && ppmShipment.Status != models.PPMShipmentStatusNeedsCloseout {
		return apperror.NewConflictError(
			ppmShipment.ID,
			fmt.Sprintf(
				"PPM shipment can't be set to %s because it's not in a %s or %s status.",
				models.PPMShipmentStatusWaitingOnCustomer,
				models.PPMShipmentStatusSubmitted,
				models.PPMShipmentStatusNeedsCloseout,
			),
		)
	}

	if ppmShipment.Shipment.Status != models.MTOShipmentStatusApproved {
		err := p.ShipmentRouter.Approve(appCtx, &ppmShipment.Shipment)

		if err != nil {
			return err
		}
	}

	ppmShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

	if ppmShipment.ApprovedAt == nil {
		ppmShipment.ApprovedAt = models.TimePointer(*ppmShipment.Shipment.ApprovedDate)
	}

	return nil
}

// SubmitCloseOutDocumentation sets the PPM shipment to the NEEDS_CLOSEOUT status
func (p *ppmShipmentRouter) SubmitCloseOutDocumentation(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	if ppmShipment.Status != models.PPMShipmentStatusWaitingOnCustomer {
		return apperror.NewConflictError(
			ppmShipment.ID,
			fmt.Sprintf(
				"PPM shipment can't be set to %s because it's not in the %s status.",
				models.PPMShipmentStatusNeedsCloseout,
				models.PPMShipmentStatusWaitingOnCustomer,
			),
		)
	}

	ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout

	if ppmShipment.SubmittedAt == nil {
		ppmShipment.SubmittedAt = models.TimePointer(time.Now())
	}

	return nil
}

// SubmitReviewedDocuments sets the PPM shipment status to the CLOSEOUT_COMPLETE if all docs approved otherwise WAITING_ON_CUSTOMER
func (p *ppmShipmentRouter) SubmitReviewedDocuments(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	if ppmShipment.Status != models.PPMShipmentStatusNeedsCloseout {
		return apperror.NewConflictError(
			ppmShipment.ID,
			fmt.Sprintf(
				"PPM shipment documents cannot be submitted because it's not in the %s status.",
				models.PPMShipmentStatusNeedsCloseout,
			),
		)
	}

	hasRejectedDocuments := false
	if len(ppmShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			if weightTicket.Status != nil && *weightTicket.Status == models.PPMDocumentStatusRejected {
				hasRejectedDocuments = true
				break
			}
		}
	}

	if len(ppmShipment.ProgearWeightTickets) >= 1 && !hasRejectedDocuments {
		for _, progear := range ppmShipment.ProgearWeightTickets {
			if progear.Status != nil && *progear.Status == models.PPMDocumentStatusRejected {
				hasRejectedDocuments = true
				break
			}
		}
	}

	if len(ppmShipment.MovingExpenses) >= 1 && !hasRejectedDocuments {
		for _, movingExpenses := range ppmShipment.MovingExpenses {
			if movingExpenses.Status != nil && *movingExpenses.Status == models.PPMDocumentStatusRejected {
				hasRejectedDocuments = true
				break
			}
		}
	}

	if hasRejectedDocuments {
		ppmShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
	} else {
		ppmShipment.Status = models.PPMShipmentStatusCloseoutComplete
	}

	return nil
}
