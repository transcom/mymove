package edisegment

func (suite *SegmentSuite) TestValidateAK4() {
	suite.Run("validate success all fields", func() {
		validAK4 := AK4{
			PositionInSegment:                       1,
			ElementPositionInSegment:                1,
			ComponentDataElementPositionInComposite: 11,
			DataElementReferenceNumber:              1111,
			DataElementSyntaxErrorCode:              "ABC",
			CopyOfBadDataElement:                    "Bad data element",
		}
		err := suite.validator.Struct(validAK4)
		suite.NoError(err)
	})

	suite.Run("validate success only required fields", func() {
		validOptionalAK4 := AK4{
			PositionInSegment:          1,
			ElementPositionInSegment:   1,
			DataElementSyntaxErrorCode: "ABC",
		}
		err := suite.validator.Struct(validOptionalAK4)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		ak4 := AK4{
			PositionInSegment:                       -1, // min
			ElementPositionInSegment:                -1, // min
			ComponentDataElementPositionInComposite: -1, // min
			DataElementReferenceNumber:              -1, // min
			DataElementSyntaxErrorCode:              "", // min
		}

		err := suite.validator.Struct(ak4)
		suite.ValidateError(err, "PositionInSegment", "min")
		suite.ValidateError(err, "ElementPositionInSegment", "min")
		suite.ValidateError(err, "ComponentDataElementPositionInComposite", "min")
		suite.ValidateError(err, "DataElementReferenceNumber", "min")
		suite.ValidateError(err, "DataElementSyntaxErrorCode", "min")
		suite.ValidateErrorLen(err, 5)
	})

	suite.Run("validate failure 2", func() {
		ak4 := AK4{
			PositionInSegment:                       100,                                                                                                    // max
			ElementPositionInSegment:                100,                                                                                                    // max
			ComponentDataElementPositionInComposite: 100,                                                                                                    // max
			DataElementReferenceNumber:              10000,                                                                                                  // max
			DataElementSyntaxErrorCode:              "THIS",                                                                                                 // max
			CopyOfBadDataElement:                    "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", // max
		}

		err := suite.validator.Struct(ak4)
		suite.ValidateError(err, "PositionInSegment", "max")
		suite.ValidateError(err, "ElementPositionInSegment", "max")
		suite.ValidateError(err, "ComponentDataElementPositionInComposite", "max")
		suite.ValidateError(err, "DataElementReferenceNumber", "max")
		suite.ValidateError(err, "DataElementSyntaxErrorCode", "max")
		suite.ValidateError(err, "CopyOfBadDataElement", "max")
		suite.ValidateErrorLen(err, 6)
	})
}

func (suite *SegmentSuite) TestStringArrayAK4() {
	suite.Run("string array all fields", func() {
		validAK4 := AK4{
			PositionInSegment:                       1,
			ElementPositionInSegment:                1,
			ComponentDataElementPositionInComposite: 11,
			DataElementReferenceNumber:              1111,
			DataElementSyntaxErrorCode:              "ABC",
			CopyOfBadDataElement:                    "Bad data element",
		}
		arrayValidAK4 := []string{"AK4", "1", "1", "11", "1111", "ABC", "Bad data element"}
		suite.Equal(arrayValidAK4, validAK4.StringArray())
	})

	suite.Run("string array only required fields", func() {
		validOptionalAK4 := AK4{
			PositionInSegment:          1,
			ElementPositionInSegment:   1,
			DataElementSyntaxErrorCode: "ABC",
		}
		arrayValidOptionalAK4 := []string{"AK4", "1", "1", "", "", "ABC", ""}
		suite.Equal(arrayValidOptionalAK4, validOptionalAK4.StringArray())
	})
}

func (suite *SegmentSuite) TestParseAK4() {
	suite.Run("parse success all fields", func() {
		arrayValidAK4 := []string{"1", "1", "11", "1111", "ABC", "Bad data element"}
		expectedAK4 := AK4{
			PositionInSegment:                       1,
			ElementPositionInSegment:                1,
			ComponentDataElementPositionInComposite: 11,
			DataElementReferenceNumber:              1111,
			DataElementSyntaxErrorCode:              "ABC",
			CopyOfBadDataElement:                    "Bad data element",
		}

		var validAK4 AK4
		err := validAK4.Parse(arrayValidAK4)
		if suite.NoError(err) {
			suite.Equal(expectedAK4, validAK4)
		}
	})

	suite.Run("parse success only required fields", func() {
		arrayValidOptionalAK4 := []string{"1", "1", "", "", "ABC", ""}
		expectedOptionalAK4 := AK4{
			PositionInSegment:          1,
			ElementPositionInSegment:   1,
			DataElementSyntaxErrorCode: "ABC",
		}

		var validOptionalAK4 AK4
		err := validOptionalAK4.Parse(arrayValidOptionalAK4)
		if suite.NoError(err) {
			suite.Equal(expectedOptionalAK4, validOptionalAK4)
		}
	})

	suite.Run("wrong number of fields", func() {
		badArrayAK4 := []string{"1", "1"}
		var badAK4 AK4
		err := badAK4.Parse(badArrayAK4)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})

	suite.Run("invalid integers", func() {
		// First four fields are integers that could fail conversion
		for i := 0; i < 4; i++ {
			badArrayAK4 := []string{"1", "1", "11", "1111", "ABC", "Bad data element"}
			badArrayAK4[i] = "abc" // can't be converted to an integer
			var badAK4 AK4
			err := badAK4.Parse(badArrayAK4)
			if suite.Error(err) {
				suite.Contains(err.Error(), "invalid syntax")
			}
		}
	})
}
