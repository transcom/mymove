package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// EIAFuelPriceLookup does lookup on the ghc diesel fuel price
type EIAFuelPriceLookup struct {
	MTOShipment models.MTOShipment
}

func (r EIAFuelPriceLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Make sure there is an actual pickup date since ActualPickupDate is nullable
	actualPickupDate := r.MTOShipment.ActualPickupDate
	if actualPickupDate == nil {
		return "", fmt.Errorf("could not find actual pickup date for MTOShipment [%s]", r.MTOShipment.ID)
	}

	// Find the GHCDieselFuelPrice object with the closest prior PublicationDate to the ActualPickupDate of the MTOShipment in question
	var ghcDieselFuelPrice models.GHCDieselFuelPrice
	err := db.Where("publication_date <= ?", actualPickupDate).Order("publication_date DESC").Last(&ghcDieselFuelPrice)
	if err != nil {
		return "", services.NewNotFoundError(uuid.Nil, "Looking for GHCDieselFuelPrice")
	}

	value := fmt.Sprintf("%d", ghcDieselFuelPrice.FuelPriceInMillicents.Int())

	return value, nil
}
