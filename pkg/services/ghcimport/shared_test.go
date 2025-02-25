package ghcimport

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_stringToInteger() {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldError bool
	}{
		{"with decimal point", "25.0", 25, false},
		{"no decimal point", "25", 25, false},
		{"not a number", "2A", 0, true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := stringToInteger(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_cleanServiceAreaNumber() {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{"with decimal point, needs leading zeros", "4.0", "004", false},
		{"no decimal point", "225", "225", false},
		{"not a number", "B3", "", true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := cleanServiceAreaNumber(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_cleanZip3() {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{"with decimal point, needs leading zeros", "15.0", "015", false},
		{"no decimal point", "309", "309", false},
		{"not a number", "30L", "", true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := cleanZip3(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_isPeakPeriod() {
	tests := []struct {
		name        string
		input       string
		expected    bool
		shouldError bool
	}{
		{"peak", "Peak", true, false},
		{"peak, upper case", "PEAK", true, false},
		{"non-peak", "NonPeak", false, false},
		{"non-peak, upper case", "NONPEAK", false, false},
		{"invalid period", "Non-Peak", false, true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := isPeakPeriod(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_getPriceParts() {
	tests := []struct {
		name                   string
		rawPrice               string
		decimalPlaces          int
		expectedIntegerPart    int
		expectedFractionalPart int
		shouldError            bool
	}{
		{"at expected decimal places", "$3.557", 3, 3, 557, false},
		{"less than max decimal places", "$3.5", 3, 0, 0, true},
		{"more than max decimal places", "$3.5777", 3, 0, 0, true},
		{"no dollar sign", "2.001", 3, 2, 1, false},
		{"very small, no dollar sign", "0.005", 3, 0, 5, false},
		{"no decimal point", "1", 3, 0, 0, true},
		{"invalid price", "$3.ABC", 3, 0, 0, true},
		{"empty string", "", 4, 0, 0, true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			integerPart, fractionalPart, err := getPriceParts(test.rawPrice, test.decimalPlaces)
			suite.Equal(test.expectedIntegerPart, integerPart)
			suite.Equal(test.expectedFractionalPart, fractionalPart)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_priceToMillicents() {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldError bool
	}{
		{"price with dollar sign", "$3.557", 355700, false},
		{"leading zeros in decimal part", "$3.005", 300500, false},
		{"two decimal places", "$3.55", 0, true},
		{"more than expected decimal places", "$3.5571", 0, true},
		{"input of zero", "0", 0, true},
		{"empty string", "", 0, true},
		{"price without dollar sign", "0.001", 100, false},
		{"price as integer", "3", 0, true},
		{"not a number", "3.53X", 0, true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := priceToMillicents(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_priceToCents() {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldError bool
	}{
		{"price with dollar sign", "$3.55", 355, false},
		{"leading zeros in decimal part", "$3.01", 301, false},
		{"more than expected decimal places", "$3.551", 0, true},
		{"input of zero", "0", 0, true},
		{"empty string", "", 0, true},
		{"price without dollar sign", "0.01", 1, false},
		{"price as integer", "3", 0, true},
		{"not a number", "3.5X", 0, true},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := priceToCents(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GHCRateEngineImportSuite) Test_getMarket() {
	tests := []struct {
		name        string
		input       string
		expected    models.Market
		shouldError bool
	}{
		{"CONUS", "CONUS", "C", false},
		{"OCONUS", "OCONUS", "O", false},
		{"XONUS", "XONUS", "invalid market", true},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			result, err := getMarket(test.input)
			suite.Equal(test.expected, result)
			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
