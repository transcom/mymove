package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL0() {
	validBilledL0 := L0{
		LadingLineItemNumber:   1,
		BilledRatedAsQuantity:  3.0,
		BilledRatedAsQualifier: "XX",
	}

	validWeightL0 := L0{
		LadingLineItemNumber: 1,
		Weight:               300.0,
		WeightQualifier:      "B",
		WeightUnitCode:       "L",
	}

	suite.T().Run("validate success billed", func(t *testing.T) {
		err := suite.validator.Struct(validBilledL0)
		suite.NoError(err)
	})

	suite.T().Run("validate success weight", func(t *testing.T) {
		err := suite.validator.Struct(validWeightL0)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		l0 := L0{
			LadingLineItemNumber:  2000, // max
			BilledRatedAsQuantity: 3.0,  // required_with
		}

		err := suite.validator.Struct(l0)
		suite.ValidateError(err, "LadingLineItemNumber", "max")
		suite.ValidateError(err, "BilledRatedAsQualifier", "required_with")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		l0 := L0{
			LadingLineItemNumber: 0,     // min
			Weight:               300.0, // required_with
		}

		err := suite.validator.Struct(l0)
		suite.ValidateError(err, "LadingLineItemNumber", "min")
		suite.ValidateError(err, "WeightQualifier", "required_with")
		suite.ValidateError(err, "WeightUnitCode", "required_with")
		suite.ValidateErrorLen(err, 3)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		l0 := validBilledL0
		l0.BilledRatedAsQualifier = "ABC" // len

		err := suite.validator.Struct(l0)
		suite.ValidateError(err, "BilledRatedAsQualifier", "len")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 4", func(t *testing.T) {
		l0 := validWeightL0
		l0.WeightQualifier = "X" // eq
		l0.WeightUnitCode = "X"  // eq

		err := suite.validator.Struct(l0)
		suite.ValidateError(err, "WeightQualifier", "eq")
		suite.ValidateError(err, "WeightUnitCode", "eq")
		suite.ValidateErrorLen(err, 2)
	})
}
