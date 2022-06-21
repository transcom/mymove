package edisegment

func (suite *SegmentSuite) TestValidateIEA() {
	validIEA := IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         1,
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validIEA)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		iea := IEA{
			NumberOfIncludedFunctionalGroups: 2, // eq
			InterchangeControlNumber:         0, // min
		}

		err := suite.validator.Struct(iea)
		suite.ValidateError(err, "NumberOfIncludedFunctionalGroups", "eq")
		suite.ValidateError(err, "InterchangeControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		iea := validIEA
		iea.InterchangeControlNumber = 1000000000 // max

		err := suite.validator.Struct(iea)
		suite.ValidateError(err, "InterchangeControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
