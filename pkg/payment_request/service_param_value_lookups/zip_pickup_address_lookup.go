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
	MTOShipment models.MTOShipment
}

func (r ZipPickupAddressLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Make sure there's a pickup and destination address since those are nullable
	pickupAddressID := r.MTOShipment.PickupAddressID
	if pickupAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
	}
	var pickupAddress models.Address
	err := db.Find(&pickupAddress, r.MTOShipment.PickupAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*r.MTOShipment.PickupAddressID, "looking for PickupAddressID")
		default:
			return "", err
		}
	}

	value := fmt.Sprintf("%s", pickupAddress.PostalCode)
	return value, nil
}
