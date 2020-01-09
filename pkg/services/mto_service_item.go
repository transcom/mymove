package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// MTOServiceItemCreator is the exported interface for creating a mto service item
//go:generate mockery -name MTOServiceItemCreator
type MTOServiceItemCreator interface {
	CreateMTOServiceItem(serviceItem *models.MTOServiceItem) (*models.MTOServiceItem, *validate.Errors, error)
}
