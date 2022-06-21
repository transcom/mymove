package edisegment

func (suite *SegmentSuite) TestValidateBX() {
	validBX := BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: "1234-1234-5",
		StandardCarrierAlphaCode:     "TEST",
		ShipmentQualifier:            "4",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validBX)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		bx := BX{
			TransactionSetPurposeCode:    "01",    // eq
			TransactionMethodTypeCode:    "K",     // eq
			ShipmentMethodOfPayment:      "QQ",    // eq
			ShipmentIdentificationNumber: "",      // min
			StandardCarrierAlphaCode:     "TEST2", // alpha
			WeightUnitCode:               "1",     // isdefault
			ShipmentQualifier:            "5",     // eq
		}

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "TransactionSetPurposeCode", "eq")
		suite.ValidateError(err, "TransactionMethodTypeCode", "eq")
		suite.ValidateError(err, "ShipmentMethodOfPayment", "eq")
		suite.ValidateError(err, "ShipmentIdentificationNumber", "min")
		suite.ValidateError(err, "StandardCarrierAlphaCode", "alpha")
		suite.ValidateError(err, "WeightUnitCode", "isdefault")
		suite.ValidateError(err, "ShipmentQualifier", "eq")
		suite.ValidateErrorLen(err, 7)
	})

	suite.Run("validate failure 2", func() {
		bx := validBX
		bx.StandardCarrierAlphaCode = "T" // min

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "StandardCarrierAlphaCode", "min")
		suite.ValidateErrorLen(err, 1)
	})

	suite.Run("validate failure 3", func() {
		bx := validBX
		bx.ShipmentIdentificationNumber = "A123456789012345678901234567890" // max
		bx.StandardCarrierAlphaCode = "TESTING"                             // max

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "ShipmentIdentificationNumber", "max")
		suite.ValidateError(err, "StandardCarrierAlphaCode", "max")
		suite.ValidateErrorLen(err, 2)
	})
}
