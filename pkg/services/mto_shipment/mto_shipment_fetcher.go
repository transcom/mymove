package mtoshipment

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
			"PPMShipment.W2Address",
			"PPMShipment.WeightTickets",
			"PPMShipment.MovingExpenses",
			"PPMShipment.ProgearWeightTickets",
			"Reweigh",
			"SITDurationUpdates",
			"StorageFacility.Address",
		). // Right now no use case for showing deleted shipments.
		Where("move_id = ?", moveID).
		Order("uses_external_vendor asc").
		Order("created_at asc").
		All(&shipments)

	if err != nil {
		return nil, err
	}

	// Need to iterate through shipments to fetch additional PPM weight ticket info
	// EagerPreload causes duplicate records because there are multiple relationships to the same table
	for i := range shipments {
		if shipments[i].ShipmentType == models.MTOShipmentTypePPM {
			shipments[i].PPMShipment.WeightTickets = shipments[i].PPMShipment.WeightTickets.FilterDeleted()
			for j := range shipments[i].PPMShipment.WeightTickets {
				// variable for convenience still modifies original shipments object
				weightTicket := &shipments[i].PPMShipment.WeightTickets[j]

				loadErr := appCtx.DB().Load(weightTicket, "EmptyDocument.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				weightTicket.EmptyDocument.UserUploads = weightTicket.EmptyDocument.UserUploads.FilterDeleted()

				loadErr = appCtx.DB().Load(weightTicket, "FullDocument.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				weightTicket.FullDocument.UserUploads = weightTicket.FullDocument.UserUploads.FilterDeleted()

				loadErr = appCtx.DB().Load(weightTicket, "ProofOfTrailerOwnershipDocument.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightTicket.ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()
			}

			shipments[i].PPMShipment.MovingExpenses = shipments[i].PPMShipment.MovingExpenses.FilterDeleted()
			for j := range shipments[i].PPMShipment.MovingExpenses {
				movingExpense := &shipments[i].PPMShipment.MovingExpenses[j]

				loadErr := appCtx.DB().Load(movingExpense, "Document.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				movingExpense.Document.UserUploads = movingExpense.Document.UserUploads.FilterDeleted()
			}

			shipments[i].PPMShipment.ProgearWeightTickets = shipments[i].PPMShipment.ProgearWeightTickets.FilterDeleted()
			for j := range shipments[i].PPMShipment.ProgearWeightTickets {
				progearWeightTicket := &shipments[i].PPMShipment.ProgearWeightTickets[j]

				loadErr := appCtx.DB().Load(progearWeightTicket, "Document.UserUploads.Upload")
				if loadErr != nil {
					return nil, loadErr
				}
				progearWeightTicket.Document.UserUploads = progearWeightTicket.Document.UserUploads.FilterDeleted()
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

func (f mtoShipmentFetcher) GetShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eagerAssociations ...string) (*models.MTOShipment, error) {
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
