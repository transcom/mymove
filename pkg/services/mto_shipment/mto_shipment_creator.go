package mtoshipment

import (
	"fmt"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/fetch"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type createMTOShipmentQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
	Transaction(fn func(tx *pop.Connection) error) error
}

type mtoShipmentCreator struct {
	db      *pop.Connection
	builder createMTOShipmentQueryBuilder
	services.Fetcher
	createNewBuilder      func(db *pop.Connection) createMTOShipmentQueryBuilder
	mtoServiceItemCreator services.MTOServiceItemCreator
}

// NewMTOShipmentCreator creates a new struct with the service dependencies
func NewMTOShipmentCreator(db *pop.Connection, builder createMTOShipmentQueryBuilder, fetcher services.Fetcher) services.MTOShipmentCreator {
	createNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
		return query.NewQueryBuilder(db)
	}

	return &mtoShipmentCreator{
		db,
		builder,
		fetch.NewFetcher(builder),
		createNewBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(builder),
	}
}

// CreateMTOShipment updates the mto shipment
func (f mtoShipmentCreator) CreateMTOShipment(shipment *models.MTOShipment, serviceItems models.MTOServiceItems) (*models.MTOShipment, error) {
	var verrs *validate.Errors
	var err error

	err = checkShipmentIDFields(shipment, serviceItems)
	if err != nil {
		return nil, err
	}

	var move models.Move
	moveID := shipment.MoveTaskOrderID

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveID),
	}

	// check if Move exists
	err = f.builder.FetchOne(&move, queryFilters)
	if err != nil {
		return nil, services.NewNotFoundError(moveID, "for move")
	}

	for _, existingShipment := range move.MTOShipments {
		hasNTSShipment := shipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom &&
			(existingShipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom && existingShipment.Status == models.MTOShipmentStatusSubmitted)

		hasNTSRShipment := shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && (existingShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && existingShipment.Status == models.MTOShipmentStatusSubmitted)

		if hasNTSShipment || hasNTSRShipment {
			return nil, services.NewInvalidInputError(existingShipment.ID, nil, nil, "Cannot create another NTS Shipment")
		}
	}

	if serviceItems != nil {
		serviceItemsList := make(models.MTOServiceItems, 0, len(serviceItems))

		for _, serviceItem := range serviceItems {
			// find the re service code id
			var reService models.ReService
			reServiceCode := serviceItem.ReService.Code
			queryFilters = []services.QueryFilter{
				query.NewQueryFilter("code", "=", reServiceCode),
			}
			err = f.builder.FetchOne(&reService, queryFilters)
			if err != nil {
				return nil, services.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service item with code: %s", reServiceCode))
			}
			// set re service for service item
			serviceItem.ReServiceID = reService.ID
			serviceItem.Status = models.MTOServiceItemStatusSubmitted

			if serviceItem.ReService.Code == models.ReServiceCodeDOSHUT || serviceItem.ReService.Code == models.ReServiceCodeDDSHUT {
				if shipment.PrimeEstimatedWeight == nil {
					return nil, services.NewConflictError(
						serviceItem.ReService.ID,
						"for creating a service item. MTOShipment associated with this service item must have a valid PrimeEstimatedWeight.",
					)
				}
			}

			serviceItemsList = append(serviceItemsList, serviceItem)
		}
		shipment.MTOServiceItems = serviceItemsList
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		// create new builder to use tx
		txBuilder := f.createNewBuilder(tx)

		// create pickup and destination addresses
		if shipment.PickupAddress != nil {
			verrs, err = txBuilder.CreateOne(shipment.PickupAddress)
			if verrs != nil || err != nil {
				return fmt.Errorf("failed to create pickup address %#v %e", verrs, err)
			}
			shipment.PickupAddressID = &shipment.PickupAddress.ID
		} else if shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom {
			return services.NewInvalidInputError(uuid.Nil, nil, nil, "PickupAddress is required to create an HHG or NTS type MTO shipment")
		}

		if shipment.DestinationAddress != nil {
			verrs, err = txBuilder.CreateOne(shipment.DestinationAddress)
			if verrs != nil || err != nil {
				return fmt.Errorf("failed to create destination address %#v %e", verrs, err)
			}
			shipment.DestinationAddressID = &shipment.DestinationAddress.ID
		}

		// check that required items to create shipment are present
		if shipment.RequestedPickupDate.IsZero() && shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom {
			return services.NewInvalidInputError(uuid.Nil, nil, nil, "RequestedPickupDate is required to create an HHG or NTS type MTO shipment")
		}

		//assign status to shipment draft by default
		if shipment.Status != models.MTOShipmentStatusSubmitted {
			shipment.Status = models.MTOShipmentStatusDraft
		}

		// create a shipment
		verrs, err = txBuilder.CreateOne(shipment)

		if verrs != nil || err != nil {
			return fmt.Errorf("failed to create shipment %s %e", verrs.Error(), err)
		}

		// create MTOAgents List
		if shipment.MTOAgents != nil {
			agentsList := make(models.MTOAgents, 0, len(shipment.MTOAgents))

			for _, agent := range shipment.MTOAgents {
				agent.MTOShipmentID = shipment.ID
				// #nosec G601 TODO needs review
				verrs, err = txBuilder.CreateOne(&agent)
				if verrs != nil && verrs.HasAny() {
					return verrs
				}
				if err != nil {
					return err
				}
				agentsList = append(agentsList, agent)
			}
			shipment.MTOAgents = agentsList
		}

		// create MTOServiceItems List
		if shipment.MTOServiceItems != nil {
			serviceItemsList := make(models.MTOServiceItems, 0, len(shipment.MTOServiceItems))

			for _, serviceItem := range shipment.MTOServiceItems {
				serviceItem.MTOShipmentID = &shipment.ID
				serviceItem.MoveTaskOrderID = shipment.MoveTaskOrderID
				// #nosec G601 TODO needs review
				verrs, err = txBuilder.CreateOne(&serviceItem)
				if verrs != nil && verrs.HasAny() {
					return verrs
				}
				if err != nil {
					return err
				}
				serviceItemsList = append(serviceItemsList, serviceItem)
			}
			shipment.MTOServiceItems = serviceItemsList
		}

		return nil
	})

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "Unable to create shipment")
	} else if err != nil {
		return nil, services.NewQueryError("unknown", err, "")
	}

	return shipment, transactionError
}

// checkShipmentIDFields checks that the client hasn't attempted to set ID fields that will be generated/auto-set
func checkShipmentIDFields(shipment *models.MTOShipment, serviceItems models.MTOServiceItems) error {
	verrs := validate.NewErrors()

	if shipment.MTOAgents != nil && len(shipment.MTOAgents) > 0 {
		for _, agent := range shipment.MTOAgents {
			if agent.ID != uuid.Nil {
				verrs.Add("agents:id", "cannot be set for new agents")
			}
			if agent.MTOShipmentID != uuid.Nil {
				verrs.Add("agents:mtoShipmentID", "cannot be set for agents created with a shipment")
			}
		}
	}

	if serviceItems != nil && len(serviceItems) > 0 {
		for _, item := range serviceItems {
			if item.ID != uuid.Nil {
				verrs.Add("mtoServiceItems:id", "cannot be set for new service items")
			}
			if item.MTOShipmentID != nil && *item.MTOShipmentID != uuid.Nil {
				verrs.Add("mtoServiceItems:mtoShipmentID", "cannot be set for service items created with a shipment")
			}
		}
	}

	addressMsg := "cannot be set for new addresses"
	if shipment.PickupAddress != nil && shipment.PickupAddress.ID != uuid.Nil {
		verrs.Add("pickupAddress:id", addressMsg)
	}
	if shipment.DestinationAddress != nil && shipment.DestinationAddress.ID != uuid.Nil {
		verrs.Add("destinationAddress:id", addressMsg)
	}
	if shipment.SecondaryPickupAddress != nil && shipment.SecondaryPickupAddress.ID != uuid.Nil {
		verrs.Add("secondaryPickupAddress:id", addressMsg)
	}
	if shipment.SecondaryDeliveryAddress != nil && shipment.SecondaryDeliveryAddress.ID != uuid.Nil {
		verrs.Add("secondaryDestinationAddress:id", addressMsg)
	}

	if verrs.HasAny() {
		return services.NewInvalidInputError(uuid.Nil, nil, verrs, "Fields that cannot be set found while creating new shipment.")
	}

	return nil
}
