package pricing

import (
	"strconv"
	"time"
)

// Test_parseDomesticLinehaulPrices
func (suite *PricingParserSuite) Test_parseDomesticLinehaulPrices() {
	const sheetIndex = 6
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

	slice, err := parseDomesticLinehaulPrices(suite.AppContextForTest(), params, sheetIndex)
	suite.NoError(err, "parseDomesticLinehaulPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(suite.AppContextForTest(), outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "6_2a_domestic_linehaul_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
