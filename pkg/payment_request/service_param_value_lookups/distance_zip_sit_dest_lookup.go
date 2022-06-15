package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// DistanceZipSITDestLookup does the lookup of distance for SIT at destination
type DistanceZipSITDestLookup struct {
	DestinationAddress      models.Address
	FinalDestinationAddress models.Address
}

func (r DistanceZipSITDestLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner

	destZip := r.DestinationAddress.PostalCode
	if len(destZip) < 5 {
		return "", fmt.Errorf("invalid destination postal code of %s", destZip)
	}

	finalDestZip := r.FinalDestinationAddress.PostalCode
	if len(finalDestZip) < 5 {
		return "", fmt.Errorf("invalid SIT final destination postal code of %s", destZip)
	}

	var distanceMiles int
	var distanceErr error
	distanceMiles, distanceErr = planner.ZipTransitDistance(appCtx, destZip, finalDestZip)
	if distanceErr != nil {
		return "", distanceErr
	}

	return strconv.Itoa(distanceMiles), nil
}
