package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL1() {
	validL1 := L1{
		FreightRate:              0,
		RateValueQualifier:       "LB",
		Charge:                   100.00,
		SpecialChargeDescription: "ABC",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validL1)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		l1 := L1{
			LadingLineItemNumber:     -3,   // min
			FreightRate:              -1,   // min
			RateValueQualifier:       "XX", // eq
			Charge:                   0,    // required
			SpecialChargeDescription: "X",  // min
		}

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "min")
		suite.ValidateError(err, "FreightRate", "min")
		suite.ValidateError(err, "RateValueQualifier", "eq")
		suite.ValidateError(err, "Charge", "required")
		suite.ValidateError(err, "SpecialChargeDescription", "min")
		suite.ValidateErrorLen(err, 5)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		l1 := validL1
		l1.LadingLineItemNumber = 1000                             // max
		l1.SpecialChargeDescription = "12345678901234567890123456" // max

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "max")
		suite.ValidateError(err, "SpecialChargeDescription", "max")
		suite.ValidateErrorLen(err, 2)
	})
}
