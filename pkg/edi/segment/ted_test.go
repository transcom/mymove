package edisegment

import (
	"strings"
	"testing"
)

func (suite *SegmentSuite) TestValidateTED() {
	suite.T().Run("validate success all fields", func(t *testing.T) {
		validTED := TED{
			ApplicationErrorConditionCode: "007",
			FreeFormMessage:               "free form message",
		}
		err := suite.validator.Struct(validTED)
		suite.NoError(err)
	})

	suite.T().Run("validate success only required fields", func(t *testing.T) {
		validOptionalTED := TED{
			ApplicationErrorConditionCode: "007",
		}
		err := suite.validator.Struct(validOptionalTED)
		suite.NoError(err)
	})

	suite.T().Run("validate failure", func(t *testing.T) {
		ted := TED{
			ApplicationErrorConditionCode: "123",                   // oneof
			FreeFormMessage:               strings.Repeat("x", 61), // max
		}

		err := suite.validator.Struct(ted)
		suite.ValidateError(err, "ApplicationErrorConditionCode", "oneof")
		suite.ValidateError(err, "FreeFormMessage", "max")
		suite.ValidateErrorLen(err, 2)
	})
}

func (suite *SegmentSuite) TestStringArrayTED() {
	suite.T().Run("string array all fields", func(t *testing.T) {
		validTED := TED{
			ApplicationErrorConditionCode: "007",
			FreeFormMessage:               "free form message",
		}
		arrayValidTED := []string{"TED", "007", "free form message"}
		suite.Equal(arrayValidTED, validTED.StringArray())
	})

	suite.T().Run("string array only required fields", func(t *testing.T) {
		validOptionalTED := TED{
			ApplicationErrorConditionCode: "007",
		}
		arrayValidOptionalTED := []string{"TED", "007", ""}
		suite.Equal(arrayValidOptionalTED, validOptionalTED.StringArray())
	})
}

func (suite *SegmentSuite) TestParseTED() {
	suite.T().Run("parse success all fields", func(t *testing.T) {
		arrayValidTED := []string{"007", "free form message"}
		expectedTED := TED{
			ApplicationErrorConditionCode: "007",
			FreeFormMessage:               "free form message",
		}

		var validTED TED
		err := validTED.Parse(arrayValidTED)
		if suite.NoError(err) {
			suite.Equal(expectedTED, validTED)
		}
	})

	suite.T().Run("parse success only required fields", func(t *testing.T) {
		arrayValidOptionalTED := []string{"007", ""}
		expectedOptionalTED := TED{
			ApplicationErrorConditionCode: "007",
		}

		var validOptionalTED TED
		err := validOptionalTED.Parse(arrayValidOptionalTED)
		if suite.NoError(err) {
			suite.Equal(expectedOptionalTED, validOptionalTED)
		}
	})

	suite.T().Run("wrong number of fields", func(t *testing.T) {
		badArrayTED := []string{"007", "hello", "world"}
		var badTED TED
		err := badTED.Parse(badArrayTED)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})
}
