package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateLX() {
	validLX := LX{
		AssignedNumber: 1,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validLX)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		lx := LX{
			AssignedNumber: 0, // min
		}

		err := suite.validator.Struct(lx)
		suite.ValidateError(err, "AssignedNumber", "min")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		lx := LX{
			AssignedNumber: 1000000, // max
		}

		err := suite.validator.Struct(lx)
		suite.ValidateError(err, "AssignedNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
