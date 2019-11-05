package main

import (
	"log"
	"time"

	"github.com/go-openapi/swag"
	"github.com/tealeg/xlsx"
)

// Test_parseDomesticLinehaulPrices
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticLinehaulPrices() {
	initDataSheetInfo()
	params := paramConfig{
		processAll:   false,
		showOutput:   false,
		xlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		xlsxSheets:   []string{"6"},
		saveToFile:   true,
		runTime:      time.Now(),
		runVerify:    true,
	}

	const sheetIndex int = 6
	dataSheet := xlsxDataSheets[sheetIndex]

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	slice, err := parseDomesticLinehaulPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticLinehaulPrices function failed")

	err = createCSV(params, sheetIndex, dataSheet.processMethods[0], slice)
	suite.NoError(err, "could not create CSV")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.runTime, nil)

	const goldenFilename string = "6_2a_domestic_linehaul_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
