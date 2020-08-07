package serviceparamvaluelookups

import (
	"database/sql"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// DistanceZip3Lookup contains zip3 lookup
type DistanceZip3Lookup struct {
	MTOShipment models.MTOShipment
}

func (r DistanceZip3Lookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db
	planner := keyData.planner

	// Make sure there's a pickup and destination address since those are nullable
	pickupAddressID := r.MTOShipment.PickupAddressID
	if pickupAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
	}
	destinationAddressID := r.MTOShipment.DestinationAddressID
	if destinationAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
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

	var destinationAddress models.Address
	err = db.Find(&destinationAddress, r.MTOShipment.DestinationAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*r.MTOShipment.DestinationAddressID, "looking for DestinationAddressID")
		default:
			return "", err
		}
	}

	// Now calculate the distance between zip3s
	pickupZip := pickupAddress.PostalCode
	destinationZip := destinationAddress.PostalCode
	distanceMiles, err := planner.Zip3TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(distanceMiles), nil
}
