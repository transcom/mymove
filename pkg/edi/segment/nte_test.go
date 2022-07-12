package edisegment

func (suite *SegmentSuite) TestValidateNTE() {
	validNTEDefault := NTE{
		Description: "Something",
	}

	validNTEAll := NTE{
		NoteReferenceCode: "ABC",
		Description:       "Something Else",
	}

	suite.Run("validate success default", func() {
		err := suite.validator.Struct(validNTEDefault)
		suite.NoError(err)
	})

	suite.Run("validate success all", func() {
		err := suite.validator.Struct(validNTEAll)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		nte := NTE{
			NoteReferenceCode: "XX", // len
			Description:       "",   // min
		}

		err := suite.validator.Struct(nte)
		suite.ValidateError(err, "NoteReferenceCode", "len")
		suite.ValidateError(err, "Description", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		nte := validNTEAll
		nte.Description = "123456789012345678901234567890123456789012345678901234567890123456789012345678901" // max

		err := suite.validator.Struct(nte)
		suite.ValidateError(err, "Description", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
