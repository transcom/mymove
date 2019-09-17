package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateBX() {
	validBX := BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: "A12345",
		StandardCarrierAlphaCode:     "TEST",
		ShipmentQualifier:            "4",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validBX)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		bx := BX{
			TransactionSetPurposeCode:    "01",      // eq
			TransactionMethodTypeCode:    "K",       // eq
			ShipmentMethodOfPayment:      "QQ",      // eq
			ShipmentIdentificationNumber: "A-12345", // alphanum
			StandardCarrierAlphaCode:     "TEST2",   // alpha
			WeightUnitCode:               "1",       // isdefault
			ShipmentQualifier:            "5",       // eq
		}

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "TransactionSetPurposeCode", "eq")
		suite.ValidateError(err, "TransactionMethodTypeCode", "eq")
		suite.ValidateError(err, "ShipmentMethodOfPayment", "eq")
		suite.ValidateError(err, "ShipmentIdentificationNumber", "alphanum")
		suite.ValidateError(err, "StandardCarrierAlphaCode", "alpha")
		suite.ValidateError(err, "WeightUnitCode", "isdefault")
		suite.ValidateError(err, "ShipmentQualifier", "eq")
		suite.ValidateErrorLen(err, 7)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		bx := validBX
		bx.ShipmentIdentificationNumber = "" // alphanum (precedence over min)
		bx.StandardCarrierAlphaCode = "T"    // min

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "ShipmentIdentificationNumber", "alphanum")
		suite.ValidateError(err, "StandardCarrierAlphaCode", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		bx := validBX
		bx.ShipmentIdentificationNumber = "A123456789012345678901234567890" // max
		bx.StandardCarrierAlphaCode = "TESTING"                             // max

		err := suite.validator.Struct(bx)
		suite.ValidateError(err, "ShipmentIdentificationNumber", "max")
		suite.ValidateError(err, "StandardCarrierAlphaCode", "max")
		suite.ValidateErrorLen(err, 2)
	})
}
