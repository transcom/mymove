package edisegment

import (
	"fmt"
	"testing"
)

func (suite *SegmentSuite) TestValidateISA() {
	validISA := ISA{
		AuthorizationInformationQualifier: "00",
		AuthorizationInformation:          "0084182369",
		SecurityInformationQualifier:      "00",
		SecurityInformation:               "0000000000",
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               fmt.Sprintf("%-15s", "MYMOVE"),
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             fmt.Sprintf("%-15s", "8004171844"),
		InterchangeDate:                   "190903",
		InterchangeTime:                   "1644",
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          1,
		AcknowledgementRequested:          1,
		UsageIndicator:                    "T",
		ComponentElementSeparator:         "|",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validISA)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		isa := ISA{
			AuthorizationInformationQualifier: "11",                               // eq
			AuthorizationInformation:          "1111111111",                       // eq
			SecurityInformationQualifier:      "11",                               // eq
			SecurityInformation:               "1111111111",                       // eq
			InterchangeSenderIDQualifier:      "QQ",                               // eq
			InterchangeSenderID:               fmt.Sprintf("%-15s", "ABCDEF"),     // eq
			InterchangeReceiverIDQualifier:    "15",                               // eq
			InterchangeReceiverID:             fmt.Sprintf("%-15s", "1234566133"), // eq
			InterchangeDate:                   "190933",                           // timeformat
			InterchangeTime:                   "344",                              // timeformat
			InterchangeControlStandards:       "Q",                                // eq
			InterchangeControlVersionNumber:   "00403",                            // eq
			InterchangeControlNumber:          0,                                  // min
			AcknowledgementRequested:          5,                                  // eq
			UsageIndicator:                    "Q",                                // oneof
			ComponentElementSeparator:         ",",                                // eq
		}

		err := suite.validator.Struct(isa)
		suite.ValidateError(err, "AuthorizationInformationQualifier", "eq")
		suite.ValidateError(err, "AuthorizationInformation", "eq")
		suite.ValidateError(err, "SecurityInformationQualifier", "eq")
		suite.ValidateError(err, "SecurityInformation", "eq")
		suite.ValidateError(err, "InterchangeSenderIDQualifier", "eq")
		suite.ValidateError(err, "InterchangeSenderID", "eq")
		suite.ValidateError(err, "InterchangeReceiverIDQualifier", "eq")
		suite.ValidateError(err, "InterchangeReceiverID", "eq")
		suite.ValidateError(err, "InterchangeDate", "timeformat")
		suite.ValidateError(err, "InterchangeTime", "timeformat")
		suite.ValidateError(err, "InterchangeControlStandards", "eq")
		suite.ValidateError(err, "InterchangeControlVersionNumber", "eq")
		suite.ValidateError(err, "InterchangeControlNumber", "min")
		suite.ValidateError(err, "AcknowledgementRequested", "oneof")
		suite.ValidateError(err, "UsageIndicator", "oneof")
		suite.ValidateError(err, "ComponentElementSeparator", "eq")
		suite.ValidateErrorLen(err, 16)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		isa := validISA
		isa.InterchangeControlNumber = 1000000000 // max

		err := suite.validator.Struct(isa)
		suite.ValidateError(err, "InterchangeControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
