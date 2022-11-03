package pricing

import (
	"strconv"
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
	suite.Run("parse non-standard location prices", func() {
		slice, err := parseNonStandardLocnPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseNonStandardLocnPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "14_3e_non_standard_locn_prices_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("attempt to parse NonStandardLocation prices with incorrect sheet index", func() {
		_, err := parseNonStandardLocnPrices(suite.AppContextForTest(), params, 13)
		if suite.Error(err) {
			suite.Equal("parseNonStandardLocnPrices expected to process sheet 14, but received sheetIndex 13", err.Error())
		}
	})

	suite.Run("verify NonStandardLocation prices", func() {
		err := verifyNonStandardLocnPrices(params, sheetIndex)
		suite.NoError(err, "verifyNonStandardLocationPrices failed")
	})

	suite.Run("attempt to verify non-standard location prices with incorrect sheet index", func() {
		err := verifyNonStandardLocnPrices(params, 13)
		if suite.Error(err) {
			suite.Equal("verifyNonStandardLocnPrices expected to process sheet 14, but received sheetIndex 13", err.Error())
		}
	})
}
