package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL10() {
	validL10 := L10{
		Weight:          100.0,
		WeightQualifier: "B",
		WeightUnitCode:  "L",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validL10)
		suite.NoError(err)
	})

	suite.T().Run("validate failure", func(t *testing.T) {
		l10 := L10{
			// Weight required
			WeightQualifier: "X",
			WeightUnitCode:  "X",
		}

		err := suite.validator.Struct(l10)
		suite.ValidateError(err, "Weight", "required")
		suite.ValidateError(err, "WeightQualifier", "eq")
		suite.ValidateError(err, "WeightUnitCode", "eq")
		suite.ValidateErrorLen(err, 3)
	})
}
