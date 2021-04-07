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

// TODO: Add more tests
