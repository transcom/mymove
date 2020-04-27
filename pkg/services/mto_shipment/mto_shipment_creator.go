package mtoshipment

import (
	"fmt"
	"github.com/transcom/mymove/pkg/services/fetch"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

//fetching an MTO
//creating a mtoShipment
//creating accessorials (find service written for this)

type createMTOShipmentQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
	Transaction(fn func(tx *pop.Connection) error) error
}

type mtoShipmentCreator struct {
	db *pop.Connection
	builder createMTOShipmentQueryBuilder
	services.Fetcher
	createNewBuilder func(db *pop.Connection) createMTOShipmentQueryBuilder
}

// NewMTOShipmentCreator creates a new struct with the service dependencies
func NewMTOShipmentCreator(db *pop.Connection, builder createMTOShipmentQueryBuilder, fetcher services.Fetcher) services.MTOShipmentCreator {
	createNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
		return query.NewQueryBuilder(db)
	}
	return &mtoShipmentCreator{db, builder, fetch.NewFetcher(builder), createNewBuilder }
}

// CreateMTOShipment updates the mto shipment
func (f mtoShipmentCreator) CreateMTOShipment(shipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
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
		return nil, services.NewNotFoundError(moveTaskOrderID, fmt.Sprintf("MoveTaskOrderID: %s", err))
	}

	err = f.db.Transaction(func(tx *pop.Connection) error {
		txBuilder := f.createNewBuilder(tx)
		verrs, err = txBuilder.CreateOne(shipment)

		if verrs != nil || err != nil {
			return fmt.Errorf("%#v %e", verrs, err)
		}

		// temp optimistic locking solution til query builder is re-tooled to handle nested updates
		encodedUpdatedAt := etag.GenerateEtag(shipment.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		if shipment.DestinationAddress != nil || shipment.PickupAddress != nil {
			addressBaseQuery := `UPDATE addresses
				SET
			`

			if shipment.DestinationAddress != nil {
				destinationAddressQuery := generateAddressQuery()
				params := generateAddressParams(shipment.DestinationAddress)
				err = f.db.RawQuery(addressBaseQuery+destinationAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}

			if shipment.PickupAddress != nil {
				pickupAddressQuery := generateAddressQuery()
				params := generateAddressParams(shipment.PickupAddress)
				err = f.db.RawQuery(addressBaseQuery+pickupAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}

			if shipment.MTOAgents != nil {
				agentQuery := `UPDATE mto_agents
					SET
				`
				for _, agent := range shipment.MTOAgents {
					updateAgentQuery := generateAgentQuery()
					params := generateMTOAgentsParams(agent)
					err = f.db.RawQuery(agentQuery+updateAgentQuery, params...).Exec()
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	var newShipment models.MTOShipment
	err = f.FetchRecord(&newShipment, queryFilters)

	if err != nil {
		return &models.MTOShipment{}, err
	}

	return &newShipment, nil
}

