package edisegment

func (suite *SegmentSuite) TestValidateFA2() {
	validFA2 := FA2{
		BreakdownStructureDetailCode: "TA",
		FinancialInformationCode:     "307",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validFA2)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		fa2 := FA2{
			BreakdownStructureDetailCode: "XX", // eq
			FinancialInformationCode:     "",   // min
		}

		err := suite.validator.Struct(fa2)
		suite.ValidateError(err, "BreakdownStructureDetailCode", "eq")
		suite.ValidateError(err, "FinancialInformationCode", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		fa2 := validFA2
		fa2.FinancialInformationCode = "123456789012345678901234567890123456789012345678901234567890123456789012345678901" // max

		err := suite.validator.Struct(fa2)
		suite.ValidateError(err, "FinancialInformationCode", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
