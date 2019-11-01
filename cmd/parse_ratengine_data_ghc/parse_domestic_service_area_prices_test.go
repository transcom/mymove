package main

import (
	"log"
	"time"

	"github.com/tealeg/xlsx"
)

// Test_parseDomesticServiceAreaPrices
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticServiceAreaPrices() {

	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: stringPointer("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"7"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	const sheetIndex int = 7
	err = parseDomesticServiceAreaPrices(params, sheetIndex, suite.tableFromSliceCreator)
	suite.NoError(err, "parseDomesticServiceAreaPrices function failed")

	outputFilename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.runTime, nil)

	const goldenFilename string = "7_2b_domestic_service_area_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)

}
