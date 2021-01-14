package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL1() {
	freightRate := 0
	validL1 := L1{
		LadingLineItemNumber: 1,
		FreightRate:          &freightRate,
		RateValueQualifier:   "LB",
		Charge:               100.00,
	}

	altValidL1 := L1{
		LadingLineItemNumber: 12,
		Charge:               100.00,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validL1)
		suite.NoError(err)
		suite.Equal([]string{"L1", "1", "0", "LB", "10000"}, validL1.StringArray())
	})

	suite.T().Run("validate alt success", func(t *testing.T) {
		err := suite.validator.Struct(altValidL1)
		suite.NoError(err)
		suite.Equal([]string{"L1", "12", "", "", "10000"}, altValidL1.StringArray())
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		freightRate := -1
		l1 := L1{
			LadingLineItemNumber: -3,           // min
			FreightRate:          &freightRate, // min
			RateValueQualifier:   "XX",         // eq
			Charge:               0,            // required
		}

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "min")
		suite.ValidateError(err, "FreightRate", "min")
		suite.ValidateError(err, "RateValueQualifier", "eq")
		suite.ValidateError(err, "Charge", "required")
		suite.ValidateErrorLen(err, 4)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		l1 := validL1
		l1.LadingLineItemNumber = 1000 // max

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		l1 := validL1
		l1.RateValueQualifier = "" // required

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "RateValueQualifier", "required_with")
		suite.ValidateErrorLen(err, 1)
	})
}
