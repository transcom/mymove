package mtoshipment

import (
	"fmt"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/fetch"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

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

	var moveTaskOrder models.MoveTaskOrder
	moveTaskOrderID := shipment.MoveTaskOrderID

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveTaskOrderID),
	}

	// check if MTO exists
	err = f.builder.FetchOne(&moveTaskOrder, queryFilters)
	if err != nil {
		return nil, services.NewNotFoundError(moveTaskOrderID, "for moveTaskOrder")
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

			println("🧸")
			fmt.Printf("service item %v", serviceItem)
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
		} else {
			// Swagger should pick this up before it ever gets here
			return services.NewInvalidInputError(uuid.Nil, nil, nil, "PickupAddress is required to create MTO shipment")
		}

		if shipment.DestinationAddress != nil {
			verrs, err = txBuilder.CreateOne(shipment.DestinationAddress)
			if verrs != nil || err != nil {
				return fmt.Errorf("failed to create destination address %#v %e", verrs, err)
			}
			shipment.DestinationAddressID = &shipment.DestinationAddress.ID
		} else {
			// Swagger should pick this up before it ever gets here
			return services.NewInvalidInputError(uuid.Nil, nil, nil, "DestinationAddress is required to create MTO shipment")
		}

		// check that required items to create shipment are present
		if shipment.RequestedPickupDate == nil {
			return services.NewInvalidInputError(uuid.Nil, nil, nil, "RequestedPickupDate is required to create MTO shipment")
		}

		//assign status to shipment submitted
		shipment.Status = models.MTOShipmentStatusSubmitted

		// create a shipment
		verrs, err = txBuilder.CreateOne(shipment)

		if verrs != nil || err != nil {
			return fmt.Errorf("failed to create shipment %s %e", verrs.Error(), err)
		}

		// create MTOAgents List
		if shipment.MTOAgents != nil {
			for _, agent := range shipment.MTOAgents {
				agent.MTOShipmentID = shipment.ID
				verrs, err = txBuilder.CreateOne(&agent)
				if err != nil {
					return err
				}
			}
		}

		// create MTOServiceItems List
		if shipment.MTOServiceItems != nil {
			serviceItemsList := make(models.MTOServiceItems, 0, len(serviceItems))

			for _, serviceItem := range shipment.MTOServiceItems {
				serviceItem.MTOShipmentID = &shipment.ID
				serviceItem.MoveTaskOrderID = shipment.MoveTaskOrderID
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
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "")
	} else if err != nil {
		return nil, services.NewQueryError("unknown", err, "")
	}

	return shipment, transactionError
}
