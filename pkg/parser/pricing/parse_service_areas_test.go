package pricing

import (
	"log"
	"time"

	"github.com/go-openapi/swag"
	"github.com/tealeg/xlsx"
)

func (suite *ParseRateEngineGHCXLSXSuite) Test_verifyServiceAreas() {
	InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		RunVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.XlsxFilename, err)
	}

	const sheetIndex int = 4

	err = verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_verifyServiceAreasWrongSheet() {
	InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		XlsxSheets:   []string{"5"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		RunVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.XlsxFilename, err)
	}

	const sheetIndex int = 4

	err = verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

// Test_parseDomesticServiceAreas
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticServiceAreas() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		RunVerify:    true,
	}

	const sheetIndex int = 4
	dataSheet := xlsxDataSheets[sheetIndex]

	xlsxFile, err := xlsx.OpenFile(*params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.XlsxFilename, err)
	}

	slice, err := parseDomesticServiceAreas(params, sheetIndex)
	suite.NoError(err, "parseDomesticServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("domestic"))
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const domesticGoldenFilename string = "4_1b_service_areas_domestic_golden.csv"
	suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
}

// Test_parseInternationalServiceAreas
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseInternationalServiceAreas() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		XlsxSheets:   []string{"4"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		RunVerify:    true,
	}

	const sheetIndex int = 4
	dataSheet := xlsxDataSheets[sheetIndex]

	xlsxFile, err := xlsx.OpenFile(*params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.XlsxFilename, err)
	}

	slice, err := parseInternationalServiceAreas(params, sheetIndex)
	suite.NoError(err, "parseInternationalServiceAreas function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("international"))
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const internationalGoldenFilename string = "4_1b_service_areas_international_golden.csv"
	suite.helperTestExpectedFileOutput(internationalGoldenFilename, outputFilename)
}
