package ghcdieselfuelprice

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceFetcher() {
	suite.Run("build correct EIA Open Data API URL", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("https://api.eia.gov/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W", "pUW34B2q8tLooWEVQpU7s9Joq672q2rP", suite.helperStubEIAData, suite.Logger()) // eiaKey: "pUW34B2q8tLooWEVQpU7s9Joq672q2rP" is a fake key

		finalEIAAPIURL, err := buildFinalEIAAPIURL(newDieselFuelPriceInfo.eiaURL, newDieselFuelPriceInfo.eiaKey)
		suite.NoError(err)

		suite.Equal("https://api.eia.gov/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W?api_key=pUW34B2q8tLooWEVQpU7s9Joq672q2rP", finalEIAAPIURL)
	})

	suite.Run("pass empty EIA Open Data API URL to Fetch data function", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("empty url", "", FetchEIAData, suite.Logger())

		_, err := newDieselFuelPriceInfo.eiaDataFetcherFunction("")
		suite.Error(err)

		suite.Equal(
			"expected finalEIAAPIURL to contain EIA Open Data API request URL, but got empty string",
			err.Error())
	})

	suite.Run("EIA Open Data API error - invalid or missing api_key", func() {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("EIA Open Data API error - invalid or missing api_key", "", suite.helperStubEIAData, suite.Logger())

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		err = checkResponseForErrors(eiaData)
		suite.Error(err)
		suite.Equal(
			"received an error from the EIA Open Data API: API_KEY_MISSING No api_key was supplied.  Please register for one at https://www.eia.gov/opendata/register.php",
			err.Error())
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

		suite.Equal("2022-12-12", newDieselFuelPriceInfo.dieselFuelPriceData.publicationDate)
		suite.Equal(4.759, newDieselFuelPriceInfo.dieselFuelPriceData.price)
	})
}
