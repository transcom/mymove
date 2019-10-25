package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// ServiceItemListFetcher is the exported interface for fetching multiple transportation offices
//go:generate mockery -name ServiceItemListFetcher
type ServiceItemListFetcher interface {
	FetchServiceItemList(moveTaskOrderID uuid.UUID) (models.ServiceItems, error)
}

// ServiceItemCreator is the exported interface for fetching multiple transportation offices
//go:generate mockery -name ServiceItemCreator
type ServiceItemCreator interface {
	CreateServiceItem(serviceItem *models.ServiceItem, moveTaskOrderIDFilter []QueryFilter) (*models.ServiceItem, *validate.Errors, error)
}
