package pricing

import (
	"strconv"
	"time"

	"github.com/go-openapi/swag"
)

func (suite *PricingParserSuite) Test_parseShipmentManagementServicesPrices() {
	const sheetIndex = 16
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

	suite.Run("parse sheet and check csv", func() {
		slice, err := parseShipmentManagementServicesPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseShipmentManagementServices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("management"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "16_4a_mgmt_coun_trans_prices_management_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("try parse wrong sheet index", func() {
		_, err := parseShipmentManagementServicesPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseShipmentManagementServicesPrices function failed") {
			suite.Equal("parseShipmentManagementServices expected to process sheet 16, but received sheetIndex 15", err.Error())
		}
	})
}

func (suite *PricingParserSuite) Test_parseCounselServicesPrices() {
	const sheetIndex = 16
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

	suite.Run("parse sheet and check csv", func() {
		slice, err := parseCounselingServicesPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseCounselingServicesPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("counsel"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "16_4a_mgmt_coun_trans_prices_counsel_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("try parse wrong sheet index", func() {
		_, err := parseCounselingServicesPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseCounselingServicesPrices function failed") {
			suite.Equal("parseCounselingServicesPrices expected to process sheet 16, but received sheetIndex 15", err.Error())
		}
	})
}

func (suite *PricingParserSuite) Test_parseTransitionPrices() {
	const sheetIndex = 16
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

	suite.Run("parse sheet and check csv", func() {
		slice, err := parseTransitionPrices(suite.AppContextForTest(), params, sheetIndex)
		suite.NoError(err, "parseTransitionPrices function failed")

		outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("transition"))
		err = createCSV(suite.AppContextForTest(), outputFilename, slice)
		suite.NoError(err, "could not create CSV")

		const goldenFilename string = "16_4a_mgmt_coun_trans_prices_transition_golden.csv"
		suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
	})

	suite.Run("try parse wrong sheet index", func() {
		_, err := parseTransitionPrices(suite.AppContextForTest(), params, sheetIndex-1)
		if suite.Error(err, "parseTransitionPrices function failed") {
			suite.Equal("parseTransitionPrices expected to process sheet 16, but received sheetIndex 15", err.Error())
		}
	})
}

func (suite *PricingParserSuite) Test_verifyManagementCounselTransitionPrices() {
	const sheetIndex = 16
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

	suite.Run("verify good sheet", func() {
		err := verifyManagementCounselTransitionPrices(params, sheetIndex)
		suite.NoError(err, "verifyManagementCounselTransitionPrices function failed")
	})

	suite.Run("verify wrong sheet", func() {
		err := verifyManagementCounselTransitionPrices(params, sheetIndex-2)
		if suite.Error(err, "verifyManagementCounselTransitionPrices function failed") {
			suite.Equal("verifyManagementCounselTransitionPrices expected to process sheet 16, but received sheetIndex 14", err.Error())
		}
	})
}
