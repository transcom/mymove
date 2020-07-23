package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// EIAFuelPriceLookup does lookup on the ghc diesel fuel price
type EIAFuelPriceLookup struct {
}

func (r EIAFuelPriceLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there is an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipmentID
	if mtoShipmentID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there is an actual pickup date since ActualPickupDate is nullable
	actualPickupDate := mtoServiceItem.MTOShipment.ActualPickupDate
	if actualPickupDate == nil {
		return "", fmt.Errorf("could not find actual pickup date for MTOShipment [%s]", mtoShipmentID)
	}

	// Find the GHCDieselFuelPrice object with the closest prior PublicationDate to the ActualPickupDate of the MTOShipment in question
	var ghcDieselFuelPrice models.GHCDieselFuelPrice
	err = db.Where("publication_date <= ?", actualPickupDate).Order("publication_date DESC").Last(&ghcDieselFuelPrice)
	if err != nil {
		return "", services.NewNotFoundError(uuid.Nil, "Looking for GHCDieselFuelPrice")
	}

	value := fmt.Sprintf("%d", ghcDieselFuelPrice.FuelPriceInMillicents.Int())

	return value, nil
}
