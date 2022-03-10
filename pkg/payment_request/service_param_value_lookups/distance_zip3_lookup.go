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

// DistanceZip3Lookup contains zip3 lookup
type DistanceZip3Lookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZip3Lookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := appCtx.DB()

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
			return "", apperror.NewNotFoundError(*mtoShipmentID, "looking for MTOShipmentID")
		default:
			return "", apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	if mtoShipment.Distance != nil {
		return strconv.Itoa(mtoShipment.Distance.Int()), nil
	}

	// Now calculate the distance between zip3s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip3TransitDistance(appCtx, pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	errorMsgForPickupZip := fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip)
	errorMsgForDestinationZip := fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip)
	if len(pickupZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf(errorMsgForPickupZip), nil, errorMsgForPickupZip)
	}
	if len(destinationZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf(errorMsgForDestinationZip), nil, errorMsgForDestinationZip)
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
