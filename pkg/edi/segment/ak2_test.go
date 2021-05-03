package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateAK2() {
	validAK2 := AK2{
		TransactionSetIdentifierCode: "858",
		TransactionSetControlNumber:  "ABCDE",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validAK2)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		ak2 := AK2{
			TransactionSetIdentifierCode: "123", // eq
			TransactionSetControlNumber:  "123", // min
		}

		err := suite.validator.Struct(ak2)
		suite.ValidateError(err, "TransactionSetIdentifierCode", "eq")
		suite.ValidateError(err, "TransactionSetControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		ak2 := validAK2
		ak2.TransactionSetControlNumber = "1234567890" // max

		err := suite.validator.Struct(ak2)
		suite.ValidateError(err, "TransactionSetControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
