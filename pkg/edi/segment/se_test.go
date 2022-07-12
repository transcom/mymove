package edisegment

func (suite *SegmentSuite) TestValidateSE() {
	validSE := SE{
		NumberOfIncludedSegments:    12345,
		TransactionSetControlNumber: "ABCDE",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validSE)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		se := SE{
			NumberOfIncludedSegments:    0,     // min
			TransactionSetControlNumber: "ABC", // min
		}

		err := suite.validator.Struct(se)
		suite.ValidateError(err, "NumberOfIncludedSegments", "min")
		suite.ValidateError(err, "TransactionSetControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		se := SE{
			NumberOfIncludedSegments:    10000000000,  // max
			TransactionSetControlNumber: "1234567890", // max
		}

		err := suite.validator.Struct(se)
		suite.ValidateError(err, "NumberOfIncludedSegments", "max")
		suite.ValidateError(err, "TransactionSetControlNumber", "max")
		suite.ValidateErrorLen(err, 2)
	})
}
