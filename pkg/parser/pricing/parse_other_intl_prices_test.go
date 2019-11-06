package pricing

import (
	"strconv"
	"time"
)

// Test_parseOtherIntlPrices
func (suite *PricingParserSuite) Test_parseOtherIntlPrices() {
	const sheetIndex = 13
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

	slice, err := parseOtherIntlPrices(params, sheetIndex)
	suite.NoError(err, "parseOtherIntlPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "13_3d_other_intl_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
