package weightticket

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// weightTicketFetcher is the concrete implementation of the services.WeightTicketFetcher interface
type weightTicketFetcher struct{}

// NewWeightTicketFetcher creates a new struct
func NewWeightTicketFetcher() services.WeightTicketFetcher {
	return &weightTicketFetcher{}
}

// GetWeightTicket fetches a WeightTicket by ID, excluding deleted weight tickets. The returned weight ticket will
// include uploads, without any deleted uploads.
func (f *weightTicketFetcher) GetWeightTicket(appCtx appcontext.AppContext, weightTicketID uuid.UUID) (*models.WeightTicket, error) {
	var weightTicket models.WeightTicket

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
		).
		Find(&weightTicket, weightTicketID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(weightTicketID, "while looking for WeightTicket")
		default:
			return nil, apperror.NewQueryError("WeightTicket", err, "unable to find WeightTicket")
		}
	}

	if appCtx.Session().IsMilApp() && weightTicket.EmptyDocument.ServiceMemberID != appCtx.Session().ServiceMemberID {
		return nil, apperror.NewForbiddenError("not authorized to access weight ticket")
	}

	weightTicket.EmptyDocument.UserUploads = weightTicket.EmptyDocument.UserUploads.FilterDeleted()
	weightTicket.FullDocument.UserUploads = weightTicket.FullDocument.UserUploads.FilterDeleted()
	weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightTicket.ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()

	return &weightTicket, nil
}

// ListWeightTickets fetches all WeightTickets for a given move, excluding deleted weight tickets. The returned weight
// tickets will include uploads, without any deleted uploads.
func (f *weightTicketFetcher) ListWeightTickets(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (models.WeightTickets, error) {
	var weightTickets models.WeightTickets

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
		).
		Where("ppm_shipment_id = ?", ppmShipmentID).
		All(&weightTickets)

	if err != nil {
		return nil, apperror.NewQueryError("WeightTicket", err, "unable to find WeightTickets")
	}

	if weightTickets == nil {
		return weightTickets, nil
	}

	if appCtx.Session().IsMilApp() && weightTickets[0].EmptyDocument.ServiceMemberID != appCtx.Session().ServiceMemberID {
		return nil, apperror.NewForbiddenError("not authorized to access weight tickets")
	}

	for i := range weightTickets {
		weightTickets[i].EmptyDocument.UserUploads = weightTickets[i].EmptyDocument.UserUploads.FilterDeleted()
		weightTickets[i].FullDocument.UserUploads = weightTickets[i].FullDocument.UserUploads.FilterDeleted()
		weightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads = weightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()
	}

	return weightTickets, nil
}
