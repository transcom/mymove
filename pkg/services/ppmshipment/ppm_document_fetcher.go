package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmDocumentFetcher is the concrete implementation of the services.PPMDocumentFetcher interface
type ppmDocumentFetcher struct{}

// NewPPMDocumentFetcher creates a new struct
func NewPPMDocumentFetcher() services.PPMDocumentFetcher {
	return &ppmDocumentFetcher{}
}

// GetPPMDocuments returns all documents associated with a PPM shipment.
func (f *ppmDocumentFetcher) GetPPMDocuments(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.PPMDocuments, error) {
	var documents models.PPMDocuments

	err := appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.WeightTicket{})).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = weight_tickets.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.WeightTickets)

	if err != nil {
		return nil, apperror.NewQueryError("WeightTicket", err, "unable to search for WeightTickets")
	}

	for i := range documents.WeightTickets {
		documents.WeightTickets[i].EmptyDocument.UserUploads = documents.WeightTickets[i].EmptyDocument.UserUploads.FilterDeleted()
		documents.WeightTickets[i].FullDocument.UserUploads = documents.WeightTickets[i].FullDocument.UserUploads.FilterDeleted()
		documents.WeightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads = documents.WeightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()
	}

	err = appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.ProgearWeightTicket{})).
		EagerPreload(
			"Document.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = progear_weight_tickets.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.ProgearExpenses)

	if err != nil {
		return nil, apperror.NewQueryError("ProgearWeightTicket", err, "unable to search for ProgearWeightTickets")
	}

	for i := range documents.ProgearExpenses {
		documents.ProgearExpenses[i].Document.UserUploads = documents.ProgearExpenses[i].Document.UserUploads.FilterDeleted()
	}

	err = appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.MovingExpense{})).
		EagerPreload(
			"Document.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = moving_expenses.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.MovingExpenses)

	if err != nil {
		return nil, apperror.NewQueryError("MovingExpense", err, "unable to search for MovingExpenses")
	}

	for i := range documents.MovingExpenses {
		documents.MovingExpenses[i].Document.UserUploads = documents.MovingExpenses[i].Document.UserUploads.FilterDeleted()
	}

	return &documents, nil
}
