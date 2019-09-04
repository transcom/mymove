package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateIEA() {
	validIEA := IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         1,
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validIEA)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		iea := IEA{
			NumberOfIncludedFunctionalGroups: 2, // eq
			InterchangeControlNumber:         0, // min
		}

		err := suite.validator.Struct(iea)
		suite.ValidateError(err, "NumberOfIncludedFunctionalGroups", "eq")
		suite.ValidateError(err, "InterchangeControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		iea := validIEA
		iea.InterchangeControlNumber = 1000000000 // max

		err := suite.validator.Struct(iea)
		suite.ValidateError(err, "InterchangeControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
