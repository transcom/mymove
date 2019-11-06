package pricing

import (
	"time"

	"github.com/go-openapi/swag"
)

func (suite *PricingParserSuite) Test_verifyServiceAreas() {
	InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 4

	err := verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

func (suite *PricingParserSuite) Test_verifyServiceAreasWrongSheet() {
	InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"5"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 4

	err := verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

// Test_parseDomesticServiceAreas
func (suite *PricingParserSuite) Test_parseDomesticServiceAreas() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 4
	dataSheet := xlsxDataSheets[sheetIndex]

	slice, err := parseDomesticServiceAreas(params, sheetIndex)
	suite.NoError(err, "parseDomesticServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("domestic"))
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const domesticGoldenFilename string = "4_1b_service_areas_domestic_golden.csv"
	suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
}

// Test_parseInternationalServiceAreas
func (suite *PricingParserSuite) Test_parseInternationalServiceAreas() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 4
	dataSheet := xlsxDataSheets[sheetIndex]

	slice, err := parseInternationalServiceAreas(params, sheetIndex)
	suite.NoError(err, "parseInternationalServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("international"))
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const internationalGoldenFilename string = "4_1b_service_areas_international_golden.csv"
	suite.helperTestExpectedFileOutput(internationalGoldenFilename, outputFilename)
}
