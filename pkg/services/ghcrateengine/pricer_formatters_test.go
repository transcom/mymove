package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestFormatCents() {
	cents := unit.Cents(1)
	result := FormatCents(cents)
	expected := "0.01"
	suite.Equal(expected, result)

	cents = unit.Cents(100)
	result = FormatCents(cents)
	expected = "1.00"
	suite.Equal(expected, result)

	cents = unit.Cents(10099)
	result = FormatCents(cents)
	expected = "100.99"
	suite.Equal(expected, result)
}

func (suite *GHCRateEngineServiceSuite) TestFormatFloat() {
	num := 1.00020000
	result := FormatFloat(num)
	expected := "1.0002"
	suite.Equal(expected, result)

	num = 1.002
	result = FormatFloat(num)
	expected = "1.002"
	suite.Equal(expected, result)
}
