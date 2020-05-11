package ghcdieselfuelprice

type DieselFuelPriceData struct {
	Date  string
	Price float64
}

func ParseEiaData(eiaData EiaData) (string, DieselFuelPriceData) {
	lastUpdated := eiaData.SeriesData[0].Updated
	var dieselFuelPriceData DieselFuelPriceData

	dieselFuelPriceData.Date = eiaData.SeriesData[0].Data[0][0].(string)
	dieselFuelPriceData.Price = eiaData.SeriesData[0].Data[0][1].(float64)

	return lastUpdated, dieselFuelPriceData
}