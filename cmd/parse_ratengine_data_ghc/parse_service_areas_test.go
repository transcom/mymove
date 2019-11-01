package main

import (
	"log"
	"time"

	"github.com/tealeg/xlsx"
)

// Test_parseServiceAreas
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseServiceAreas() {
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
	err = parseServiceAreas(params, sheetIndex, suite.tableFromSliceCreator)
	suite.NoError(err, "parseDomesticServiceAreas function failed")

	outputFilename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.runTime, stringPointer("domestic"))

	const goldenFilename string = "4_1b_service_areas_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
