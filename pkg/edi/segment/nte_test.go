package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateNTE() {
	validNTEDefault := NTE{
		Description: "Something",
	}

	validNTEAll := NTE{
		NoteReferenceCode: "ABC",
		Description:       "Something Else",
	}

	suite.T().Run("validate success default", func(t *testing.T) {
		err := suite.validator.Struct(validNTEDefault)
		suite.NoError(err)
	})

	suite.T().Run("validate success all", func(t *testing.T) {
		err := suite.validator.Struct(validNTEAll)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		nte := NTE{
			NoteReferenceCode: "XX", // len
			Description:       "",   // min
		}

		err := suite.validator.Struct(nte)
		suite.ValidateError(err, "NoteReferenceCode", "len")
		suite.ValidateError(err, "Description", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		nte := validNTEAll
		nte.Description = "123456789012345678901234567890123456789012345678901234567890123456789012345678901" // max

		err := suite.validator.Struct(nte)
		suite.ValidateError(err, "Description", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
