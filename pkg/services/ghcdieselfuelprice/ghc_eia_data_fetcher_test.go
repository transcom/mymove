package ghcdieselfuelprice

import (
	"fmt"
	"testing"
)

func TestFetchEiaData(t *testing.T) {
	// TODO: Checking the response code is 200
	// TODO: Checking that a 200 response with bad data is handled correctly

	t.Run("fetching series data from EIA Open Data API", func(t *testing.T) {
		eiaData , err := FetchEiaData("http://api.eia.gov/s1eries/?api_key=3c1c9ce6bd4dcaf619f5db940d150ac6&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W")
		if len(eiaData.SeriesData) == 0 {
			t.Error("Failed to fetch series data from EIA Open Data API")
		}
		fmt.Println(eiaData)
		fmt.Println(err)
	})
}