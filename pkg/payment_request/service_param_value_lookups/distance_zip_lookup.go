package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// DistanceZipLookup contains zip3 lookup
type DistanceZipLookup struct {
	PickupAddress      models.Address
	DestinationAddress models.Address
}

func (r DistanceZipLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	planner := keyData.planner
	db := appCtx.DB()
	var distanceMiles int
	var totalDistanceMiles int
	hasApprovedDestinationSIT := false

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

	err = appCtx.DB().EagerPreload(
		"MTOServiceItems",
		"Distance",
		"PickupAddress",
		"DestinationAddress",
	).Find(&mtoShipment, mtoShipment.ID)
	if err != nil {
		return "", err
	}

	// Now calculate the distance between zips
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode

	// if the shipment is international, we need to change the respective ZIP to use the port ZIP and not the address ZIP
	if mtoShipment.MarketCode == models.MarketCodeInternational {
		portZip, portType, err := models.GetPortLocationInfoForShipment(appCtx.DB(), *mtoShipmentID)
		if err != nil {
			return "", err
		}
		if portZip != nil && portType != nil {
			// if the port type is POEFSC this means the shipment is CONUS -> OCONUS (pickup -> port)
			// if the port type is PODFSC this means the shipment is OCONUS -> CONUS (port -> destination)
			if *portType == models.ReServiceCodePOEFSC.String() {
				destinationZip = *portZip
			} else if *portType == models.ReServiceCodePODFSC.String() {
				pickupZip = *portZip
			}
		} else {
			return "", apperror.NewNotFoundError(*mtoShipmentID, "looking for port ZIP for shipment")
		}
	}

	// if the shipment is international, we need to change the respective ZIP to use the port ZIP and not the address ZIP
	if mtoShipment.MarketCode == models.MarketCodeInternational {
		portZip, portType, err := models.GetPortLocationInfoForShipment(appCtx.DB(), *mtoShipmentID)
		if err != nil {
			return "", err
		}
		if portZip != nil && portType != nil {
			// if the port type is POEFSC this means the shipment is CONUS -> OCONUS (pickup -> port)
			// if the port type is PODFSC this means the shipment is OCONUS -> CONUS (port -> destination)
			if *portType == models.ReServiceCodePOEFSC.String() {
				destinationZip = *portZip
			} else if *portType == models.ReServiceCodePODFSC.String() {
				pickupZip = *portZip
			}
		} else {
			return "", apperror.NewNotFoundError(*mtoShipmentID, "looking for port ZIP for shipment")
		}
	}
	errorMsgForPickupZip := fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip)
	errorMsgForDestinationZip := fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip)
	if len(pickupZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("%s", errorMsgForPickupZip), nil, errorMsgForPickupZip)
	}
	if len(destinationZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("%s", errorMsgForDestinationZip), nil, errorMsgForDestinationZip)
	}

	for _, si := range mtoShipment.MTOServiceItems {
		siCopy := si
		err := appCtx.DB().EagerPreload("ReService").Find(&siCopy, siCopy.ID)
		if err != nil {
			return "", err
		}

		switch siCopy.ReService.Code {
		case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC:
			if siCopy.Status == models.MTOServiceItemStatusApproved {
				hasApprovedDestinationSIT = true
			}
		}
	}

	internationalShipment := mtoShipment.MarketCode == models.MarketCodeInternational
	if mtoShipment.Distance != nil && mtoShipment.ShipmentType != models.MTOShipmentTypePPM && !internationalShipment {
		return strconv.Itoa(mtoShipment.Distance.Int()), nil
	}

	if pickupZip == destinationZip {
		distanceMiles = 1
		totalDistanceMiles = distanceMiles
	} else if hasApprovedDestinationSIT {
		// from pickup zip to delivery zip
		totalDistanceMiles, err = planner.ZipTransitDistance(appCtx, mtoShipment.PickupAddress.PostalCode, mtoShipment.DestinationAddress.PostalCode, false, internationalShipment)
		if err != nil {
			return "", err
		}
		// from pickup zip to Destination SIT zip
		distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, destinationZip, false, internationalShipment)
		if err != nil {
			return "", err
		}
	} else {
		distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, destinationZip, false, internationalShipment)
		if err != nil {
			return "", err
		}
		totalDistanceMiles = distanceMiles
	}

	miles := unit.Miles(totalDistanceMiles)
	if mtoShipment.Distance == nil || mtoShipment.Distance.Int() != totalDistanceMiles {
		mtoShipment.Distance = &miles
		err = db.Save(&mtoShipment)
		if err != nil {
			return "", err
		}
	}

	return strconv.Itoa(distanceMiles), nil
}
