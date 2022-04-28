package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItemCustomerContact creates a single customer contact and associated set relationships
func MakeMTOServiceItemCustomerContact(db *pop.Connection, assertions Assertions) models.MTOServiceItemCustomerContact {
	MTOServiceItem := assertions.MTOServiceItem
	if isZeroUUID(MTOServiceItem.ID) {
		MTOServiceItem = MakeMTOServiceItem(db, assertions)
	}

	MTOServiceItemCustomerContact := models.MTOServiceItemCustomerContact{
		MTOServiceItemID:           MTOServiceItem.ID,
		MTOServiceItem:             MTOServiceItem,
		Type:                       models.CustomerContactTypeFirst,
		TimeMilitary:               "0400Z",
		FirstAvailableDeliveryDate: time.Now(),
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItemCustomerContact, assertions.MTOServiceItemCustomerContact)

	mustCreate(db, &MTOServiceItemCustomerContact, assertions.Stub)

	return MTOServiceItemCustomerContact
}
