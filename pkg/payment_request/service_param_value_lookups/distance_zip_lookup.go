package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DistanceZipLookup contains zip lookup
type DistanceZipLookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZipLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := keyData.db

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

	// Now calculate the distance between zips
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode

	if len(strings.TrimSpace(pickupZip)) != 5 && len(strings.TrimSpace(destinationZip)) != 5 {
		return "", services.NewBadDataError(
			fmt.Sprintf("Both ZIPs are not of length 5 pickupZIP %s and destinationZIP %s", pickupZip, destinationZip))
	}
	distanceMiles, err := distanceZip(planner, pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	/*
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
	*/

	miles := unit.Miles(distanceMiles)
	mtoShipment.Distance = &miles
	err = db.Save(&mtoShipment)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(distanceMiles), nil
}
