package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

// DistanceZipSITOriginLookup does the lookup of distance for SIT at origin
type DistanceZipSITOriginLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r DistanceZipSITOriginLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	originalAddress := *r.ServiceItem.SITOriginHHGOriginalAddress
	actualAddress := *r.ServiceItem.SITOriginHHGActualAddress

	// If the zip3s are identical, we do a zip3 distance calc (which uses RM).
	// If they are different, we do a zip5 distance calc (which uses DTOD).

	originZip := originalAddress.PostalCode
	if len(originZip) < 5 {
		return "", fmt.Errorf("invalid origin postal code of %s", originZip)
	}

	originZip3 := originZip[:3]

	actualOriginZip := actualAddress.PostalCode
	if len(actualOriginZip) < 5 {
		return "", fmt.Errorf("invalid SIT origin postal code of %s", actualOriginZip)
	}

	actualOriginZip3 := actualOriginZip[:3]

	var distanceMiles int
	var distanceErr error
	if originZip3 == actualOriginZip3 {
		distanceMiles, distanceErr = planner.Zip5TransitDistance(originZip, actualOriginZip)
	} else {
		distanceMiles, distanceErr = planner.Zip3TransitDistance(originZip, actualOriginZip)
	}
	if distanceErr != nil {
		return "", distanceErr
	}

	return strconv.Itoa(distanceMiles), nil
}
