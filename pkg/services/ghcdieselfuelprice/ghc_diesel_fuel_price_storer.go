package ghcdieselfuelprice

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func priceInMillicents(price float64) unit.Millicents {
	priceInMillicents := unit.Millicents(int(price * 100000))

	return priceInMillicents
}

func publicationDateInTime(publicationDate string) (time.Time, error) {
	publicationDateInTime, err := time.Parse("2006-01-02", publicationDate)

	return publicationDateInTime, err
}

// RunStorer stores the final EIA weekly average diesel fuel price data in the ghc_diesel_fuel_price table
func (d *DieselFuelPriceInfo) RunStorer(appCtx appcontext.AppContext) error {
	priceInMillicents := priceInMillicents(d.dieselFuelPriceData.price)

	publicationDate, err := publicationDateInTime(d.dieselFuelPriceData.publicationDate)
	if err != nil {
		return err
	}

	var newGHCDieselFuelPrice models.GHCDieselFuelPrice

	newGHCDieselFuelPrice.PublicationDate = publicationDate
	newGHCDieselFuelPrice.FuelPriceInMillicents = priceInMillicents

	var lastGHCDieselFuelPrice models.GHCDieselFuelPrice

	err = appCtx.DB().Where("publication_date = ?", publicationDate).First(&lastGHCDieselFuelPrice)
	if err != nil {
		appCtx.Logger().Info("no existing GHCDieselFuelPrice record found with", zap.String("publication_date", publicationDate.String()))

		verrs, err := appCtx.DB().ValidateAndCreate(&newGHCDieselFuelPrice)
		if err != nil {
			return fmt.Errorf("failed to create ghcDieselFuelPrice: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ghcDieselFuelPrice: %w", verrs)
		}
	} else if priceInMillicents != lastGHCDieselFuelPrice.FuelPriceInMillicents {
		appCtx.Logger().Info("Updating existing GHCDieselFuelPrice record found with", zap.String("publication_date", publicationDate.String()))
		lastGHCDieselFuelPrice.FuelPriceInMillicents = priceInMillicents

		verrs, err := appCtx.DB().ValidateAndUpdate(&lastGHCDieselFuelPrice)
		if err != nil {
			return fmt.Errorf("failed to update ghcDieselFuelPrice: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ghcDieselFuelPrice: %w", verrs)
		}
	} else {
		appCtx.Logger().Info(
			"Existing GHCDieselFuelPrice record found with matching fuel prices",
			zap.String("publication_date", publicationDate.String()),
			zap.String("fuel price", priceInMillicents.ToDollarString()))
	}

	return nil
}
