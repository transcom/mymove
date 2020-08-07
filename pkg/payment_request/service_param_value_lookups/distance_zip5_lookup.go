package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

// DistanceZip5Lookup contains zip5 lookup
type DistanceZip5Lookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZip5Lookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner

	// Now calculate the distance between zip5s
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	distanceMiles, err := planner.Zip5TransitDistance(pickupZip, destinationZip)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(distanceMiles), nil
}
