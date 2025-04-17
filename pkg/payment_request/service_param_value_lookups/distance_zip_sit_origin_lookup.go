package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// DistanceZipSITOriginLookup does the lookup of distance for SIT at origin
type DistanceZipSITOriginLookup struct {
	ServiceItem models.MTOServiceItem
	MTOShipment models.MTOShipment
}

func (r DistanceZipSITOriginLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner

	originZip, err := keyData.ServiceParamValue(appCtx, models.ServiceItemParamNameZipSITOriginHHGOriginalAddress)
	if err != nil {
		return "", err
	}
	if len(originZip) < 5 {
		return "", fmt.Errorf("invalid origin postal code of %s", originZip)
	}

	var actualOriginZip string
	actualOriginZip, err = keyData.ServiceParamValue(appCtx, models.ServiceItemParamNameZipSITOriginHHGActualAddress)
	if err != nil {
		return "", err
	}
	if len(actualOriginZip) < 5 {
		return "", fmt.Errorf("invalid SIT origin postal code of %s", actualOriginZip)
	}

	var distanceMiles int
	var distanceErr error

	// Verify if original, actual and MTOShipment.PickupAddress are all the same
	if originZip == actualOriginZip && r.MTOShipment.PickupAddress.PostalCode == actualOriginZip {
		distanceMiles = 1
	} else {
		distanceMiles, distanceErr = planner.ZipTransitDistance(appCtx, r.MTOShipment.PickupAddress.PostalCode, actualOriginZip)
	}
	if distanceErr != nil {
		return "", distanceErr
	}

	return strconv.Itoa(distanceMiles), nil
}
