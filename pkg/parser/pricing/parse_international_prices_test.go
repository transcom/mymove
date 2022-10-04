package pricing

import (
	"strconv"
	"time"
)

// Test_parseOconusToOconusPrices
func (suite *PricingParserSuite) Test_OconusToOconusPrices() {
	const sheetIndex = 10
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
	suite.Run("parse oconusToOconus prices", func() {
		slice, err := parseOconusToOconusPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseOconusToOconusPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "10_3a_oconus_to_oconus_prices_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("attempt to parse oconusToOconus prices with incorrect sheet index", func() {
		_, err := parseOconusToOconusPrices(suite.AppContextForTest(), params, 7)
		if suite.Error(err) {
			suite.Equal("parseOconusToOconusPrices expected to process sheet 10, but received sheetIndex 7", err.Error())
		}
	})

	suite.Run("verify oconusToOconus prices", func() {
		err := verifyIntlOconusToOconusPrices(params, sheetIndex)
		suite.NoError(err, "verifyIntlOconusToOconusPrices failed")
	})

	suite.Run("attempt to verify oconusToOconus prices with incorrect sheet index", func() {
		err := verifyIntlOconusToOconusPrices(params, 7)
		if suite.Error(err) {
			suite.Equal("verifyInternationalPrices expected to process sheet 10, but received sheetIndex 7", err.Error())
		}
	})
}

// Test_parseConusToOconusPrices
func (suite *PricingParserSuite) Test_ConusToOconusPrices() {
	const sheetIndex = 11
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
	suite.Run("parse conusToOconus prices", func() {
		slice, err := parseConusToOconusPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseConusToOconusPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "11_3b_conus_to_oconus_prices_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("attempt to parse conusToOconus prices with incorrect sheet index", func() {
		_, err := parseConusToOconusPrices(suite.AppContextForTest(), params, 7)
		if suite.Error(err) {
			suite.Equal("parseConusToOconusPrices expected to process sheet 11, but received sheetIndex 7", err.Error())
		}
	})

	suite.Run("verify conusToOconus prices", func() {
		err := verifyIntlConusToOconusPrices(params, sheetIndex)
		suite.NoError(err, "verifyIntlConusToOconusPrices failed")
	})

	suite.Run("attempt to verify conusToOconus prices with incorrect sheet index", func() {
		err := verifyIntlConusToOconusPrices(params, 7)
		if suite.Error(err) {
			suite.Equal("verifyInternationalPrices expected to process sheet 11, but received sheetIndex 7", err.Error())
		}
	})
}

// Test_parseOconusToConusPrices
func (suite *PricingParserSuite) Test_OconusToConusPrices() {
	const sheetIndex = 12
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
	suite.Run("parse oconusToConus prices", func() {
		slice, err := parseOconusToConusPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseOconusToConusPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "12_3c_oconus_to_conus_prices_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("attempt to parse oconusToConus prices with incorrect sheet index", func() {
		_, err := parseOconusToConusPrices(suite.AppContextForTest(), params, 7)
		if suite.Error(err) {
			suite.Equal("parseOconusToConusPrices expected to process sheet 12, but received sheetIndex 7", err.Error())
		}
	})

	suite.Run("verify oconusToConus prices", func() {
		err := verifyIntlOconusToConusPrices(params, sheetIndex)
		suite.NoError(err, "verifyIntlOconusToConusPrices failed")
	})

	suite.Run("attempt to verify oconusToConus prices with incorrect sheet index", func() {
		err := verifyIntlOconusToConusPrices(params, 7)
		if suite.Error(err) {
			suite.Equal("verifyInternationalPrices expected to process sheet 12, but received sheetIndex 7", err.Error())
		}
	})
}
