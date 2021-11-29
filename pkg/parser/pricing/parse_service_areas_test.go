package pricing

import (
	"strconv"
	"time"

	"github.com/go-openapi/swag"
)

func (suite *PricingParserSuite) Test_verifyServiceAreas() {
	const sheetIndex = 4
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

	err := verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

func (suite *PricingParserSuite) Test_verifyServiceAreasWrongSheet() {
	const sheetIndex = 5
	InitDataSheetInfo()

	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: suite.xlsxFilename,
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	err := verifyServiceAreas(params, sheetIndex)
	if suite.Error(err) {
		suite.Equal("verifyServiceAreas expected to process sheet 4, but received sheetIndex 5", err.Error())
	}
}

// Test_parseDomesticServiceAreas
func (suite *PricingParserSuite) Test_parseDomesticServiceAreas() {
	const sheetIndex = 4
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

	slice, err := parseDomesticServiceAreas(suite.AppContextForTest(), params, sheetIndex)
	suite.NoError(err, "parseDomesticServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("domestic"))
	err = createCSV(suite.AppContextForTest(), outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const domesticGoldenFilename string = "4_1b_service_areas_domestic_golden.csv"
	suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
}

// Test_parseInternationalServiceAreas
func (suite *PricingParserSuite) Test_parseInternationalServiceAreas() {
	const sheetIndex = 4
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

	slice, err := parseInternationalServiceAreas(suite.AppContextForTest(), params, sheetIndex)
	suite.NoError(err, "parseInternationalServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("international"))
	err = createCSV(suite.AppContextForTest(), outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const internationalGoldenFilename string = "4_1b_service_areas_international_golden.csv"
	suite.helperTestExpectedFileOutput(internationalGoldenFilename, outputFilename)
}
