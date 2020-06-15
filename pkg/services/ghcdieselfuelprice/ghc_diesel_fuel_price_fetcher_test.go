package ghcdieselfuelprice

import (
	"testing"
)

func (suite *GhcDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceFetcher() {
	suite.T().Run("build correct EIA Open Data API URL", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("https://api.eia.gov/series/", "pUW34B2q8tLooWEVQpU7s9Joq672q2rP", suite.helperStubEiaData)

		eiaFinalURL, err := BuildEiaAPIURL(dieselFuelPriceStorer.eiaURL, dieselFuelPriceStorer.eiaKey)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaFinalURL = eiaFinalURL

		suite.Equal("https://api.eia.gov/series/?api_key=pUW34B2q8tLooWEVQpU7s9Joq672q2rP&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W", dieselFuelPriceStorer.eiaFinalURL)
	})

	// TODO: Figure out how to test FetchEiaData function without needing to make API call
	suite.T().Run("fetch EIA data from EIA Open Data API", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("https://api.eia.gov/series/", "3c1c9ce6bd4dcaf619f5db940d150ac6", FetchEiaData)

		eiaFinalURL, err := BuildEiaAPIURL(dieselFuelPriceStorer.eiaURL, dieselFuelPriceStorer.eiaKey)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaFinalURL = eiaFinalURL

		_, err = dieselFuelPriceStorer.eiaDataFetcherFunction(dieselFuelPriceStorer.eiaFinalURL)
		suite.NoError(err)
	})

	suite.T().Run("EIA Open Data API error - invalid or missing api_key", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("EIA Open Data API error - invalid or missing api_key", "", suite.helperStubEiaData)

		eiaData, err := dieselFuelPriceStorer.eiaDataFetcherFunction(dieselFuelPriceStorer.eiaURL)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaData = eiaData

		_, err = ExtractDieselFuelPriceData(dieselFuelPriceStorer.eiaData)
		suite.Error(err)

	})

	suite.T().Run("EIA Open Data API error - invalid series_id", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("EIA Open Data API error - invalid series_id", "", suite.helperStubEiaData)

		eiaData, err := dieselFuelPriceStorer.eiaDataFetcherFunction(dieselFuelPriceStorer.eiaURL)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaData = eiaData

		_, err = ExtractDieselFuelPriceData(dieselFuelPriceStorer.eiaData)
		suite.Error(err)
	})

	suite.T().Run("nil series data", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("nil series data", "", suite.helperStubEiaData)

		eiaData, err := dieselFuelPriceStorer.eiaDataFetcherFunction(dieselFuelPriceStorer.eiaURL)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaData = eiaData

		_, err = ExtractDieselFuelPriceData(dieselFuelPriceStorer.eiaData)
		suite.Error(err)
	})

	suite.T().Run("extract diesel fuel price data", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer("extract diesel fuel price data", "", suite.helperStubEiaData)

		eiaData, err := dieselFuelPriceStorer.eiaDataFetcherFunction(dieselFuelPriceStorer.eiaURL)
		suite.NoError(err)
		dieselFuelPriceStorer.eiaData = eiaData

		_, err = ExtractDieselFuelPriceData(dieselFuelPriceStorer.eiaData)
		suite.NoError(err)
	})
}
