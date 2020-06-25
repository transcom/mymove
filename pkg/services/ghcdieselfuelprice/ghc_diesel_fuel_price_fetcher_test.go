package ghcdieselfuelprice

import (
	"testing"
)

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceFetcher() {
	suite.T().Run("build correct EIA Open Data API URL", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("https://api.eia.gov/series/", "pUW34B2q8tLooWEVQpU7s9Joq672q2rP", suite.helperStubEIAData, suite.logger) // eiaKey: "pUW34B2q8tLooWEVQpU7s9Joq672q2rP" is a fake key

		finalEIAAPIURL, err := buildFinalEIAAPIURL(newDieselFuelPriceInfo.eiaURL, newDieselFuelPriceInfo.eiaKey)
		suite.NoError(err)

		suite.Equal("https://api.eia.gov/series/?api_key=pUW34B2q8tLooWEVQpU7s9Joq672q2rP&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W", finalEIAAPIURL)
	})

	suite.T().Run("EIA Open Data API error - invalid or missing api_key", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("EIA Open Data API error - invalid or missing api_key", "", suite.helperStubEIAData, suite.logger)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("EIA Open Data API error - invalid series_id", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("EIA Open Data API error - invalid series_id", "", suite.helperStubEIAData, suite.logger)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("nil series data", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("nil series data", "", suite.helperStubEIAData, suite.logger)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("extract diesel fuel price data", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("extract diesel fuel price data", "", suite.helperStubEIAData, suite.logger)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.NoError(err)
	})

	suite.T().Run("run fetcher", func(t *testing.T) {
		newDieselFuelPriceInfo := NewDieselFuelPriceInfo("run fetcher", "", suite.helperStubEIAData, suite.logger)

		err := newDieselFuelPriceInfo.RunFetcher()
		suite.NoError(err)

		suite.Equal("2020-06-22T18:16:52-0400", newDieselFuelPriceInfo.dieselFuelPriceData.lastUpdated)
		suite.Equal("20200622", newDieselFuelPriceInfo.dieselFuelPriceData.publicationDate)
		suite.Equal(2.425, newDieselFuelPriceInfo.dieselFuelPriceData.price)
	})
}
