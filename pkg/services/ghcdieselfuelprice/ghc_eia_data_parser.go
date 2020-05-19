package ghcdieselfuelprice

import (
	"github.com/pkg/errors"
)

type DieselFuelPriceData struct {
	PublishedDate  string
	Price          float64
}

func ParseEiaData(eiaData EiaData) (string, DieselFuelPriceData, error) {
	lastUpdated := eiaData.SeriesData[0].Updated
	var dieselFuelPriceData DieselFuelPriceData

	publishedDate, ok := eiaData.SeriesData[0].Data[0][0].(string)
	if !ok {
		return lastUpdated, dieselFuelPriceData, errors.New("Published date returned from EIA Open Data API failed string type assertion")
	}
	dieselFuelPriceData.PublishedDate = publishedDate

	price, ok := eiaData.SeriesData[0].Data[0][1].(float64)
	if !ok {
		return lastUpdated, dieselFuelPriceData, errors.New("Price returned from EIA Open Data API failed float64 type assertion")
	}
	dieselFuelPriceData.Price = price

	return lastUpdated, dieselFuelPriceData, errors.New("Unable to parse EIA data from EIA Open Data API")
}