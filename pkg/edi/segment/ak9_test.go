package edisegment

import (
	"fmt"
	"testing"
)

func (suite *SegmentSuite) TestValidateAK9() {
	suite.T().Run("validate success all fields", func(t *testing.T) {
		validAK9 := AK9{
			FunctionalGroupAcknowledgeCode:      "A",
			NumberOfTransactionSetsIncluded:     1,
			NumberOfReceivedTransactionSets:     2,
			NumberOfAcceptedTransactionSets:     3,
			FunctionalGroupSyntaxErrorCodeAK905: "AAA",
			FunctionalGroupSyntaxErrorCodeAK906: "BBB",
			FunctionalGroupSyntaxErrorCodeAK907: "CCC",
			FunctionalGroupSyntaxErrorCodeAK908: "DDD",
			FunctionalGroupSyntaxErrorCodeAK909: "EEE",
		}
		err := suite.validator.Struct(validAK9)
		suite.NoError(err)
	})

	suite.T().Run("validate success only required fields", func(t *testing.T) {
		validAK9 := AK9{
			FunctionalGroupAcknowledgeCode:  "E",
			NumberOfTransactionSetsIncluded: 1,
			NumberOfReceivedTransactionSets: 2,
			NumberOfAcceptedTransactionSets: 3,
		}
		err := suite.validator.Struct(validAK9)
		suite.NoError(err)
	})

	suite.T().Run("validate success for all valid FunctionalGroupAcknowledgeCode values", func(t *testing.T) {
		validAK9 := AK9{
			FunctionalGroupAcknowledgeCode:  "A",
			NumberOfTransactionSetsIncluded: 1,
			NumberOfReceivedTransactionSets: 2,
			NumberOfAcceptedTransactionSets: 3,
		}
		allowedValues := []string{"A", "E", "P", "R"}
		for _, val := range allowedValues {
			validAK9.FunctionalGroupAcknowledgeCode = val
			err := suite.validator.Struct(validAK9)
			suite.NoError(err, fmt.Sprintf("Failed to validate allowed value: \"%s\"", val))
		}
	})

	suite.T().Run("validate failure for invalid FunctionalGroupAcknowledgeCode", func(t *testing.T) {
		validAK9 := AK9{
			FunctionalGroupAcknowledgeCode:  "B",
			NumberOfTransactionSetsIncluded: 1,
			NumberOfReceivedTransactionSets: 2,
			NumberOfAcceptedTransactionSets: 3,
		}
		err := suite.validator.Struct(validAK9)
		suite.ValidateError(err, "FunctionalGroupAcknowledgeCode", "oneof")
	})

	suite.T().Run("failure due to missing required fields", func(t *testing.T) {

		ak9 := AK9{
			FunctionalGroupSyntaxErrorCodeAK905: "AAA",
			FunctionalGroupSyntaxErrorCodeAK906: "BBB",
			FunctionalGroupSyntaxErrorCodeAK907: "CCC",
			FunctionalGroupSyntaxErrorCodeAK908: "DDD",
			FunctionalGroupSyntaxErrorCodeAK909: "EEE",
		}
		err := suite.validator.Struct(ak9)
		suite.ValidateError(err, "FunctionalGroupAcknowledgeCode", "oneof")
		suite.ValidateError(err, "NumberOfTransactionSetsIncluded", "min")
		suite.ValidateError(err, "NumberOfReceivedTransactionSets", "min")
		suite.ValidateError(err, "NumberOfAcceptedTransactionSets", "min")
		suite.ValidateErrorLen(err, 4)
	})

	suite.T().Run("validate failure max", func(t *testing.T) {
		// length of characters are more than max
		ak9 := AK9{
			FunctionalGroupAcknowledgeCode:      "AA",
			NumberOfTransactionSetsIncluded:     1000000,
			NumberOfReceivedTransactionSets:     1000000,
			NumberOfAcceptedTransactionSets:     1000000,
			FunctionalGroupSyntaxErrorCodeAK905: "AAAA",
			FunctionalGroupSyntaxErrorCodeAK906: "BBBB",
			FunctionalGroupSyntaxErrorCodeAK907: "CCCC",
			FunctionalGroupSyntaxErrorCodeAK908: "DDDD",
			FunctionalGroupSyntaxErrorCodeAK909: "EEEE",
		}

		err := suite.validator.Struct(ak9)
		suite.ValidateError(err, "FunctionalGroupAcknowledgeCode", "oneof")
		suite.ValidateError(err, "NumberOfTransactionSetsIncluded", "max")
		suite.ValidateError(err, "NumberOfReceivedTransactionSets", "max")
		suite.ValidateError(err, "NumberOfAcceptedTransactionSets", "max")
		suite.ValidateError(err, "FunctionalGroupSyntaxErrorCodeAK905", "max")
		suite.ValidateError(err, "FunctionalGroupSyntaxErrorCodeAK906", "max")
		suite.ValidateError(err, "FunctionalGroupSyntaxErrorCodeAK907", "max")
		suite.ValidateError(err, "FunctionalGroupSyntaxErrorCodeAK908", "max")
		suite.ValidateError(err, "FunctionalGroupSyntaxErrorCodeAK909", "max")
		suite.ValidateErrorLen(err, 9)
	})

	suite.T().Run("validate failure min", func(t *testing.T) {
		// length of characters are less than min
		ak9 := AK9{
			FunctionalGroupAcknowledgeCode:  "",
			NumberOfTransactionSetsIncluded: 0,
			NumberOfReceivedTransactionSets: 0,
			NumberOfAcceptedTransactionSets: 0,
		}

		err := suite.validator.Struct(ak9)
		suite.ValidateError(err, "FunctionalGroupAcknowledgeCode", "oneof")
		suite.ValidateError(err, "NumberOfTransactionSetsIncluded", "min")
		suite.ValidateError(err, "NumberOfReceivedTransactionSets", "min")
		suite.ValidateError(err, "NumberOfAcceptedTransactionSets", "min")
		suite.ValidateErrorLen(err, 4)
	})
}

func (suite *SegmentSuite) TestStringArrayAK9() {
	suite.T().Run("string array all fields", func(t *testing.T) {
		validAK9 := AK9{
			FunctionalGroupAcknowledgeCode:      "A",
			NumberOfTransactionSetsIncluded:     1,
			NumberOfReceivedTransactionSets:     2,
			NumberOfAcceptedTransactionSets:     3,
			FunctionalGroupSyntaxErrorCodeAK905: "AAA",
			FunctionalGroupSyntaxErrorCodeAK906: "BBB",
			FunctionalGroupSyntaxErrorCodeAK907: "CCC",
			FunctionalGroupSyntaxErrorCodeAK908: "DDD",
			FunctionalGroupSyntaxErrorCodeAK909: "EEE",
		}
		arrayValidAK9 := []string{"AK9", "A", "1", "2", "3", "AAA", "BBB", "CCC", "DDD", "EEE"}
		suite.Equal(arrayValidAK9, validAK9.StringArray())
	})

	suite.T().Run("string array only required fields", func(t *testing.T) {
		validOptionalAK9 := AK9{
			FunctionalGroupAcknowledgeCode:  "A",
			NumberOfTransactionSetsIncluded: 1,
			NumberOfReceivedTransactionSets: 2,
			NumberOfAcceptedTransactionSets: 3,
		}
		arrayValidOptionalAK9 := []string{"AK9", "A", "1", "2", "3", "", "", "", "", ""}
		suite.Equal(arrayValidOptionalAK9, validOptionalAK9.StringArray())
	})
}

func (suite *SegmentSuite) TestParseAK9() {
	suite.T().Run("parse success all fields", func(t *testing.T) {
		arrayValidAK9 := []string{"A", "1", "2", "3", "AAA", "BBB", "CCC", "DDD", "EEE"}
		expectedAK9 := AK9{
			FunctionalGroupAcknowledgeCode:      "A",
			NumberOfTransactionSetsIncluded:     1,
			NumberOfReceivedTransactionSets:     2,
			NumberOfAcceptedTransactionSets:     3,
			FunctionalGroupSyntaxErrorCodeAK905: "AAA",
			FunctionalGroupSyntaxErrorCodeAK906: "BBB",
			FunctionalGroupSyntaxErrorCodeAK907: "CCC",
			FunctionalGroupSyntaxErrorCodeAK908: "DDD",
			FunctionalGroupSyntaxErrorCodeAK909: "EEE",
		}

		var validAK9 AK9
		err := validAK9.Parse(arrayValidAK9)
		if suite.NoError(err) {
			suite.Equal(expectedAK9, validAK9)
		}
	})

	suite.T().Run("parse success only required fields", func(t *testing.T) {
		arrayValidOptionalAK9 := []string{"A", "1", "2", "3", "", "", "", "", ""}
		expectedOptionalAK9 := AK9{
			FunctionalGroupAcknowledgeCode:  "A",
			NumberOfTransactionSetsIncluded: 1,
			NumberOfReceivedTransactionSets: 2,
			NumberOfAcceptedTransactionSets: 3,
		}

		var validOptionalAK9 AK9
		err := validOptionalAK9.Parse(arrayValidOptionalAK9)
		if suite.NoError(err) {
			suite.Equal(expectedOptionalAK9, validOptionalAK9)
		}
	})

	suite.T().Run("wrong number of elements", func(t *testing.T) {
		badArrayAK9 := []string{"A", "abc"}
		var badAK9 AK9
		err := badAK9.Parse(badArrayAK9)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of elements")
		}
	})

	suite.T().Run("parse fails for invalid ints", func(t *testing.T) {
		var validOptionalAK9 AK9
		arrayInvalidIntsAK9 := []string{"A", "g", "2", "3", "", "", "", "", ""}

		err := validOptionalAK9.Parse(arrayInvalidIntsAK9)
		if suite.Error(err) {
			suite.Contains(err.Error(), "invalid syntax")
		}

		arrayInvalidIntsAK9[1] = "1"
		arrayInvalidIntsAK9[2] = "AAA"

		err = validOptionalAK9.Parse(arrayInvalidIntsAK9)
		if suite.Error(err) {
			suite.Contains(err.Error(), "invalid syntax")
		}
		arrayInvalidIntsAK9[2] = "2"
		arrayInvalidIntsAK9[3] = "3.0"

		err = validOptionalAK9.Parse(arrayInvalidIntsAK9)
		if suite.Error(err) {
			suite.Contains(err.Error(), "invalid syntax")
		}
	})
}
