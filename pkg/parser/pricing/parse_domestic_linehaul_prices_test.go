package pricing

import (
	"log"
	"time"

	"github.com/go-openapi/swag"
	"github.com/tealeg/xlsx"
)

// Test_parseDomesticLinehaulPrices
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticLinehaulPrices() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: swag.String("fixtures/pricing_template_2019-09-19_fake-data.xlsx"),
		XlsxSheets:   []string{"6"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		RunVerify:    true,
	}

	const sheetIndex int = 6
	dataSheet := xlsxDataSheets[sheetIndex]

	xlsxFile, err := xlsx.OpenFile(*params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.XlsxFilename, err)
	}

	slice, err := parseDomesticLinehaulPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticLinehaulPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "6_2a_domestic_linehaul_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
