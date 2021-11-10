package pricing

import (
	"strconv"
	"time"
)

// Test_parseDomesticServiceAreaPrices
func (suite *PricingParserSuite) Test_parseDomesticServiceAreaPrices() {
	const sheetIndex = 7
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

	slice, err := parseDomesticServiceAreaPrices(suite.AppContextForTest(), params, sheetIndex)
	suite.NoError(err, "parseDomesticServiceAreaPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(suite.AppContextForTest(), outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "7_2b_domestic_service_area_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
