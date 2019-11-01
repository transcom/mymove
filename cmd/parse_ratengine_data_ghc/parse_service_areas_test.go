package main

import (
	"log"
	"time"

	"github.com/tealeg/xlsx"
)

func (suite *ParseRateEngineGHCXLSXSuite) Test_verifyServiceAreas() {
	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: stringPointer("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"4"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	const sheetIndex int = 4

	err = verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_verifyServiceAreasWrongSheet() {
	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: stringPointer("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"5"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	const sheetIndex int = 4

	err = verifyServiceAreas(params, sheetIndex)
	suite.NoError(err, "verifyServiceAreas function failed")
}

// Test_parseDomesticServiceAreas
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticServiceAreas() {
	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: stringPointer("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"4"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	const sheetIndex int = 4
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime, stringPointer("domestic"))
	if csvWriter != nil {
		defer csvWriter.close()
	}

	err = parseDomesticServiceAreas(params, sheetIndex, suite.tableFromSliceCreator, csvWriter)
	suite.NoError(err, "parseDomesticServiceAreas function failed")

	outputFilename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.runTime, stringPointer("domestic"))

	const domesticGoldenFilename string = "4_1b_service_areas_domestic_golden.csv"
	suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
}

// Test_parseInternationalServiceAreas
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseInternationalServiceAreas() {
	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: stringPointer("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"4"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	const sheetIndex int = 4
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime, stringPointer("international"))
	if csvWriter != nil {
		defer csvWriter.close()
	}

	err = parseInternationalServiceAreas(params, sheetIndex, suite.tableFromSliceCreator, csvWriter)
	suite.NoError(err, "parseInternationalServiceAreas function failed")

	outputFilename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.runTime, stringPointer("international"))

	const internationalGoldenFilename string = "4_1b_service_areas_international_golden.csv"
	suite.helperTestExpectedFileOutput(internationalGoldenFilename, outputFilename)
}
