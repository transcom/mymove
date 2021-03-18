package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateAK5() {
	suite.T().Run("validate success all fields", func(t *testing.T) {
		validAK5 := AK5{
			TransactionSetAcknowledgmentCode:   "A",
			TransactionSetSyntaxErrorCodeAK502: "abc",
			TransactionSetSyntaxErrorCodeAK503: "def",
			TransactionSetSyntaxErrorCodeAK504: "ghi",
			TransactionSetSyntaxErrorCodeAK505: "jkl",
			TransactionSetSyntaxErrorCodeAK506: "mno",
		}
		err := suite.validator.Struct(validAK5)
		suite.NoError(err)
	})

	suite.T().Run("validate success only required fields", func(t *testing.T) {
		validAK5 := AK5{
			TransactionSetAcknowledgmentCode: "A",
		}
		err := suite.validator.Struct(validAK5)
		suite.NoError(err)
	})

	suite.T().Run("failure due to missing required fields", func(t *testing.T) {

		ak5 := AK5{
			TransactionSetSyntaxErrorCodeAK502: "abc",
			TransactionSetSyntaxErrorCodeAK503: "def",
			TransactionSetSyntaxErrorCodeAK504: "ghi",
			TransactionSetSyntaxErrorCodeAK505: "jkl",
			TransactionSetSyntaxErrorCodeAK506: "mno",
		}
		err := suite.validator.Struct(ak5)
		suite.ValidateError(err, "TransactionSetAcknowledgmentCode", "len")
	})

	suite.T().Run("validate failure max", func(t *testing.T) {
		// length of characters are more than max
		ak5 := AK5{
			TransactionSetAcknowledgmentCode:   "AAAAA",
			TransactionSetSyntaxErrorCodeAK502: "abcz",
			TransactionSetSyntaxErrorCodeAK503: "defz",
			TransactionSetSyntaxErrorCodeAK504: "ghiz",
			TransactionSetSyntaxErrorCodeAK505: "jklz",
			TransactionSetSyntaxErrorCodeAK506: "mnoz",
		}

		err := suite.validator.Struct(ak5)
		suite.ValidateError(err, "TransactionSetAcknowledgmentCode", "len")
		suite.ValidateError(err, "TransactionSetSyntaxErrorCodeAK502", "max")
		suite.ValidateError(err, "TransactionSetSyntaxErrorCodeAK503", "max")
		suite.ValidateError(err, "TransactionSetSyntaxErrorCodeAK504", "max")
		suite.ValidateError(err, "TransactionSetSyntaxErrorCodeAK505", "max")
		suite.ValidateError(err, "TransactionSetSyntaxErrorCodeAK506", "max")
		suite.ValidateErrorLen(err, 6)
	})

	suite.T().Run("validate failure min", func(t *testing.T) {
		// length of characters are less than min
		ak5 := AK5{
			TransactionSetAcknowledgmentCode:   "",
			TransactionSetSyntaxErrorCodeAK502: "",
			TransactionSetSyntaxErrorCodeAK503: "",
			TransactionSetSyntaxErrorCodeAK504: "",
			TransactionSetSyntaxErrorCodeAK505: "",
			TransactionSetSyntaxErrorCodeAK506: "",
		}

		err := suite.validator.Struct(ak5)
		suite.ValidateError(err, "TransactionSetAcknowledgmentCode", "len")
		suite.ValidateErrorLen(err, 1)
	})
}

func (suite *SegmentSuite) TestStringArrayAK5() {
	suite.T().Run("string array all fields", func(t *testing.T) {
		validAK5 := AK5{
			TransactionSetAcknowledgmentCode:   "A",
			TransactionSetSyntaxErrorCodeAK502: "abc",
			TransactionSetSyntaxErrorCodeAK503: "def",
			TransactionSetSyntaxErrorCodeAK504: "ghi",
			TransactionSetSyntaxErrorCodeAK505: "jkl",
			TransactionSetSyntaxErrorCodeAK506: "mno",
		}
		arrayValidAK5 := []string{"AK5", "A", "abc", "def", "ghi", "jkl", "mno"}
		suite.Equal(arrayValidAK5, validAK5.StringArray())
	})

	suite.T().Run("string array only required fields", func(t *testing.T) {
		validOptionalAK5 := AK5{
			TransactionSetAcknowledgmentCode: "A",
		}
		arrayValidOptionalAK5 := []string{"AK5", "A", "", "", "", "", ""}
		suite.Equal(arrayValidOptionalAK5, validOptionalAK5.StringArray())
	})
}

func (suite *SegmentSuite) TestParseAK5() {
	suite.T().Run("parse success all fields", func(t *testing.T) {
		arrayValidAK5 := []string{"A", "abc", "def", "ghi", "jkl", "mno"}
		expectedAK5 := AK5{
			TransactionSetAcknowledgmentCode:   "A",
			TransactionSetSyntaxErrorCodeAK502: "abc",
			TransactionSetSyntaxErrorCodeAK503: "def",
			TransactionSetSyntaxErrorCodeAK504: "ghi",
			TransactionSetSyntaxErrorCodeAK505: "jkl",
			TransactionSetSyntaxErrorCodeAK506: "mno",
		}

		var validAK5 AK5
		err := validAK5.Parse(arrayValidAK5)
		if suite.NoError(err) {
			suite.Equal(expectedAK5, validAK5)
		}
	})

	suite.T().Run("parse success only required fields", func(t *testing.T) {
		arrayValidOptionalAK5 := []string{"A", "", "", "", "", ""}
		expectedOptionalAK5 := AK5{
			TransactionSetAcknowledgmentCode: "A",
		}

		var validOptionalAK5 AK5
		err := validOptionalAK5.Parse(arrayValidOptionalAK5)
		if suite.NoError(err) {
			suite.Equal(expectedOptionalAK5, validOptionalAK5)
		}
	})

	suite.T().Run("wrong number of fields", func(t *testing.T) {
		badArrayAK5 := []string{"A", "abc", "def", "ghi", "jkl", "mno", "zzz"}
		var badAK5 AK5
		err := badAK5.Parse(badArrayAK5)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})
}
