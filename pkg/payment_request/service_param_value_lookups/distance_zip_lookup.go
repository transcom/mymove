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

	err = appCtx.DB().EagerPreload("DeliveryAddressUpdate", "DeliveryAddressUpdate.OriginalAddress", "DeliveryAddressUpdate.NewAddress", "MTOServiceItems", "Distance").Find(&mtoShipment, mtoShipment.ID)
	if err != nil {
		return "", err
	}

	// Now calculate the distance between zips
	pickupZip := r.PickupAddress.PostalCode
	destinationZip := r.DestinationAddress.PostalCode
	errorMsgForPickupZip := fmt.Sprintf("Shipment must have valid pickup zipcode. Received: %s", pickupZip)
	errorMsgForDestinationZip := fmt.Sprintf("Shipment must have valid destination zipcode. Received: %s", destinationZip)
	if len(pickupZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("%s", errorMsgForPickupZip), nil, errorMsgForPickupZip)
	}
	if len(destinationZip) < 5 {
		return "", apperror.NewInvalidInputError(*mtoShipmentID, fmt.Errorf("%s", errorMsgForDestinationZip), nil, errorMsgForDestinationZip)
	}

	serviceCode := keyData.MTOServiceItem.ReService.Code
	switch serviceCode {
	case models.ReServiceCodeDLH, models.ReServiceCodeDSH, models.ReServiceCodeFSC:
		for _, si := range mtoShipment.MTOServiceItems {
			siCopy := si
			err := appCtx.DB().EagerPreload("ReService", "ApprovedAt").Find(&siCopy, siCopy.ID)
			if err != nil {
				return "", err
			}

			switch siCopy.ReService.Code {
			case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC:
				if mtoShipment.DeliveryAddressUpdate != nil && mtoShipment.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
					if siCopy.ApprovedAt != nil {
						if mtoShipment.DeliveryAddressUpdate.UpdatedAt.After(*siCopy.ApprovedAt) {
							destinationZip = mtoShipment.DeliveryAddressUpdate.OriginalAddress.PostalCode
						} else {
							destinationZip = mtoShipment.DeliveryAddressUpdate.NewAddress.PostalCode
						}
					}
				}
			}
		}

		if mtoShipment.DeliveryAddressUpdate != nil && mtoShipment.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
			distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, mtoShipment.DeliveryAddressUpdate.NewAddress.PostalCode, false)
			if err != nil {
				return "", err
			}
			return strconv.Itoa(distanceMiles), nil
		}
	}

	if mtoShipment.Distance != nil && mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
		return strconv.Itoa(mtoShipment.Distance.Int()), nil
	}

	if pickupZip == destinationZip {
		distanceMiles = 1
	} else {
		distanceMiles, err = planner.ZipTransitDistance(appCtx, pickupZip, destinationZip, false)
		if err != nil {
			return "", err
		}
	}

	miles := unit.Miles(distanceMiles)
	mtoShipment.Distance = &miles
	err = db.Save(&mtoShipment)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(distanceMiles), nil
}
