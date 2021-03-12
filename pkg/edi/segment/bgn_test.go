package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateBGN() {
	suite.T().Run("validate success all fields", func(t *testing.T) {
		validBGN := BGN{
			TransactionSetPurposeCode: "11",
			ReferenceIdentification:   "hello",
			Date:                      "20210310",
		}
		err := suite.validator.Struct(validBGN)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		bgn := BGN{
			TransactionSetPurposeCode: "10",       // eq
			ReferenceIdentification:   "",         // min
			Date:                      "20211313", // datetime
		}

		err := suite.validator.Struct(bgn)
		suite.ValidateError(err, "TransactionSetPurposeCode", "eq")
		suite.ValidateError(err, "ReferenceIdentification", "min")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateErrorLen(err, 3)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		bgn := BGN{
			TransactionSetPurposeCode: "11",
			ReferenceIdentification:   "long string that exceeds max length", // max
			Date:                      "20210310",
		}

		err := suite.validator.Struct(bgn)
		suite.ValidateError(err, "ReferenceIdentification", "max")
		suite.ValidateErrorLen(err, 1)
	})
}

func (suite *SegmentSuite) TestStringArrayBGN() {
	suite.T().Run("string array all fields", func(t *testing.T) {
		validBGN := BGN{
			TransactionSetPurposeCode: "11",
			ReferenceIdentification:   "hello",
			Date:                      "20210310",
		}
		arrayValidBGN := []string{"BGN", "11", "hello", "20210310"}
		suite.Equal(arrayValidBGN, validBGN.StringArray())
	})
}

func (suite *SegmentSuite) TestParseBGN() {
	suite.T().Run("parse success all fields", func(t *testing.T) {
		arrayValidBGN := []string{"11", "hello", "20210310"}
		expectedBGN := BGN{
			TransactionSetPurposeCode: "11",
			ReferenceIdentification:   "hello",
			Date:                      "20210310",
		}

		var validBGN BGN
		err := validBGN.Parse(arrayValidBGN)
		if suite.NoError(err) {
			suite.Equal(expectedBGN, validBGN)
		}
	})

	suite.T().Run("wrong number of fields", func(t *testing.T) {
		badArrayBGN := []string{"11", "hello"}
		var badBGN BGN
		err := badBGN.Parse(badArrayBGN)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})
}
