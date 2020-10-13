package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidatePER() {
	validPERDefault := PER{
		ContactFunctionCode: "IC",
	}

	validPERAll := PER{
		ContactFunctionCode:          "IC",
		Name:                         "Cross Dock",
		CommunicationNumberQualifier: "TE",
		CommunicationNumber:          "5551234567",
	}

	suite.T().Run("validate success default", func(t *testing.T) {
		err := suite.validator.Struct(validPERDefault)
		suite.NoError(err)
	})

	suite.T().Run("validate success all", func(t *testing.T) {
		err := suite.validator.Struct(validPERAll)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		per := PER{
			Name:                         "Cross Dock Testing With Too Long a Name",
			CommunicationNumberQualifier: "BX",
			CommunicationNumber:          "111111111111111111111111111111111111111115551234567",
		}

		err := suite.validator.Struct(per)
		suite.ValidateError(err, "ContactFunctionCode", "required")
		suite.ValidateError(err, "Name", "max")
		suite.ValidateError(err, "CommunicationNumberQualifier", "eq")
		suite.ValidateError(err, "CommunicationNumber", "max")
		suite.ValidateErrorLen(err, 4)
	})
}
