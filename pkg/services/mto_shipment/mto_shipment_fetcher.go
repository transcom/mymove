package mtoshipment

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
)

type mtoShipmentFetcher struct {
}

// NewMTOShipmentFetcher creates a new MTOShipmentFetcher struct that supports ListMTOShipments
func NewMTOShipmentFetcher() services.MTOShipmentFetcher {
	return &mtoShipmentFetcher{}
}

func (f mtoShipmentFetcher) ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) ([]models.MTOShipment, error) {
	var move models.Move
	err := appCtx.DB().Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "move not found")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	var shipments []models.MTOShipment
	err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"MTOServiceItems.ReService",
			"MTOAgents",
			"PickupAddress",
			"SecondaryPickupAddress",
			"DestinationAddress",
			"SecondaryDeliveryAddress",
			"MTOServiceItems.Dimensions",
			"PPMShipment.WeightTickets.EmptyDocument.UserUploads.Upload",
			// Can't EagerPreload "Reweigh" due to a Pop bug (see below)
			// "Reweigh",
			"SITExtensions",
			"StorageFacility.Address",
		). // Right now no use case for showing deleted shipments.
		Where("move_id = ?", moveID).
		Order("uses_external_vendor asc").
		Order("created_at asc").
		All(&shipments)

	if err != nil {
		return nil, err
	}

	// Due to a Pop bug, we cannot EagerPreload "Reweigh" or "PPMShipment" likely because it is a pointer and
	// a "has_one" field.  This seems similar to other EagerPreload issues we've found (and
	// sometimes fixed): https://github.com/gobuffalo/pop/issues?q=author%3Areggieriser
	for i := range shipments {
		loadErr := appCtx.DB().Load(&shipments[i], "Reweigh")
		if loadErr != nil {
			return nil, err
		}

		if shipments[i].ShipmentType == models.MTOShipmentTypePPM {
			for j := range shipments[i].PPMShipment.WeightTickets {
				// variable for convience still modifies original shipments object
				weightTicket := &shipments[i].PPMShipment.WeightTickets[j]

				weightTicket.EmptyDocument.UserUploads = weightticket.FilterDeletedValued(weightTicket.EmptyDocument.UserUploads)

				loadErr := appCtx.DB().Load(weightTicket, "FullDocument.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				weightTicket.FullDocument.UserUploads = weightticket.FilterDeletedValued(weightTicket.FullDocument.UserUploads)

				loadErr = appCtx.DB().Load(weightTicket, "ProofOfTrailerOwnershipDocument.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightticket.FilterDeletedValued(weightTicket.ProofOfTrailerOwnershipDocument.UserUploads)
			}
		}
	}

	return shipments, nil
}

func FindShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eagerAssociations ...string) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	findShipmentQuery := appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope())

	if len(eagerAssociations) > 0 {
		findShipmentQuery.Eager(eagerAssociations...)
	}

	err := findShipmentQuery.Find(&shipment, shipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}
