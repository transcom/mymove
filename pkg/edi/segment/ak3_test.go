package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateAK3() {
	validAK3 := AK3{
		SegmentIDCode:                   "ID",
		SegmentPositionInTransactionSet: 12345,
		LoopIdentifierCode:              "CODE",
		SegmentSyntaxErrorCode:          "ERR",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validAK3)
		suite.NoError(err)
		suite.Equal([]string{"AK3", "ID", "12345", "CODE", "ERR"}, validAK3.StringArray())
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		ak3 := AK3{
			SegmentIDCode:                   "", // min
			SegmentPositionInTransactionSet: 0,  // min
			LoopIdentifierCode:              "", // min
			SegmentSyntaxErrorCode:          "", // min
		}

		err := suite.validator.Struct(ak3)
		suite.ValidateError(err, "SegmentIDCode", "min")
		suite.ValidateError(err, "SegmentPositionInTransactionSet", "min")
		suite.ValidateError(err, "LoopIdentifierCode", "min")
		suite.ValidateError(err, "SegmentSyntaxErrorCode", "min")
		suite.ValidateErrorLen(err, 4)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		ak3 := AK3{
			SegmentIDCode:                   "XXXX",    // max
			SegmentPositionInTransactionSet: 9999999,   // max
			LoopIdentifierCode:              "XXXXXXX", // max
			SegmentSyntaxErrorCode:          "XXXX",    // max
		}

		err := suite.validator.Struct(ak3)
		suite.ValidateError(err, "SegmentIDCode", "max")
		suite.ValidateError(err, "SegmentPositionInTransactionSet", "max")
		suite.ValidateError(err, "LoopIdentifierCode", "max")
		suite.ValidateError(err, "SegmentSyntaxErrorCode", "max")
		suite.ValidateErrorLen(err, 4)
	})
}
