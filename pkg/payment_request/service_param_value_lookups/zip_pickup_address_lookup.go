package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ZipPickupAddressLookup does lookup on the postal code for the pickup address
type ZipPickupAddressLookup struct {
}

func (r ZipPickupAddressLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("MTOShipment", "MTOShipment.PickupAddress").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipmentID
	if mtoShipmentID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there's a pickup address zip code since pickupAddress is nullable
	pickupAddress := mtoServiceItem.MTOShipment.PickupAddress

	if pickupAddress == nil {
		return "", fmt.Errorf("could not find pickup address for MTOShipment [%s]", mtoShipmentID)
	}

	zipPickupAddress := pickupAddress.PostalCode

	value := fmt.Sprintf("%s", zipPickupAddress)
	return value, nil
}
