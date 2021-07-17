package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateFA2() {
	validFA2 := FA2{
		BreakdownStructureDetailCode: "TA",
		FinancialInformationCode:     "307",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validFA2)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		fa2 := FA2{
			BreakdownStructureDetailCode: "XX", // eq
			FinancialInformationCode:     "",   // min
		}

		err := suite.validator.Struct(fa2)
		suite.ValidateError(err, "BreakdownStructureDetailCode", "oneof")
		suite.ValidateError(err, "FinancialInformationCode", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		fa2 := validFA2
		fa2.FinancialInformationCode = "123456789012345678901234567890123456789012345678901234567890123456789012345678901" // max

		err := suite.validator.Struct(fa2)
		suite.ValidateError(err, "FinancialInformationCode", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
