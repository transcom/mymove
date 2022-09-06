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

func (r EIAFuelPriceLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	db := appCtx.DB()

	// Make sure there is an actual pickup date since ActualPickupDate is nullable
	actualPickupDate := r.MTOShipment.ActualPickupDate
	if actualPickupDate == nil {
		return "", fmt.Errorf("could not find actual pickup date for MTOShipment [%s]", r.MTOShipment.ID)
	}

	// Find the GHCDieselFuelPrice object with the closest prior PublicationDate to the ActualPickupDate of the MTOShipment in question
	var ghcDieselFuelPrice models.GHCDieselFuelPrice
	err := db.Where("publication_date <= ?", actualPickupDate).Order("publication_date DESC").Last(&ghcDieselFuelPrice)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(uuid.Nil, "Looking for GHCDieselFuelPrice")
		default:
			return "", apperror.NewQueryError("GHCDieselFuelPrice", err, "")
		}
	}

	value := fmt.Sprintf("%d", ghcDieselFuelPrice.FuelPriceInMillicents.Int())

	return value, nil
}
