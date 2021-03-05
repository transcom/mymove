package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateAK3() {
	freightRate := 0
	validAK3 := AK3{
		LadingLineItemNumber: 1,
		FreightRate:          &freightRate,
		RateValueQualifier:   "LB",
		Charge:               10000,
	}

	altValidAK3 := AK3{
		LadingLineItemNumber: 12,
		Charge:               10000,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validAK3)
		suite.NoError(err)
		suite.Equal([]string{"AK3", "1", "0", "LB", "10000"}, validAK3.StringArray())
	})

	suite.T().Run("validate alt success", func(t *testing.T) {
		err := suite.validator.Struct(altValidAK3)
		suite.NoError(err)
		suite.Equal([]string{"AK3", "12", "", "", "10000"}, altValidAK3.StringArray())
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		freightRate := -1
		l1 := AK3{
			// LadingLineItemNumber:          // required
			FreightRate:        &freightRate, // min
			RateValueQualifier: "XX",         // eq
			// Charge:                        // required
		}

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "required")
		suite.ValidateError(err, "FreightRate", "min")
		suite.ValidateError(err, "RateValueQualifier", "eq")
		suite.ValidateError(err, "Charge", "required")
		suite.ValidateErrorLen(err, 4)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		l1 := validAK3
		l1.LadingLineItemNumber = 1000 // max
		l1.Charge = 1000000000000      // max

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "max")
		suite.ValidateError(err, "Charge", "max")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		l1 := validAK3
		l1.LadingLineItemNumber = -3 // min
		l1.RateValueQualifier = ""   // required
		l1.Charge = -1000000000000   // min

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "min")
		suite.ValidateError(err, "RateValueQualifier", "required_with")
		suite.ValidateError(err, "Charge", "min")
		suite.ValidateErrorLen(err, 3)
	})
}
