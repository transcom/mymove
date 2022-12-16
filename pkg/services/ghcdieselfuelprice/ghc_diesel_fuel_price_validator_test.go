package ghcdieselfuelprice

import "fmt"

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceValidator() {
	type validatorScenario struct {
		name     string
		errorMsg string
	}
	validatorScenarios := []validatorScenario{
		{"empty fuel data", "received empty array of fuel data"},
		{"invalid duo area", "Expected DuoArea to be NUS, received INVALID"},
		{"invalid area name", "Expected AreaName to be U.S., received INVALID"},
		{"invalid product", "Expected Product to be EPD2D, received INVALID"},
		{"invalid process", "Expected Process to be PTE, received INVALID"},
		{"invalid series", "Expected Series to be EMD_EPD2D_PTE_NUS_DPG, received INVALID"},
		{"invalid units", "Expected Units to be $/GAL, received INVALID"},
		{"invalid date format", "Expected DateFormat to be YYYY-MM-DD, received INVALID"},
		{"invalid frequency", "Expected Frequency to be weekly, received INVALID"},
	}

	for _, scenario := range validatorScenarios {
		suite.Run(fmt.Sprintf("validation scenario %s", scenario.name), func() {
			newDieselFuelPriceInfo := NewDieselFuelPriceInfo(scenario.name, "", suite.helperStubEIAData, suite.Logger())

			eiaData, err := newDieselFuelPriceInfo.eiaDataFetcherFunction(newDieselFuelPriceInfo.eiaURL)
			suite.NoError(err)

			verr := eiaData.validateEIAData()
			suite.Error(verr)
			suite.Equal(scenario.errorMsg, verr.Error())
		})
	}
}
