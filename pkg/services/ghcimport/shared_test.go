package ghcimport

import (
	"testing"
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
		suite.T().Run(test.name, func(t *testing.T) {
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
		suite.T().Run(test.name, func(t *testing.T) {
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
		suite.T().Run(test.name, func(t *testing.T) {
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
		suite.T().Run(test.name, func(t *testing.T) {
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

func (suite *GHCRateEngineImportSuite) Test_priceStringToFloat() {
	tests := []struct {
		name        string
		input       string
		expected    float64
		shouldError bool
	}{
		{"price with dollar sign", "$3.557", 3.557, false},
		{"price with no dollar sign", "3.557", 3.557, false},
		{"price as integer", "3", 3, false},
		{"not a number", "3.53X", 0, true},
	}
	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			result, err := priceStringToFloat(test.input)
			suite.Equal(test.expected, result)
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
		{"price with no dollar sign", "3.557", 355700, false},
		{"price as integer", "3", 300000, false},
		{"not a number", "3.53X", 0, true},
	}
	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
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
