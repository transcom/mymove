package pricing

import (
	"strconv"
	"time"

	"github.com/go-openapi/swag"
)

func (suite *PricingParserSuite) Test_parseDomesticMoveAccessorialPrices() {
	const sheetIndex = 17
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

	slice, err := parseDomesticMoveAccessorialPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticMoveAccessorialPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("domestic"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "17_5a_access_and_add_prices_domestic_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_parseInternationalMoveAccessorialPrices() {
	const sheetIndex = 17
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

	slice, err := parseInternationalMoveAccessorialPrices(params, sheetIndex)
	suite.NoError(err, "parseInternationalMoveAccessorialPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("international"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "17_5a_access_and_add_prices_international_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_parseDomesticInternationalAdditionalPrices() {
	const sheetIndex = 17
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

	slice, err := parseDomesticInternationalAdditionalPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticInternationalAdditionalPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("additional"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "17_5a_access_and_add_prices_additional_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

func (suite *PricingParserSuite) Test_verifyAccessAndAddPrices() {
	const sheetIndex = 17
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

	err := verifyAccessAndAddPrices(params, sheetIndex)
	suite.NoError(err, "verifyAccessAndAddPrices function failed")
}

func (suite *PricingParserSuite) Test_verifyAccessAndAddPricesWithWrongSheet() {
	const sheetIndex = 15
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

	err := verifyAccessAndAddPrices(params, sheetIndex)
	suite.Error(err, "verifyAccessAndAddPrices function failed")
	suite.Equal("verifyAccessAndAddPrices expected to process sheet 17, but received sheetIndex 15", err.Error())
}
