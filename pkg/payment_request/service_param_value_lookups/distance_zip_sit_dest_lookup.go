package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// DistanceZipSITDestLookup does the lookup of distance for SIT at destination
type DistanceZipSITDestLookup struct {
	DestinationAddress models.Address
}

func (r DistanceZipSITDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := keyData.db

	// Get the MTOServiceItem and associated MTOShipment and addresses
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.
		// Eager("SITDestinationFinalAddress").
		Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// If the zip3s are identical, we do a zip3 distance calc (which uses RM).
	// If they are different, we do a zip5 distance calc (which uses DTOD).

	destZip := r.DestinationAddress.PostalCode
	if len(destZip) < 5 {
		return "", fmt.Errorf("invalid destination postal code of %s", destZip)
	}
	destZip3 := destZip[:3]

	// sitDestZip := mtoServiceItem.SITDestinationFinalAddress.PostalCode
	sitDestZip := "30907" // Placeholder for now
	if len(sitDestZip) < 5 {
		return "", fmt.Errorf("invalid SIT destination postal code of %s", destZip)
	}
	sitDestZip3 := sitDestZip[:3]

	var distanceMiles int
	var distanceErr error
	if destZip3 == sitDestZip3 {
		distanceMiles, distanceErr = planner.Zip5TransitDistance(destZip, sitDestZip)
	} else {
		distanceMiles, distanceErr = planner.Zip3TransitDistance(destZip, sitDestZip)
	}
	if distanceErr != nil {
		return "", distanceErr
	}

	// TODO: Do we need to store the distance anywhere like the other distance lookups?

	return strconv.Itoa(distanceMiles), nil
}
