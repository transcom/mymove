package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateAK1() {
	validAK1 := AK1{
		FunctionalIdentifierCode: "SI",
		GroupControlNumber:       1234567,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validAK1)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		ak1 := AK1{
			FunctionalIdentifierCode: "XX", // eq
			GroupControlNumber:       0,    // min
		}

		err := suite.validator.Struct(ak1)
		suite.ValidateError(err, "FunctionalIdentifierCode", "eq")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		ak1 := AK1{
			FunctionalIdentifierCode: "SI",
			GroupControlNumber:       999999999999999, // max
		}

		err := suite.validator.Struct(ak1)
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
