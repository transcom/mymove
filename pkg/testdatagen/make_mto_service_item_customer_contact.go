package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItemCustomerContact creates a single customer contact and associated set relationships
func MakeMTOServiceItemCustomerContact(db *pop.Connection, assertions Assertions) (models.MTOServiceItemCustomerContact, error) {
	MTOServiceItem := assertions.MTOServiceItem
	if isZeroUUID(MTOServiceItem.ID) {
		var err error
		MTOServiceItem, err = MakeMTOServiceItem(db, assertions)
		if err != nil {
			return models.MTOServiceItemCustomerContact{}, err
		}
	}

	MTOServiceItemCustomerContact := models.MTOServiceItemCustomerContact{
		Type:                       models.CustomerContactTypeFirst,
		DateOfContact:              time.Now(),
		TimeMilitary:               "0400Z",
		FirstAvailableDeliveryDate: time.Now(),
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItemCustomerContact, assertions.MTOServiceItemCustomerContact)

	mustCreate(db, &MTOServiceItemCustomerContact, assertions.Stub)

	return MTOServiceItemCustomerContact, nil
}
