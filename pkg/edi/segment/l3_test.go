package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL3() {
	validWeightL3 := L3{
		Weight:          300.0,
		WeightQualifier: "B",
		PriceCents:      100,
	}

	suite.T().Run("validate success weight", func(t *testing.T) {
		err := suite.validator.Struct(validWeightL3)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		l3 := L3{
			Weight:     300.0,
			PriceCents: 100,
		}

		err := suite.validator.Struct(l3)
		suite.ValidateError(err, "WeightQualifier", "required_with")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		l3 := L3{
			Weight:          300.0,
			WeightQualifier: "INVALID",
			PriceCents:      100,
		}

		err := suite.validator.Struct(l3)
		suite.ValidateError(err, "WeightQualifier", "eq")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		l3 := L3{
			Weight:          99999999999, // 10 digits
			WeightQualifier: "B",
			PriceCents:      9999999999999, // 13 digits
		}

		err := suite.validator.Struct(l3)
		suite.ValidateError(err, "Weight", "max")
		suite.ValidateError(err, "PriceCents", "max")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		l3 := L3{
			Weight:     -1,
			PriceCents: -9999999999999, // 13 digits
		}

		err := suite.validator.Struct(l3)
		suite.ValidateError(err, "Weight", "min")
		suite.ValidateError(err, "WeightQualifier", "required_with")
		suite.ValidateError(err, "PriceCents", "min")
		suite.ValidateErrorLen(err, 3)
	})
}
