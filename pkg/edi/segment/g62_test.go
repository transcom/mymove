package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateG62() {
	validG62ActualPickupDateTime := G62{
		DateQualifier: 86,
		Date:          "20200909",
		TimeQualifier: 8,
		Time:          "1617",
	}
	validG62RequestedPickupDateTime := G62{
		DateQualifier: 68,
		Date:          "20200909",
		TimeQualifier: 5,
		Time:          "1617",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validG62ActualPickupDateTime)
		suite.NoError(err)
		err = suite.validator.Struct(validG62RequestedPickupDateTime)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		g62 := G62{
			DateQualifier: 42,         // oneof
			Date:          "20190945", // timeformat
			TimeQualifier: 42,         // oneof
			Time:          "2517",     // timeformat
		}

		err := suite.validator.Struct(g62)
		suite.ValidateError(err, "DateQualifier", "oneof")
		suite.ValidateError(err, "Date", "timeformat")
		suite.ValidateError(err, "TimeQualifier", "oneof")
		suite.ValidateError(err, "Time", "timeformat")
		suite.ValidateErrorLen(err, 4)
	})
}
