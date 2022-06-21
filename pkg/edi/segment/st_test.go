package edisegment

func (suite *SegmentSuite) TestValidateST() {
	validST := ST{
		TransactionSetIdentifierCode: "858",
		TransactionSetControlNumber:  "ABCDE",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validST)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		st := ST{
			TransactionSetIdentifierCode: "123", // eq
			TransactionSetControlNumber:  "123", // min
		}

		err := suite.validator.Struct(st)
		suite.ValidateError(err, "TransactionSetIdentifierCode", "oneof")
		suite.ValidateError(err, "TransactionSetControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		st := validST
		st.TransactionSetControlNumber = "1234567890" // max

		err := suite.validator.Struct(st)
		suite.ValidateError(err, "TransactionSetControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
