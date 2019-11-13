package pricing

import (
	"strconv"
	"testing"
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

	suite.T().Run("normal operation", func(t *testing.T) {
		err := verifyPriceEscalationDiscount(params, sheetIndex)
		suite.NoError(err, "verifyPriceEscalationDiscount function failed")
	})

	suite.T().Run("passing in a bad sheet index", func(t *testing.T) {
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

	suite.T().Run("normal operation", func(t *testing.T) {
		slice, err := parsePriceEscalationDiscount(params, sheetIndex)
		suite.NoError(err, "parsePriceEscalationDiscount function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, nil)
		err = createCSV(outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const domesticGoldenFilename string = "18_5b_price_escalation_discount_golden.csv"
		suite.helperTestExpectedFileOutput(domesticGoldenFilename, outputFilename)
	})

	suite.T().Run("passing in a bad sheet index", func(t *testing.T) {
		_, err := parsePriceEscalationDiscount(params, 15)
		if suite.Error(err) {
			suite.Equal("parsePriceEscalationDiscount expected to process sheet 18, but received sheetIndex 15", err.Error())
		}
	})
}
