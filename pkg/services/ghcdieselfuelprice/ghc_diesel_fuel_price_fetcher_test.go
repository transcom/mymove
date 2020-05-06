package ghcdieselfuelprice

func (suite *GhcDieselFuelPriceServiceSuite) TestFetchDieselFuelPrice() {
	_ ,err := FetchDieselFuelPrice("https://api.eia.gov/series/?api_key=3c1c9ce6bd4dcaf619f5db940d150ac6&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.M")

	suite.NoError(err)
}