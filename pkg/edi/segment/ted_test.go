package edisegment

import (
	"strings"
)

func (suite *SegmentSuite) TestValidateTED() {
	suite.Run("validate success all fields", func() {
		validTED := TED{
			ApplicationErrorConditionCode: "007",
			FreeFormMessage:               "free form message",
		}
		err := suite.validator.Struct(validTED)
		suite.NoError(err)
	})

	suite.Run("validate success only required fields", func() {
		validOptionalTED := TED{
			ApplicationErrorConditionCode: "007",
		}
		err := suite.validator.Struct(validOptionalTED)
		suite.NoError(err)
	})

	suite.Run("validate failure", func() {
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
	suite.Run("string array all fields", func() {
		validTED := TED{
			ApplicationErrorConditionCode: "007",
			FreeFormMessage:               "free form message",
		}
		arrayValidTED := []string{"TED", "007", "free form message"}
		suite.Equal(arrayValidTED, validTED.StringArray())
	})

	suite.Run("string array only required fields", func() {
		validOptionalTED := TED{
			ApplicationErrorConditionCode: "007",
		}
		arrayValidOptionalTED := []string{"TED", "007", ""}
		suite.Equal(arrayValidOptionalTED, validOptionalTED.StringArray())
	})
}

func (suite *SegmentSuite) TestParseTED() {
	suite.Run("parse success all fields", func() {
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

	suite.Run("parse success only required fields", func() {
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

	suite.Run("wrong number of fields", func() {
		badArrayTED := []string{"007", "hello", "world"}
		var badTED TED
		err := badTED.Parse(badArrayTED)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})
}
