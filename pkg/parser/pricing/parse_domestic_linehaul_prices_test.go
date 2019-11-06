package pricing

import (
	"time"
)

// Test_parseDomesticLinehaulPrices
func (suite *ParseRateEngineGHCXLSXSuite) Test_parseDomesticLinehaulPrices() {
	xlsxDataSheets := InitDataSheetInfo()
	params := ParamConfig{
		ProcessAll:   false,
		ShowOutput:   false,
		XlsxFilename: &suite.xlsxFilename,
		XlsxSheets:   []string{"6"},
		SaveToFile:   true,
		RunTime:      time.Now(),
		XlsxFile:     suite.xlsxFile,
		RunVerify:    true,
	}

	const sheetIndex int = 6
	dataSheet := xlsxDataSheets[sheetIndex]

	slice, err := parseDomesticLinehaulPrices(params, sheetIndex)
	suite.NoError(err, "parseDomesticLinehaulPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, sheetIndex, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "6_2a_domestic_linehaul_prices_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
