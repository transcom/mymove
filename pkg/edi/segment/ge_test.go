package edisegment

func (suite *SegmentSuite) TestValidateGE() {
	validGE := GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              1234567,
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validGE)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		ge := GE{
			NumberOfTransactionSetsIncluded: 2, // eq
			GroupControlNumber:              0, // min
		}

		err := suite.validator.Struct(ge)
		suite.ValidateError(err, "NumberOfTransactionSetsIncluded", "eq")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		ge := validGE
		ge.GroupControlNumber = 1000000000 // max

		err := suite.validator.Struct(ge)
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
