package ghcdieselfuelprice

import (
	"fmt"
	"testing"
)

func TestParseEiaData(t *testing.T) {
	t.Run("parsing diesel fuel price data from EIA data returned by FetchDieselFuelPrices function", func(t *testing.T) {
		// TODO: Figure out how to handle When the series data is an empty array, drilling down into the JSON data throws an error
		eiaData ,_ := FetchEiaData("http://api.eia.gov/series/?api_key=3c1c9ce6bd4dcaf619f5db940d150ac6&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W")
		lastUpdated, dieselFuelPriceData, err := ParseEiaData(eiaData)
		fmt.Println(lastUpdated)
		fmt.Println(dieselFuelPriceData)
		fmt.Println(err)
	})
}