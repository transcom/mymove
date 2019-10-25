package serviceitem

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type createServiceItemQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
}

type serviceItemCreator struct {
	builder createServiceItemQueryBuilder
}

func (o *serviceItemCreator) CreateServiceItem(serviceItem *models.ServiceItem, moveTaskOrderIDFilter []services.QueryFilter) (*models.ServiceItem, *validate.Errors, error) {
	// Use FetchOne to see if we have a move task order that matches the provided id
	var moveTaskOrder models.MoveTaskOrder
	err := o.builder.FetchOne(&moveTaskOrder, moveTaskOrderIDFilter)

	if err != nil {
		return nil, nil, err
	}

	verrs, err := o.builder.CreateOne(serviceItem)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return serviceItem, nil, nil
}

func NewServiceItemCreator(builder createServiceItemQueryBuilder) services.ServiceItemCreator {
	return &serviceItemCreator{builder}
}
