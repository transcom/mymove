package pricing

import (
	"strconv"
	"time"
)

// Test_parseOconusToOconusPrices
func (suite *PricingParserSuite) Test_parseOconusToOconusPrices() {
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

	slice, err := parseOconusToOconusPrices(params, sheetIndex)
	suite.NoError(err, "parseOconusToOconusPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "10_3a_oconus_to_oconus_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_parseConusToOconusPrices() {
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

	slice, err := parseConusToOconusPrices(params, sheetIndex)
	suite.NoError(err, "parseConusToOconusPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "11_3b_conus_to_oconus_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_parseOconusToConusPrices() {
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

	slice, err := parseOconusToConusPrices(params, sheetIndex)
	suite.NoError(err, "parseOconusToConusPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "12_3c_oconus_to_conus_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
