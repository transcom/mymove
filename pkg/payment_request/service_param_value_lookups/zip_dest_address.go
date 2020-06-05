package serviceparamvaluelookups

import (
	"fmt"

	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ZipDestAddress does lookup on actual weight billed
type ZipDestAddress struct {
}

func (r ZipDestAddress) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService", "MTOShipment", "MTOShipment.DestinationAddress").Find(&mtoServiceItem, mtoServiceItemID)
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

	// Make sure there's a destination address since those are nullable
	destAddressID := mtoServiceItem.MTOShipment.DestinationAddressID
	if destAddressID == nil || *destAddressID == uuid.Nil {
		//check for string of all zeros
		return "", services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
	}

	value := fmt.Sprintf("%+v", mtoServiceItem.MTOShipment.DestinationAddress.PostalCode)
	return value, nil
}
