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

	slice, err := parseShipmentManagementServicesPrices(params, sheetIndex)
	suite.NoError(err, "parseShipmentManagementServices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("management"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "16_4a_mgmt_coun_trans_prices_management_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
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

	slice, err := parseCounselingServicesPrices(params, sheetIndex)
	suite.NoError(err, "parseCounselingServicesPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("counsel"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "16_4a_mgmt_coun_trans_prices_counsel_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
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

	slice, err := parseTransitionPrices(params, sheetIndex)
	suite.NoError(err, "parseTransitionPrices function failed")

	outputFilename := dataSheet.generateOutputFilename(sheetIndex, params.RunTime, swag.String("transition"))
	err = createCSV(outputFilename, slice)
	suite.NoError(err, "could not create CSV")

	const goldenFilename string = "16_4a_mgmt_coun_trans_prices_transition_golden.csv"
	suite.helperTestExpectedFileOutput(goldenFilename, outputFilename)
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

	err := verifyManagementCounselTransitionPrices(params, sheetIndex)
	suite.NoError(err, "verifyManagementCounselTransitionPrices function failed")
}

func (suite *PricingParserSuite) Test_verifyManagementCounselTransitionPricesWithWrongSheet() {
	const sheetIndex = 15
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

	err := verifyManagementCounselTransitionPrices(params, sheetIndex)
	suite.Error(err, "verifyManagementCounselTransitionPrices function failed")
	suite.Equal("verifyManagementCounselTransitionPrices expected to process sheet 16, but received sheetIndex 15", err.Error())
}
