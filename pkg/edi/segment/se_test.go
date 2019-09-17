package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateSE() {
	validSE := SE{
		NumberOfIncludedSegments:    12345,
		TransactionSetControlNumber: "ABCDE",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validSE)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		se := SE{
			NumberOfIncludedSegments:    0,     // min
			TransactionSetControlNumber: "ABC", // min
		}

		err := suite.validator.Struct(se)
		suite.ValidateError(err, "NumberOfIncludedSegments", "min")
		suite.ValidateError(err, "TransactionSetControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		se := SE{
			NumberOfIncludedSegments:    10000000000,  // max
			TransactionSetControlNumber: "1234567890", // max
		}

		err := suite.validator.Struct(se)
		suite.ValidateError(err, "NumberOfIncludedSegments", "max")
		suite.ValidateError(err, "TransactionSetControlNumber", "max")
		suite.ValidateErrorLen(err, 2)
	})
}
