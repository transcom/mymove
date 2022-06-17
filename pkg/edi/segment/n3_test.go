package edisegment

func (suite *SegmentSuite) TestValidateN3() {
	validN3Line1Only := N3{
		AddressInformation1: "ABC",
	}

	validN3BothLines := N3{
		AddressInformation1: "ABC",
		AddressInformation2: "XYZ",
	}

	suite.Run("validate success line 1 only", func() {
		err := suite.validator.Struct(validN3Line1Only)
		suite.NoError(err)
	})

	suite.Run("validate success both lines", func() {
		err := suite.validator.Struct(validN3BothLines)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		n3 := N3{
			AddressInformation1: "",                                                         // min
			AddressInformation2: "12345678901234567890123456789012345678901234567890123456", // max
		}

		err := suite.validator.Struct(n3)
		suite.ValidateError(err, "AddressInformation1", "min")
		suite.ValidateError(err, "AddressInformation2", "max")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		n3 := N3{
			AddressInformation1: "12345678901234567890123456789012345678901234567890123456", // max
		}

		err := suite.validator.Struct(n3)
		suite.ValidateError(err, "AddressInformation1", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
