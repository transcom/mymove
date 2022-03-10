package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
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

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := keyData.mtoShipmentID
	if mtoShipmentID == nil {
		return "", apperror.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	var mtoShipment models.MTOShipment
	err := db.Find(&mtoShipment, keyData.mtoShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(mtoServiceItemID, "looking for MTOShipmentID")
		default:
			return "", apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	if mtoShipment.Distance != nil {
		return strconv.Itoa(mtoShipment.Distance.Int()), nil
	}

	// Now calculate the distance between zip5s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip5TransitDistance(appCtx, pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	if len(pickupZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("Shipment must have valid pickup zipcode. Received: %s", pickupZip), nil, fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip))
	}
	if len(destinationZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("Shipment must have valid destination zipcode. Received: %s", destinationZip), nil, fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip))
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
