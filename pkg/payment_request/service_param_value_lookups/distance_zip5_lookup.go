package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DistanceZip5Lookup contains zip5 lookup
type DistanceZip5Lookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZip5Lookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := appCtx.DB()

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

	// Now calculate the distance between zip5s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip5TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	if len(pickupZip) < 5 {
		return "", services.NewInvalidInputError(*mtoServiceItem.MTOShipmentID, fmt.Errorf("Shipment must have valid pickup zipcode. Received: %s", pickupZip), nil, fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip))
	}
	if len(destinationZip) < 5 {
		return "", services.NewInvalidInputError(*mtoServiceItem.MTOShipmentID, fmt.Errorf("Shipment must have valid destination zipcode. Received: %s", destinationZip), nil, fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip))
	}

	pickupZip3 := pickupZip[:3]
	destinationZip3 := destinationZip[:3]

	if pickupZip3 == destinationZip3 {
		miles := unit.Miles(distanceMiles)
		mtoShipment.Distance = &miles
		err := db.Save(&mtoShipment)
		if err != nil {
			return "", err
		}
	}

	return strconv.Itoa(distanceMiles), nil
}
