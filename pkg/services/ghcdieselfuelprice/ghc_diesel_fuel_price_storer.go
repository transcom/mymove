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

func publicationDateInTime(d DieselFuelPriceInfo) (time.Time, error) {
	layout := getEIADateFormatMap()[d.eiaData.ResponseData.DateFormat]
	publicationDateInTime, err := time.Parse(layout, d.dieselFuelPriceData.publicationDate)

	return publicationDateInTime, err
}

// RunStorer stores the final EIA weekly average diesel fuel price data in the ghc_diesel_fuel_price table
func (d *DieselFuelPriceInfo) RunStorer(appCtx appcontext.AppContext) error {
	priceInMillicents := priceInMillicents(d.dieselFuelPriceData.price)

	publicationDate, err := publicationDateInTime(*d)
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
		newGHCDieselFuelPrice.EffectiveDate = publicationDate.AddDate(0, 0, 1)

		dayOfWeek := publicationDate.Weekday().String()
		appCtx.Logger().Info("day_of_week", zap.String("day_of_week", dayOfWeek))
		var daysAdded int
		//fuel prices are generally published on mondays and then by business rule should expire on monday no matter what- but in case its published on a different day, we will still always expire on the following monday
		switch dayOfWeek {
		case "Monday":
			daysAdded = 7
		case "Tuesday":
			daysAdded = 6
		//very unlikely to get past here- monday is the noraml publish day- tuesday if monday is holiday.. but adding other weekedays just in case
		case "Wednesday":
			daysAdded = 6
		case "Thursday":
			daysAdded = 4
		case "Friday":
			daysAdded = 3
		}

		newGHCDieselFuelPrice.EndDate = publicationDate.AddDate(0, 0, daysAdded)
		appCtx.Logger().Info("effective_date", zap.String("effective_date", newGHCDieselFuelPrice.EffectiveDate.String()))
		appCtx.Logger().Info("end_date", zap.String("EndDate", newGHCDieselFuelPrice.EndDate.String()))

		verrs, err := appCtx.DB().ValidateAndCreate(&newGHCDieselFuelPrice)
		if err != nil {
			return fmt.Errorf("failed to create ghcDieselFuelPrice: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ghcDieselFuelPrice: %w", verrs)
		}
	} else if priceInMillicents != lastGHCDieselFuelPrice.FuelPriceInMillicents {
		appCtx.Logger().Info("existing GHCDieselFuelPrice record found with", zap.String("publication_date", publicationDate.String()))

		//no longer updating prices throughout the week- only accept the first published price per week
		if err != nil {
			return fmt.Errorf("failed to update ghcDieselFuelPrice: %w", err)
		}

	} else {
		appCtx.Logger().Info(
			"Existing GHCDieselFuelPrice record found with matching fuel prices",
			zap.String("publication_date", publicationDate.String()),
			zap.String("fuel price", priceInMillicents.ToDollarString()))
	}

	return nil
}
