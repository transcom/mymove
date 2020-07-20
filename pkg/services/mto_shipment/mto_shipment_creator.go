package mtoshipment

import (
	"fmt"

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
	createNewBuilder func(db *pop.Connection) createMTOShipmentQueryBuilder
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

		return nil
	})

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "")
	} else if err != nil {
		return nil, services.NewQueryError("unknown", err, "")
	}

	return shipment, transactionError
}
