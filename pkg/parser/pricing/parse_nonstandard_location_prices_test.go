package pricing

import (
	"strconv"
	"testing"
	"time"
)

// Test_parseNonStandardLocnPrices
func (suite *PricingParserSuite) Test_NonStandardLocnPrices() {
	const sheetIndex = 14
	xlsxDataSheets := InitDataSheetInfo()
	dataSheet := xlsxDataSheets[sheetIndex]

	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: suite.xlsxFilename,
		XlsxSheets:   []string{strconv.Itoa(sheetIndex)},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}
	suite.T().Run("parse non-standard location prices", func(t *testing.T) {
		slice, err := parseNonStandardLocnPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseNonStandardLocnPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "14_3e_non_standard_locn_prices_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.T().Run("attempt to parse NonStandardLocation prices with incorrect sheet index", func(t *testing.T) {
		_, err := parseNonStandardLocnPrices(suite.AppContextForTest(), params, 13)
		if suite.Error(err) {
			suite.Equal("parseNonStandardLocnPrices expected to process sheet 14, but received sheetIndex 13", err.Error())
		}
	})

	suite.T().Run("verify NonStandardLocation prices", func(t *testing.T) {
		err := verifyNonStandardLocnPrices(params, sheetIndex)
		suite.NoError(err, "verifyNonStandardLocationPrices failed")
	})

	suite.T().Run("attempt to verify non-standard location prices with incorrect sheet index", func(t *testing.T) {
		err := verifyNonStandardLocnPrices(params, 13)
		if suite.Error(err) {
			suite.Equal("verifyNonStandardLocnPrices expected to process sheet 14, but received sheetIndex 13", err.Error())
		}
	})
}
