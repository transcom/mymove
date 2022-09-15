package ghcdieselfuelprice

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceFetcher() {
	suite.Run("build correct EIA Open Data API URL", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("https://api.eia.gov/series/", "pUW34B2q8tLooWEVQpU7s9Joq672q2rP", suite.helperStubEIAData, suite.Logger()) // eiaKey: "pUW34B2q8tLooWEVQpU7s9Joq672q2rP" is a fake key

		finalEIAAPIURL, err := buildFinalEIAAPIURL(newDieselFuelPriceInfo.eiaURL, newDieselFuelPriceInfo.eiaKey)
		suite.NoError(err)

		suite.Equal("https://api.eia.gov/series/?api_key=pUW34B2q8tLooWEVQpU7s9Joq672q2rP&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W", finalEIAAPIURL)
	})

	suite.Run("EIA Open Data API error - invalid or missing api_key", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("EIA Open Data API error - invalid or missing api_key", "", suite.helperStubEIAData, suite.Logger())

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.Run("EIA Open Data API error - invalid series_id", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("EIA Open Data API error - invalid series_id", "", suite.helperStubEIAData, suite.Logger())

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.Run("nil series data", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("nil series data", "", suite.helperStubEIAData, suite.Logger())

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.Run("extract diesel fuel price data", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("extract diesel fuel price data", "", suite.helperStubEIAData, suite.Logger())

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.NoError(err)
	})

	suite.Run("run fetcher", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("run fetcher", "", suite.helperStubEIAData, suite.Logger())

		err := newDieselFuelPriceInfo.RunFetcher(suite.AppContextForTest())
		suite.NoError(err)

		suite.Equal("2020-06-22T18:16:52-0400", newDieselFuelPriceInfo.dieselFuelPriceData.lastUpdated)
		suite.Equal("20200622", newDieselFuelPriceInfo.dieselFuelPriceData.publicationDate)
		suite.Equal(2.425, newDieselFuelPriceInfo.dieselFuelPriceData.price)
	})
}
