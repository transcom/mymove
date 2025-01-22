package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// EIAFuelPriceLookup does lookup on the ghc diesel fuel price
type EIAFuelPriceLookup struct {
	MTOShipment models.MTOShipment
}

func (r EIAFuelPriceLookup) lookup(appCtx appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	db := appCtx.DB()

	// Make sure there is an actual pickup date since ActualPickupDate is nullable
	actualPickupDate := r.MTOShipment.ActualPickupDate
	if actualPickupDate == nil {
		return "", fmt.Errorf("not found looking for shipment pickup date")
	}

	// Find the GHCDieselFuelPrice object effective before the shipment's ActualPickupDate and ends after the ActualPickupDate
	var ghcDieselFuelPrice models.GHCDieselFuelPrice
	err := db.Where("? BETWEEN effective_date and end_date", actualPickupDate).Order("publication_date DESC").First(&ghcDieselFuelPrice) //only want the first published price per week
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// If no published price is found, look for the first published price after the actual pickup date
			err = db.Where("publication_date <= ?", actualPickupDate).Order("publication_date DESC").Last(&ghcDieselFuelPrice)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return "", apperror.NewNotFoundError(uuid.Nil, "Looking for GHCDieselFuelPrice")
				default:
					return "", apperror.NewQueryError("GHCDieselFuelPrice", err, "")
				}
			}
		default:
			return "", apperror.NewQueryError("GHCDieselFuelPrice", err, "")
		}
	}

	value := fmt.Sprintf("%d", ghcDieselFuelPrice.FuelPriceInMillicents.Int())

	return value, nil
}
