package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

// DistanceZipSITDestLookup does the lookup of distance for SIT at destination
type DistanceZipSITDestLookup struct {
	DestinationAddress      models.Address
	FinalDestinationAddress models.Address
}

func (r DistanceZipSITDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner

	// If the zip3s are identical, we do a zip3 distance calc (which uses RM).
	// If they are different, we do a zip5 distance calc (which uses DTOD).

	destZip := r.DestinationAddress.PostalCode
	if len(destZip) < 5 {
		return "", fmt.Errorf("invalid destination postal code of %s", destZip)
	}
	destZip3 := destZip[:3]

	finalDestZip := r.FinalDestinationAddress.PostalCode
	if len(finalDestZip) < 5 {
		return "", fmt.Errorf("invalid SIT final destination postal code of %s", destZip)
	}
	finalDestZip3 := finalDestZip[:3]

	var distanceMiles int
	var distanceErr error
	if destZip3 == finalDestZip3 {
		distanceMiles, distanceErr = planner.Zip5TransitDistance(destZip, finalDestZip)
	} else {
		distanceMiles, distanceErr = planner.Zip3TransitDistance(destZip, finalDestZip)
	}
	if distanceErr != nil {
		return "", distanceErr
	}

	return strconv.Itoa(distanceMiles), nil
}
