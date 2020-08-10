package serviceparamvaluelookups

import (
	"database/sql"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DistanceZip5Lookup contains zip5 lookup
type DistanceZip5Lookup struct {
}

func (r DistanceZip5Lookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db
	planner := keyData.planner

	// Get the MTOServiceItem and associated MTOShipment and addresses
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.
		Eager("MTOShipment", "MTOShipment.PickupAddress", "MTOShipment.DestinationAddress").
		Find(&mtoServiceItem, mtoServiceItemID)
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

	mtoShipment := mtoServiceItem.MTOShipment
	if mtoShipment.Distance != nil {
		return strconv.Itoa(mtoShipment.Distance.Int()), nil
	}

	// Make sure there's a pickup and destination address since those are nullable
	pickupAddressID := mtoServiceItem.MTOShipment.PickupAddressID
	if pickupAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
	}
	destinationAddressID := mtoServiceItem.MTOShipment.DestinationAddressID
	if destinationAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
	}

	// Now calculate the distance between zip5s
	pickupZip := mtoServiceItem.MTOShipment.PickupAddress.PostalCode
	destinationZip := mtoServiceItem.MTOShipment.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip5TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	if distanceMiles < 50 {
		miles := unit.Miles(distanceMiles)
		mtoShipment.Distance = &miles
		db.Save(&mtoShipment)
	}

	return strconv.Itoa(distanceMiles), nil
}
