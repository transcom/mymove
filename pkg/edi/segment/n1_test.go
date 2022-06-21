package edisegment

func (suite *SegmentSuite) TestValidateN1() {
	validN1 := N1{
		EntityIdentifierCode:        "SF",
		Name:                        "ABC",
		IdentificationCodeQualifier: "27",
		IdentificationCode:          "XX",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validN1)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		n1 := N1{
			EntityIdentifierCode:        "XX", // oneof
			Name:                        "",   // min
			IdentificationCodeQualifier: "27", // required_with
		}

		err := suite.validator.Struct(n1)
		suite.ValidateError(err, "EntityIdentifierCode", "oneof")
		suite.ValidateError(err, "Name", "min")
		suite.ValidateError(err, "IdentificationCode", "required_with")
		suite.ValidateErrorLen(err, 3)
	})

	suite.Run("validate failure 2", func() {
		n1 := validN1
		n1.Name = "1234567890123456789012345678901234567890123456789012345678901" // max
		n1.IdentificationCodeQualifier = "19"                                     // oneof
		n1.IdentificationCode = "1"                                               // min

		err := suite.validator.Struct(n1)
		suite.ValidateError(err, "Name", "max")
		suite.ValidateError(err, "IdentificationCodeQualifier", "oneof")
		suite.ValidateError(err, "IdentificationCode", "min")
		suite.ValidateErrorLen(err, 3)
	})

	suite.Run("validate failure 3", func() {
		n1 := validN1
		n1.IdentificationCode = "123456789012345678901234567890123456789012345678901234567890123456789012345678901" // max

		err := suite.validator.Struct(n1)
		suite.ValidateError(err, "IdentificationCode", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
