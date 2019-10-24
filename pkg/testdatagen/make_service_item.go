package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceItem creates a single office user and associated TransportOffice
func MakeServiceItem(db *pop.Connection, assertions Assertions) models.ServiceItem {
	mto := assertions.ServiceItem.MoveTaskOrder

	if assertions.ServiceItem.MoveTaskOrderID == uuid.Nil {
		mto = MakeMoveTaskOrder(db, assertions)
	}

	serviceItem := models.ServiceItem{
		MoveTaskOrderID: mto.ID,
	}

	mergeModels(&serviceItem, assertions.ServiceItem)

	mustCreate(db, &serviceItem)

	return serviceItem
}

// MakeDefaultServiceItem makes an ServiceItem with default values
func MakeDefaultServiceItem(db *pop.Connection) models.ServiceItem {
	return MakeServiceItem(db, Assertions{})
}
