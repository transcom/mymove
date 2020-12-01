package ghcdieselfuelprice

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func priceInMillicents(price float64) unit.Millicents {
	priceInMillicents := unit.Millicents(int(price * 100000))

	return priceInMillicents
}

func publicationDateInTime(publicationDate string) (time.Time, error) {
	publicationDateInTime, err := time.Parse("20060102", publicationDate)

	return publicationDateInTime, err
}

// RunStorer stores the final EIA weekly average diesel fuel price data in the ghc_diesel_fuel_price table
func (d *DieselFuelPriceInfo) RunStorer(dbTx *pop.Connection) error {
	priceInMillicents := priceInMillicents(d.dieselFuelPriceData.price)

	publicationDate, err := publicationDateInTime(d.dieselFuelPriceData.publicationDate)
	if err != nil {
		return err
	}

	var newGHCDieselFuelPrice models.GHCDieselFuelPrice

	newGHCDieselFuelPrice.PublicationDate = publicationDate
	newGHCDieselFuelPrice.FuelPriceInMillicents = priceInMillicents

	var lastGHCDieselFuelPrice models.GHCDieselFuelPrice

	err = dbTx.Where("publication_date = ?", publicationDate).First(&lastGHCDieselFuelPrice)
	if err != nil {
		d.logger.Info("no existing GHCDieselFuelPrice record found with", zap.String("publication_date", publicationDate.String()))

		verrs, err := dbTx.ValidateAndCreate(&newGHCDieselFuelPrice)
		if err != nil {
			return fmt.Errorf("failed to create ghcDieselFuelPrice: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ghcDieselFuelPrice: %w", verrs)
		}
	} else if priceInMillicents != lastGHCDieselFuelPrice.FuelPriceInMillicents {
		lastGHCDieselFuelPrice.FuelPriceInMillicents = priceInMillicents

		verrs, err := dbTx.ValidateAndUpdate(&lastGHCDieselFuelPrice)
		if err != nil {
			return fmt.Errorf("failed to update ghcDieselFuelPrice: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ghcDieselFuelPrice: %w", verrs)
		}
	}

	return nil
}