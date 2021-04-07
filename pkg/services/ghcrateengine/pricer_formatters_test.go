package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestFormatTimestamp() {
	testTime := time.Date(2021, time.June, 5, 7, 33, 11, 456, time.UTC)
	suite.Equal("2021-06-05T07:33:11Z", FormatTimestamp(testTime))
}

func (suite *GHCRateEngineServiceSuite) TestFormatDate() {
	testDate := time.Date(2021, time.June, 5, 7, 33, 11, 456, time.UTC)
	suite.Equal("2021-06-05", FormatDate(testDate))
}

func (suite *GHCRateEngineServiceSuite) TestFormatCents() {
	testCases := []struct {
		inputCents unit.Cents
		expected   string
	}{
		{unit.Cents(1), "0.01"},
		{unit.Cents(100), "1.00"},
		{unit.Cents(10099), "100.99"},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, FormatCents(tc.inputCents))
	}
}

func (suite *GHCRateEngineServiceSuite) TestFormatBool() {
	testCases := []struct {
		inputBool bool
		expected  string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, FormatBool(tc.inputBool))
	}
}

func (suite *GHCRateEngineServiceSuite) TestFormatFloat() {
	testCases := []struct {
		inputFloat float64
		precision  int
		expected   string
	}{
		{1.234, 2, "1.23"},
		{1.234567, 3, "1.235"},
		{1.23456789, -1, "1.23456789"},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, FormatFloat(tc.inputFloat, tc.precision))
	}
}
