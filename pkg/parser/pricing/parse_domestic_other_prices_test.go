package pricing

import (
	"strconv"
	"time"
)

// Test_parseDomesticOtherPricesPack
func (suite *PricingParserSuite) Test_parseDomesticOtherPrices() {
	const sheetIndex = 8
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

	suite.Run("parseDomesticOtherPricesPack", func() {
		slice, err := parseDomesticOtherPricesPack(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseDomesticOtherPricesPack function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "8_2c_domestic_other_prices_pack_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("parseDomesticOtherPricesSit", func() {
		slice, err := parseDomesticOtherPricesSit(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseDomesticOtherPricesSit function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "8_2c_domestic_other_prices_sit_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})
}

// Test_verifyDomesticOtherPrices
func (suite *PricingParserSuite) Test_verifyDomesticOtherPrices() {
	sheetIndex := 8
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

	suite.Run("verifyDomesticOtherPrices success", func() {
		err := verifyDomesticOtherPrices(params, sheetIndex)
		suite.NoError(err)
	})

	suite.Run("verifyDomesticOtherPrices with invalid sheetIndex", func() {
		err := verifyDomesticOtherPrices(params, 7)
		if suite.Error(err) {
			suite.Equal("verifyDomesticOtherPrices expected to process sheet 8, but received sheetIndex 7", err.Error())
		}
	})
}
