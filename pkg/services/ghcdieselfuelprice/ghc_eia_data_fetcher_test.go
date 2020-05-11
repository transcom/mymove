package ghcdieselfuelprice

import (
	"testing"
)

func TestFetchEiaData(t *testing.T) {
	t.Run("fetching series data from EIA Open Data API", func(t *testing.T) {
		eiaData ,_ := FetchEiaData("http://api.eia.gov/series/?api_key=3c1c9ce6bd4dcaf619f5db940d150ac6&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W")
		if len(eiaData.SeriesData) == 0 {
			t.Error("Failed to fetch series data from EIA Open Data API")
		}
	})
}