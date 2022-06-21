package edisegment

func (suite *SegmentSuite) TestValidateN4() {
	validN4 := N4{
		CityName:            "Augusta",
		StateOrProvinceCode: "GA",
		PostalCode:          "30907",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validN4)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		n4 := N4{
			CityName:            "A",       // min
			StateOrProvinceCode: "Georgia", // len
			PostalCode:          "27",      // min
			CountryCode:         "U",       // min
			LocationQualifier:   "ABC",     // isdefault
			LocationIdentifier:  "XYZ",     // isdefault
		}

		err := suite.validator.Struct(n4)
		suite.ValidateError(err, "CityName", "min")
		suite.ValidateError(err, "StateOrProvinceCode", "len")
		suite.ValidateError(err, "PostalCode", "min")
		suite.ValidateError(err, "CountryCode", "min")
		suite.ValidateError(err, "LocationQualifier", "isdefault")
		suite.ValidateError(err, "LocationIdentifier", "isdefault")
		suite.ValidateErrorLen(err, 6)
	})

	suite.Run("validate failure 2", func() {
		n4 := validN4
		n4.CityName = "1234567890123456789012345678901" // max
		n4.PostalCode = "01234567890123456"             // max
		n4.CountryCode = "0123"                         // max

		err := suite.validator.Struct(n4)
		suite.ValidateError(err, "CityName", "max")
		suite.ValidateError(err, "PostalCode", "max")
		suite.ValidateError(err, "CountryCode", "max")
		suite.ValidateErrorLen(err, 3)
	})
}
