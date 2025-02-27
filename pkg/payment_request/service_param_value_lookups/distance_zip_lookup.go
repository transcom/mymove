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
		"PPMShipment.PickupAddress",
		"PPMShipment.DestinationAddress",
	).Find(&mtoShipment, mtoShipment.ID)
	if err != nil {
		return "", err
	}

	// Now calculate the distance between zips
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode

	isInternationalShipment := mtoShipment.MarketCode == models.MarketCodeInternational

	// if the shipment is international, we need to change the respective ZIP to use the port ZIP and not the address ZIP
	if isInternationalShipment {
		if mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
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
		} else {
			// PPMs get reimbursed for their travel from CONUS <-> Port ZIPs, but only for the Tacoma Port
			portLocation, err := models.FetchPortLocationByCode(appCtx.DB(), "4E1") // Tacoma port code
			if err != nil {
				return "", fmt.Errorf("unable to find port zip with code %s", "4E1")
			}
			if mtoShipment.PPMShipment != nil && mtoShipment.PPMShipment.PickupAddress != nil && mtoShipment.PPMShipment.DestinationAddress != nil {
				// need to figure out if we are going to go Port -> CONUS or CONUS -> Port
				pickupOconus := *mtoShipment.PPMShipment.PickupAddress.IsOconus
				destOconus := *mtoShipment.PPMShipment.DestinationAddress.IsOconus
				if pickupOconus && !destOconus {
					// Port ZIP -> CONUS ZIP
					pickupZip = portLocation.UsPostRegionCity.UsprZipID
					destinationZip = mtoShipment.PPMShipment.DestinationAddress.PostalCode
				} else if !pickupOconus && destOconus {
					// CONUS ZIP -> Port ZIP
					pickupZip = mtoShipment.PPMShipment.PickupAddress.PostalCode
					destinationZip = portLocation.UsPostRegionCity.UsprZipID
				} else {
					// OCONUS -> OCONUS mileage they don't get reimbursed for this
					return strconv.Itoa(0), nil
				}
			} else {
				return "", fmt.Errorf("missing required PPM & address information for shipment with id %s", mtoShipmentID)
			}
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

	if pickupZip == destinationZip {
		distanceMiles = 1
		totalDistanceMiles = distanceMiles
	} else if hasApprovedDestinationSIT {
		// from pickup zip to delivery zip
		totalDistanceMiles, err = planner.ZipTransitDistance(appCtx, mtoShipment.PickupAddress.PostalCode, mtoShipment.DestinationAddress.PostalCode)
		if err != nil {
			return "", err
		}
		// from pickup zip to Destination SIT zip
		distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, destinationZip)
		if err != nil {
			return "", err
		}
	} else {
		distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, destinationZip)
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
