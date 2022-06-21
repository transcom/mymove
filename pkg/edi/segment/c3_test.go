package edisegment

func (suite *SegmentSuite) TestValidateC3() {
	validC3 := C3{
		CurrencyCodeC301: "USD",
		ExchangeRate:     "x",
		CurrencyCodeC303: "USD",
		CurrencyCodeC304: "EUR",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validC3)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		c3 := C3{
			CurrencyCodeC301: "TTTT", // max
			ExchangeRate:     "K",    // none
			CurrencyCodeC303: "QQQQ", // max
			CurrencyCodeC304: "RRRR", // max
		}

		err := suite.validator.Struct(c3)
		suite.ValidateError(err, "CurrencyCodeC301", "max")
		suite.ValidateError(err, "CurrencyCodeC303", "max")
		suite.ValidateError(err, "CurrencyCodeC304", "max")
		suite.ValidateErrorLen(err, 3)
	})
}
