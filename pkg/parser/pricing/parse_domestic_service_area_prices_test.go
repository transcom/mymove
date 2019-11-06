package pricing

import (
	"time"
)

// Test_parseDomesticServiceAreaPrices
func (suite *PricingParserSuite) Test_parseDomesticServiceAreaPrices() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"7"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 7
	dataSheet := xlsxDataSheets[sheetIndex]

	slice, err := parseDomesticServiceAreaPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticServiceAreaPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "7_2b_domestic_service_area_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
