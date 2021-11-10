package pricing

import (
	"strconv"
	"testing"
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

	suite.T().Run("parse sheet and check csv", func(t *testing.T) {
		slice, err := parseDomesticMoveAccessorialPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseDomesticMoveAccessorialPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("domestic"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "17_5a_access_and_add_prices_domestic_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.T().Run("try parse wrong sheet index", func(t *testing.T) {
		_, err := parseDomesticMoveAccessorialPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseDomesticMoveAccessorialPrices function failed") {
			suite.Equal("parseDomesticMoveAccessorialPrices expected to process sheet 17, but received sheetIndex 16", err.Error())
		}
	})
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

	suite.T().Run("parse sheet and check csv", func(t *testing.T) {
		slice, err := parseInternationalMoveAccessorialPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseInternationalMoveAccessorialPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("international"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "17_5a_access_and_add_prices_international_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.T().Run("try parse wrong sheet index", func(t *testing.T) {
		_, err := parseInternationalMoveAccessorialPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseInternationalMoveAccessorialPrices function failed") {
			suite.Equal("parseInternationalMoveAccessorialPrices expected to process sheet 17, but received sheetIndex 16", err.Error())
		}
	})
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

	suite.T().Run("parse sheet and check csv", func(t *testing.T) {
		slice, err := parseDomesticInternationalAdditionalPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseDomesticInternationalAdditionalPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("additional"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "17_5a_access_and_add_prices_additional_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.T().Run("try parse wrong sheet index", func(t *testing.T) {
		_, err := parseDomesticInternationalAdditionalPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseDomesticInternationalAdditionalPrices function failed") {
			suite.Equal("parseDomesticInternationalAdditionalPrices expected to process sheet 17, but received sheetIndex 16", err.Error())
		}
	})
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

	suite.T().Run("verify good sheet", func(t *testing.T) {
		err := verifyAccessAndAddPrices(params, sheetIndex)
		suite.NoError(err, "verifyAccessAndAddPrices function failed")
	})

	suite.T().Run("verify wrong sheet", func(t *testing.T) {
		err := verifyAccessAndAddPrices(params, sheetIndex-2)
		if suite.Error(err, "verifyAccessAndAddPrices function failed") {
			suite.Equal("verifyAccessAndAddPrices expected to process sheet 17, but received sheetIndex 15", err.Error())
		}
	})
}
