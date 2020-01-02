package pricing

import (
	"strconv"
	"time"
)

// Test_parseDomesticOtherPricesPack
func (suite *PricingParserSuite) Test_parseDomesticOtherPricesPack() {
	const sheetIndex = 8
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

	slice, err := parseDomesticOtherPricesPack(params, sheetIndex)
	suite.NoError(err, "parseDomesticOtherPricesPack function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "8_2c_domestic_other_prices_pack_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}

// Test_parseDomesticOtherPricesSit
func (suite *PricingParserSuite) Test_parseDomesticOtherPricesSit() {
	const sheetIndex = 8
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

	slice, err := parseDomesticOtherPricesSit(params, sheetIndex)
	suite.NoError(err, "parseDomesticOtherPricesSit function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "8_2c_domestic_other_prices_sit_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
}
