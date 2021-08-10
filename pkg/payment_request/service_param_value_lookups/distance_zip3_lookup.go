package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DistanceZip3Lookup contains zip3 lookup
type DistanceZip3Lookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZip3Lookup) lookup(appCfg appconfig.AppConfig, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := appCfg.DB()

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

	// Now calculate the distance between zip3s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip3TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	errorMsgForPickupZip := fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip)
	errorMsgForDestinationZip := fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip)
	if len(pickupZip) < 5 {
		return "", services.NewInvalidInputError(*mtoServiceItem.MTOShipmentID, fmt.Errorf(errorMsgForPickupZip), nil, errorMsgForPickupZip)
	}
	if len(destinationZip) < 5 {
		return "", services.NewInvalidInputError(*mtoServiceItem.MTOShipmentID, fmt.Errorf(errorMsgForDestinationZip), nil, errorMsgForDestinationZip)
	}

	pickupZip3 := pickupZip[:3]
	destinationZip3 := destinationZip[:3]
	if pickupZip3 != destinationZip3 {
		miles := unit.Miles(distanceMiles)
		mtoShipment.Distance = &miles
		err = db.Save(&mtoShipment)
		if err != nil {
			return "", err
		}
	}

	return strconv.Itoa(distanceMiles), nil
}
