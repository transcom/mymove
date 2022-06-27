package edisegment

func (suite *SegmentSuite) TestValidateAK3() {
	validAK3 := AK3{
		SegmentIDCode:                   "ID",
		SegmentPositionInTransactionSet: 12345,
		LoopIdentifierCode:              "CODE",
		SegmentSyntaxErrorCode:          "ERR",
	}

	suite.Run("validate success all fields", func() {
		err := suite.validator.Struct(validAK3)
		suite.NoError(err)
	})

	altValidAK3 := AK3{
		SegmentIDCode:                   "ID",
		SegmentPositionInTransactionSet: 12345,
	}

	suite.Run("validate success with only required fields", func() {
		err := suite.validator.Struct(altValidAK3)
		suite.NoError(err)
	})

	suite.Run("validate failure for min", func() {
		ak3 := AK3{
			SegmentIDCode:                   "", // min
			SegmentPositionInTransactionSet: 0,  // min
		}

		err := suite.validator.Struct(ak3)
		suite.ValidateError(err, "SegmentIDCode", "min")
		suite.ValidateError(err, "SegmentPositionInTransactionSet", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure for max", func() {
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

func (suite *SegmentSuite) TestStringArrayAK3() {
	suite.Run("string array all fields", func() {
		validAK3 := AK3{
			SegmentIDCode:                   "ID",
			SegmentPositionInTransactionSet: 12345,
			LoopIdentifierCode:              "CODE",
			SegmentSyntaxErrorCode:          "ERR",
		}
		arrayValidAK3 := []string{"AK3", "ID", "12345", "CODE", "ERR"}
		suite.Equal(arrayValidAK3, validAK3.StringArray())
	})

	suite.Run("string array only required fields", func() {
		validAK3 := AK3{
			SegmentIDCode:                   "ID",
			SegmentPositionInTransactionSet: 12345,
		}
		arrayValidAK3 := []string{"AK3", "ID", "12345", "", ""}
		suite.Equal(arrayValidAK3, validAK3.StringArray())
	})
}

func (suite *SegmentSuite) TestParseAK3() {
	suite.Run("parse success all fields", func() {
		arrayValidAK3 := []string{"ID", "12345", "CODE", "ERR"}

		expectedAK3 := AK3{
			SegmentIDCode:                   "ID",
			SegmentPositionInTransactionSet: 12345,
			LoopIdentifierCode:              "CODE",
			SegmentSyntaxErrorCode:          "ERR",
		}

		var validAK3 AK3
		err := validAK3.Parse(arrayValidAK3)
		if suite.NoError(err) {
			suite.Equal(expectedAK3, validAK3)
		}
	})

	suite.Run("parse success on required fields", func() {
		arrayValidAK3 := []string{"ID", "12345", "", ""}

		expectedAK3 := AK3{
			SegmentIDCode:                   "ID",
			SegmentPositionInTransactionSet: 12345,
		}

		var validAK3 AK3
		err := validAK3.Parse(arrayValidAK3)
		if suite.NoError(err) {
			suite.Equal(expectedAK3, validAK3)
		}
	})

	suite.Run("wrong number of elements", func() {
		badArrayAK3 := []string{"11"}
		var badAK3 AK3
		err := badAK3.Parse(badArrayAK3)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of elements")
		}
	})

	suite.Run("wrong number of elements greater than max", func() {
		badArrayAK3 := []string{"11", "12", "by", "goo", "fooz"}
		var badAK3 AK3
		err := badAK3.Parse(badArrayAK3)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of elements")
		}
	})

	suite.Run("fail when SegmentPositionInTransactionSet not a valid int", func() {
		badArrayAK3 := []string{"ID", "12345.4", "", ""}

		var badAK3 AK3
		err := badAK3.Parse(badArrayAK3)
		if suite.Error(err) {
			suite.Contains(err.Error(), "invalid syntax")
		}
	})
}
