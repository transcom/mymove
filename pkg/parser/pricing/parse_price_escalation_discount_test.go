package pricing

import (
	"strconv"
	"time"
)

func (suite *PricingParserSuite) Test_verifyPriceEscalationDiscount() {
	const sheetIndex = 18
	InitDataSheetInfo()

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

	suite.Run("normal operation", func() {
		err := verifyPriceEscalationDiscount(params, sheetIndex)
		suite.NoError(err, "verifyPriceEscalationDiscount function failed")
	})

	suite.Run("passing in a bad sheet index", func() {
		err := verifyPriceEscalationDiscount(params, 7)
		if suite.Error(err) {
			suite.Equal("verifyPriceEscalationDiscount expected to process sheet 18, but received sheetIndex 7", err.Error())
		}
	})
}

func (suite *PricingParserSuite) Test_parsePriceEscalationDiscount() {
	const sheetIndex = 18
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

	suite.Run("normal operation", func() {
		slice, err := parsePriceEscalationDiscount(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parsePriceEscalationDiscount function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const domesticGoldenFilename string = "18_5b_price_escalation_discount_golden.csv"
		suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
	})

	suite.Run("passing in a bad sheet index", func() {
		_, err := parsePriceEscalationDiscount(suite.AppContextForTest(), params, 15)
		if suite.Error(err) {
			suite.Equal("parsePriceEscalationDiscount expected to process sheet 18, but received sheetIndex 15", err.Error())
		}
	})
}
