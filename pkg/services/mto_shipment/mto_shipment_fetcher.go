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
	moveQuery := appCtx.DB().Q()

	if appCtx.Session().IsMilApp() {
		serviceMemberID := appCtx.Session().ServiceMemberID
		moveQuery = moveQuery.
			LeftJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", serviceMemberID)
	}

	err := moveQuery.Find(&move, moveID)

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
			"PickupAddress.Country",
			"SecondaryPickupAddress.Country",
			"TertiaryPickupAddress.Country",
			"DestinationAddress.Country",
			"SecondaryDeliveryAddress.Country",
			"TertiaryDeliveryAddress.Country",
			"MTOServiceItems.Dimensions",
			"BoatShipment",
			"MobileHome",
			"PPMShipment.W2Address",
			"PPMShipment.WeightTickets",
			"PPMShipment.MovingExpenses",
			"PPMShipment.ProgearWeightTickets",
			"PPMShipment.PickupAddress",
			"PPMShipment.SecondaryPickupAddress",
			"PPMShipment.TertiaryPickupAddress",
			"PPMShipment.DestinationAddress",
			"PPMShipment.SecondaryDestinationAddress",
			"PPMShipment.TertiaryDestinationAddress",
			"DeliveryAddressUpdate",
			"DeliveryAddressUpdate.OriginalAddress",
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

		if shipments[i].DeliveryAddressUpdate != nil {
			// Cannot EagerPreload the address update `NewAddress` due to POP bug
			// See: https://transcom.github.io/mymove-docs/docs/backend/setup/using-eagerpreload-in-pop#eager-vs-eagerpreload-inconsistency
			loadErr := appCtx.DB().Load(shipments[i].DeliveryAddressUpdate, "NewAddress", "SitOriginalAddress")
			if loadErr != nil {
				return nil, apperror.NewQueryError("DeliveryAddressUpdate", loadErr, "")
			}
		}

		var agents []models.MTOAgent
		err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", shipments[i].ID).All(&agents)
		if err != nil {
			return nil, err
		}
		shipments[i].MTOAgents = agents
	}

	return shipments, nil
}

func FindShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eagerAssociations ...string) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	findShipmentQuery := appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope())
	if len(eagerAssociations) > 0 {
		findShipmentQuery.Eager(eagerAssociations...)
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		findShipmentQuery.
			LeftJoin("moves", "moves.id = mto_shipments.move_id").
			LeftJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
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

// This allows us to gather all possible parent and child shipments in the diverted shipment chain
func (f mtoShipmentFetcher) GetDiversionChain(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*[]models.MTOShipment, error) {
	var allShipmentsInChain []models.MTOShipment
	var initialShipment models.MTOShipment

	// Grab the initial shipment the reweight was requested for
	err := appCtx.DB().Find(&initialShipment, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(shipmentID, "while looking for shipment")
		}
		return nil, apperror.NewQueryError("MTOShipment", err, "")
	}

	allShipmentsInChain = append(allShipmentsInChain, initialShipment)

	// Loop over the "parent" shipments by DivertedFromShipmentID until no more IDs are found (No more parent shipments are found)
	currentShipmentID := initialShipment.DivertedFromShipmentID
	for currentShipmentID != nil {
		var parentShipment models.MTOShipment
		err := appCtx.DB().Find(&parentShipment, *currentShipmentID)
		if err != nil {
			if err == sql.ErrNoRows {
				// If sql ErrNoRows pops up we know there are no more parent shipments
				break
			}
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
		allShipmentsInChain = append(allShipmentsInChain, parentShipment)
		currentShipmentID = parentShipment.DivertedFromShipmentID
	}

	// Loop over the "child" shipments by parent ID to child DivertedFromShipmentID until no more child shipments are found
	// The loop will break when no more child shipments can be found
	currentShipmentID = &initialShipment.ID
	for {
		var childShipment models.MTOShipment
		err := appCtx.DB().Where("diverted_from_shipment_id = ?", *currentShipmentID).First(&childShipment)
		if err != nil {
			// If sql ErrNoRows pops up we know there are no more child shipments
			if err == sql.ErrNoRows {
				break
			}
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}

		if childShipment.DivertedFromShipmentID == nil || *childShipment.DivertedFromShipmentID != *currentShipmentID {
			// No more child shipments
			break
		}

		allShipmentsInChain = append(allShipmentsInChain, childShipment)
		currentShipmentID = &childShipment.ID
	}

	return &allShipmentsInChain, nil
}
