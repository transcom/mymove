package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// RequestedPickupDateLookup does lookup on requested pickup date
type RequestedPickupDateLookup struct {
}

func (r RequestedPickupDateLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
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

	// Make sure there's a requested pickup date since that's nullable
	requestedPickupDate := mtoServiceItem.MTOShipment.RequestedPickupDate
	if requestedPickupDate == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find a requested pickup date for MTOShipmentID [%s]", mtoShipmentID)
	}

	return requestedPickupDate.Format("2006-01-02"), nil
}
