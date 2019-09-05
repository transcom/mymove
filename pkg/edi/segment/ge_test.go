package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateGE() {
	validGE := GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              1234567,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validGE)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		ge := GE{
			NumberOfTransactionSetsIncluded: 2, // eq
			GroupControlNumber:              0, // min
		}

		err := suite.validator.Struct(ge)
		suite.ValidateError(err, "NumberOfTransactionSetsIncluded", "eq")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		ge := validGE
		ge.GroupControlNumber = 1000000000 // max

		err := suite.validator.Struct(ge)
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
