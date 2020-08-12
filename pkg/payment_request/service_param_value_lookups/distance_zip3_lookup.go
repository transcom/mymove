package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

// DistanceZip3Lookup contains zip3 lookup
type DistanceZip3Lookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZip3Lookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner

	// Now calculate the distance between zip3s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip3TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(distanceMiles), nil
}
