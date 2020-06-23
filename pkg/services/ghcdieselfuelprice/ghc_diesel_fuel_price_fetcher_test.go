package ghcdieselfuelprice

import (
	"testing"
)

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceFetcher() {
	suite.T().Run("build correct EIA Open Data API URL", func(t *testing.T) {
		newDieselFuelPriceInfo := newDieselFuelPriceInfo("https://api.eia.gov/series/", "pUW34B2q8tLooWEVQpU7s9Joq672q2rP", suite.logger, suite.helperStubEIAData)  // eiaKey: "pUW34B2q8tLooWEVQpU7s9Joq672q2rP" is a fake key

		finalEIAAPIURL, err := buildFinalEIAAPIURL(newDieselFuelPriceInfo.eiaURL, newDieselFuelPriceInfo.eiaKey)
		suite.NoError(err)

		suite.Equal("https://api.eia.gov/series/?api_key=pUW34B2q8tLooWEVQpU7s9Joq672q2rP&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W", finalEIAAPIURL)
	})

	suite.T().Run("EIA Open Data API error - invalid or missing api_key", func(t *testing.T) {
		newDieselFuelPriceInfo := newDieselFuelPriceInfo("EIA Open Data API error - invalid or missing api_key", "", suite.logger, suite.helperStubEIAData)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("EIA Open Data API error - invalid series_id", func(t *testing.T) {
		newDieselFuelPriceInfo := newDieselFuelPriceInfo("EIA Open Data API error - invalid series_id", "", suite.logger, suite.helperStubEIAData)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("nil series data", func(t *testing.T) {
		newDieselFuelPriceInfo := newDieselFuelPriceInfo("nil series data", "", suite.logger, suite.helperStubEIAData)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.Error(err)
	})

	suite.T().Run("extract diesel fuel price data", func(t *testing.T) {
		newDieselFuelPriceInfo := newDieselFuelPriceInfo("extract diesel fuel price data", "", suite.logger, suite.helperStubEIAData)

		eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
		suite.NoError(err)

		_, err = extractDieselFuelPriceData(eiaData)
		suite.NoError(err)
	})
}
