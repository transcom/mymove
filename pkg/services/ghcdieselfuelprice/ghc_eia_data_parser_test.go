package ghcdieselfuelprice

import (
	"fmt"
	"testing"
)

func TestParseDieselFuelPrices(t *testing.T) {
	t.Run("parsing diesel fuel price data from data returned by FetchDieselFuelPrices function", func(t *testing.T) {
		eiaData ,_ := FetchEiaData("http://api.eia.gov/series/?api_key=3c1c9ce6bd4dcaf619f5db940d150ac6&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W")
		lastUpdated, dieselFuelPriceData := ParseEiaData(eiaData)
		fmt.Println(lastUpdated)
		fmt.Println(dieselFuelPriceData)
	})
}