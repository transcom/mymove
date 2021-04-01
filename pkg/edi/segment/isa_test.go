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
		InterchangeSenderID:               fmt.Sprintf("%-15s", "MILMOVE"),
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

	altValidISA := ISA{
		AuthorizationInformationQualifier: "00",
		SecurityInformationQualifier:      "00",
		InterchangeSenderIDQualifier:      "12",
		InterchangeSenderID:               fmt.Sprintf("%-15s", "8004171844"),
		InterchangeReceiverIDQualifier:    "ZZ",
		InterchangeReceiverID:             fmt.Sprintf("%-15s", "MILMOVE"),
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
		err := suite.validator.Struct(altValidISA)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		isa := ISA{
			AuthorizationInformationQualifier: "11",                           // eq
			AuthorizationInformation:          "1111111111",                   // eq
			SecurityInformationQualifier:      "11",                           // eq
			SecurityInformation:               "1111111111",                   // eq
			InterchangeSenderIDQualifier:      "QQ",                           // oneof
			InterchangeSenderID:               fmt.Sprintf("%-17s", "ABCDEF"), // eq
			InterchangeReceiverIDQualifier:    "15",                           // oneof
			InterchangeReceiverID:             fmt.Sprintf("%-1s", "13"),      // len
			InterchangeDate:                   "190933",                       // datetime
			InterchangeTime:                   "344",                          // datetime
			InterchangeControlStandards:       "Q",                            // eq
			InterchangeControlVersionNumber:   "00403",                        // eq
			InterchangeControlNumber:          0,                              // min
			AcknowledgementRequested:          5,                              // eq
			UsageIndicator:                    "Q",                            // oneof
			ComponentElementSeparator:         ",",                            // eq
		}

		err := suite.validator.Struct(isa)
		suite.ValidateError(err, "AuthorizationInformationQualifier", "eq")
		suite.ValidateError(err, "AuthorizationInformation", "eq")
		suite.ValidateError(err, "SecurityInformationQualifier", "eq")
		suite.ValidateError(err, "SecurityInformation", "eq")
		suite.ValidateError(err, "InterchangeSenderIDQualifier", "oneof")
		suite.ValidateError(err, "InterchangeSenderID", "len")
		suite.ValidateError(err, "InterchangeReceiverIDQualifier", "oneof")
		suite.ValidateError(err, "InterchangeReceiverID", "len")
		suite.ValidateError(err, "InterchangeDate", "datetime")
		suite.ValidateError(err, "InterchangeTime", "datetime")
		suite.ValidateError(err, "InterchangeControlStandards", "eq")
		suite.ValidateError(err, "InterchangeControlVersionNumber", "eq")
		suite.ValidateError(err, "InterchangeControlNumber", "min")
		suite.ValidateError(err, "AcknowledgementRequested", "oneof")
		suite.ValidateError(err, "UsageIndicator", "oneof")
		suite.ValidateError(err, "ComponentElementSeparator", "oneof")
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
