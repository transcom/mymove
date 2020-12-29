package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL1() {
	validL1 := L1{
		FreightRate:        0,
		RateValueQualifier: "LB",
		Charge:             100.00,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validL1)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		l1 := L1{
			LadingLineItemNumber: -3,   // min
			FreightRate:          -1,   // min
			RateValueQualifier:   "XX", // eq
			Charge:               0,    // required
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
}
