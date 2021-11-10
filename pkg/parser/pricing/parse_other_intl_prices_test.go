package pricing

import (
	"strconv"
	"time"
)

// Test_parseOtherIntlPrices
func (suite *PricingParserSuite) Test_parseOtherIntlPrices() {
	const sheetIndex = 13
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

	slice, err := parseOtherIntlPrices(suite.AppContextForTest(), params, sheetIndex)
	suite.NoError(err, "parseOtherIntlPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(suite.AppContextForTest(), outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "13_3d_other_intl_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_parseOtherIntlPricesWrongSheet() {
	const sheetIndex = 12

	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: suite.xlsxFilename,
		XlsxSheets:   []string{"13"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	err := verifyOtherIntlPrices(params, sheetIndex)
	if suite.Error(err) {
		suite.Equal("verifyOtherIntlPrices expected to process sheet 13, but received sheetIndex 12", err.Error())
	}
}

func (suite *PricingParserSuite) Test_verifyOtherIntlPrices() {
	const sheetIndex = 13
	InitDataSheetInfo()

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

	err := verifyOtherIntlPrices(params, sheetIndex)
	suite.NoError(err, "verifyOtherIntlPrices function failed")
}

func (suite *PricingParserSuite) Test_verifyOtherIntlPricesWrongSheet() {
	const sheetIndex = 12
	InitDataSheetInfo()

	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: suite.xlsxFilename,
		XlsxSheets:   []string{"13"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	err := verifyOtherIntlPrices(params, sheetIndex)
	if suite.Error(err) {
		suite.Equal("verifyOtherIntlPrices expected to process sheet 13, but received sheetIndex 12", err.Error())
	}
}
