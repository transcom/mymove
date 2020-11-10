package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateN9() {
	validN9Default := N9{
		ReferenceIdentificationQualifier: "DY",
		ReferenceIdentification:          "12345",
	}

	validN9All := N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          "XYZ",
		FreeFormDescription:              "Something",
		Date:                             "20190903",
	}

	suite.T().Run("validate success default", func(t *testing.T) {
		err := suite.validator.Struct(validN9Default)
		suite.NoError(err)
	})

	suite.T().Run("validate success all", func(t *testing.T) {
		err := suite.validator.Struct(validN9All)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		n9 := N9{
			ReferenceIdentificationQualifier: "XX",                                             // oneof
			ReferenceIdentification:          "",                                               // min
			FreeFormDescription:              "1234567890123456789012345678901234567890123456", // max
			Date:                             "20190933",                                       // datetime
		}

		err := suite.validator.Struct(n9)
		suite.ValidateError(err, "ReferenceIdentificationQualifier", "oneof")
		suite.ValidateError(err, "ReferenceIdentification", "min")
		suite.ValidateError(err, "FreeFormDescription", "max")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateErrorLen(err, 4)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		n9 := validN9All
		n9.ReferenceIdentification = "1234567890123456789012345678901" // max

		err := suite.validator.Struct(n9)
		suite.ValidateError(err, "ReferenceIdentification", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
