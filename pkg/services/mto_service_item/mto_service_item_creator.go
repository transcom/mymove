package mtoserviceitem

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type createMTOServiceItemQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
}

type mtoServiceItemCreator struct {
	builder createMTOServiceItemQueryBuilder
}

// CreateMTOServiceItem creates an MTO Service Item
func (o *mtoServiceItemCreator) CreateMTOServiceItem(serviceItem *models.MTOServiceItem) (*models.MTOServiceItem, *validate.Errors, error) {
	var moveTaskOrder models.MoveTaskOrder
	moveTaskOrderID := serviceItem.MoveTaskOrderID
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveTaskOrderID),
	}
	// check if MTO exists
	err := o.builder.FetchOne(&moveTaskOrder, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(moveTaskOrderID, fmt.Sprintf("MoveTaskOrderID: %s", err))
	}

	// check if shipment exists
	var mtoShipment models.MTOShipment
	var mtoShipmentID uuid.UUID
	if serviceItem.MTOShipmentID != nil {
		mtoShipmentID = *serviceItem.MTOShipmentID
	}
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoShipmentID),
	}
	err = o.builder.FetchOne(&mtoShipment, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(mtoShipmentID, fmt.Sprintf("MTOShipmentID: %s", err))
	}

	// find the re service code id
	var reService models.ReService
	reServiceCode := serviceItem.ReService.Code
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}
	err = o.builder.FetchOne(&reService, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(uuid.Nil, fmt.Sprintf("failed to find service item code: %s; %s", reServiceCode, err))
	}

	// set re service for service item
	serviceItem.ReServiceID = reService.ID

	if serviceItem.ReService.Code == models.ReServiceCodeDOSHUT || serviceItem.ReService.Code == models.ReServiceCodeDDSHUT {
		if mtoShipment.PrimeEstimatedWeight == nil {
			return nil, nil, services.NewInvalidInputError(mtoShipmentID, nil, nil,
				fmt.Sprintf("MTOShipment with id: %s is missing the estimated weight required for this service item", mtoShipmentID))
		}
	}

	verrs, err := o.builder.CreateOne(serviceItem)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return serviceItem, nil, nil
}

// NewMTOServiceItemCreator returns a new MTO service item creator
func NewMTOServiceItemCreator(builder createMTOServiceItemQueryBuilder) services.MTOServiceItemCreator {
	return &mtoServiceItemCreator{builder}
}
