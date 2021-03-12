package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateOTI() {
	validOTI := OTI{
		ApplicationAcknowledgementCode:   "TA",
		ReferenceIdentificationQualifier: "BM",
		ReferenceIdentification:          "ABC",
		ApplicationSendersCode:           "MILMOVE",
		ApplicationReceiversCode:         "RECEIVER",
		Date:                             "20210311",
		Time:                             "1057",
		GroupControlNumber:               12345,
		TransactionSetControlNumber:      "ABCDE",
	}

	suite.T().Run("validate success all fields", func(t *testing.T) {
		err := suite.validator.Struct(validOTI)
		suite.NoError(err)
	})

	suite.T().Run("validate success only required fields", func(t *testing.T) {
		validOptionalOTI := OTI{
			ApplicationAcknowledgementCode:   "TA",
			ReferenceIdentificationQualifier: "BM",
			ReferenceIdentification:          "ABC",
		}
		err := suite.validator.Struct(validOptionalOTI)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		oti := OTI{
			ApplicationAcknowledgementCode:   "XX",       // oneof
			ReferenceIdentificationQualifier: "XX",       // oneof
			ReferenceIdentification:          "",         // min
			ApplicationSendersCode:           "X",        // min
			ApplicationReceiversCode:         "X",        // min
			Date:                             "20211311", // datetime
			Time:                             "2557",     // datetime
			GroupControlNumber:               -1,         // min
			TransactionSetControlNumber:      "ABC",      // min
		}

		err := suite.validator.Struct(oti)
		suite.ValidateError(err, "ApplicationAcknowledgementCode", "oneof")
		suite.ValidateError(err, "ReferenceIdentificationQualifier", "oneof")
		suite.ValidateError(err, "ReferenceIdentification", "min")
		suite.ValidateError(err, "ApplicationSendersCode", "min")
		suite.ValidateError(err, "ApplicationReceiversCode", "min")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateError(err, "Time", "datetime")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateError(err, "TransactionSetControlNumber", "min")
		suite.ValidateErrorLen(err, 9)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		oti := validOTI
		oti.ReferenceIdentification = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" // max
		oti.ApplicationSendersCode = "MILMOVEMILMOVEMILMOVE"            // max
		oti.ApplicationReceiversCode = "RECEIVERRECEIVER"               // max
		oti.GroupControlNumber = 1000000000                             // max
		oti.TransactionSetControlNumber = "ABCDEABCDE"                  // max

		err := suite.validator.Struct(oti)
		suite.ValidateError(err, "ReferenceIdentification", "max")
		suite.ValidateError(err, "ApplicationSendersCode", "max")
		suite.ValidateError(err, "ApplicationReceiversCode", "max")
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateError(err, "TransactionSetControlNumber", "max")
		suite.ValidateErrorLen(err, 5)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		oti := validOTI
		oti.GroupControlNumber = 0 // required_with

		err := suite.validator.Struct(oti)
		suite.ValidateError(err, "GroupControlNumber", "required_with")
		suite.ValidateErrorLen(err, 1)
	})
}

func (suite *SegmentSuite) TestStringArrayOTI() {
	suite.T().Run("string array all fields", func(t *testing.T) {
		validOTI := OTI{
			ApplicationAcknowledgementCode:   "TA",
			ReferenceIdentificationQualifier: "BM",
			ReferenceIdentification:          "ABC",
			ApplicationSendersCode:           "MILMOVE",
			ApplicationReceiversCode:         "RECEIVER",
			Date:                             "20210311",
			Time:                             "1057",
			GroupControlNumber:               12345,
			TransactionSetControlNumber:      "ABCDE",
		}
		arrayValidOTI := []string{"OTI", "TA", "BM", "ABC", "MILMOVE", "RECEIVER", "20210311", "1057", "12345", "ABCDE"}
		suite.Equal(arrayValidOTI, validOTI.StringArray())
	})

	suite.T().Run("string array only required fields", func(t *testing.T) {
		validOptionalOTI := OTI{
			ApplicationAcknowledgementCode:   "TA",
			ReferenceIdentificationQualifier: "BM",
			ReferenceIdentification:          "ABC",
		}
		arrayValidOptionalOTI := []string{"OTI", "TA", "BM", "ABC", "", "", "", "", "", ""}
		suite.Equal(arrayValidOptionalOTI, validOptionalOTI.StringArray())
	})
}

func (suite *SegmentSuite) TestParseOTI() {
	suite.T().Run("parse success all fields", func(t *testing.T) {
		arrayValidOTI := []string{"TA", "BM", "ABC", "MILMOVE", "RECEIVER", "20210311", "1057", "12345", "ABCDE"}
		expectedOTI := OTI{
			ApplicationAcknowledgementCode:   "TA",
			ReferenceIdentificationQualifier: "BM",
			ReferenceIdentification:          "ABC",
			ApplicationSendersCode:           "MILMOVE",
			ApplicationReceiversCode:         "RECEIVER",
			Date:                             "20210311",
			Time:                             "1057",
			GroupControlNumber:               12345,
			TransactionSetControlNumber:      "ABCDE",
		}

		var validOTI OTI
		err := validOTI.Parse(arrayValidOTI)
		if suite.NoError(err) {
			suite.Equal(expectedOTI, validOTI)
		}
	})

	suite.T().Run("parse success only required fields", func(t *testing.T) {
		arrayValidOptionalOTI := []string{"TA", "BM", "ABC", "", "", "", "", "", ""}
		expectedOptionalOTI := OTI{
			ApplicationAcknowledgementCode:   "TA",
			ReferenceIdentificationQualifier: "BM",
			ReferenceIdentification:          "ABC",
		}

		var validOptionalOTI OTI
		err := validOptionalOTI.Parse(arrayValidOptionalOTI)
		if suite.NoError(err) {
			suite.Equal(expectedOptionalOTI, validOptionalOTI)
		}
	})

	suite.T().Run("wrong number of fields", func(t *testing.T) {
		badArrayOTI := []string{"TA", "BM"}
		var badOTI OTI
		err := badOTI.Parse(badArrayOTI)
		if suite.Error(err) {
			suite.Contains(err.Error(), "Wrong number of fields")
		}
	})

	suite.T().Run("invalid integers", func(t *testing.T) {
		badArrayOTI := []string{"TA", "BM", "ABC", "MILMOVE", "RECEIVER", "20210311", "1057", "A12345", "ABCDE"}
		var badOTI OTI
		err := badOTI.Parse(badArrayOTI)
		if suite.Error(err) {
			suite.Contains(err.Error(), "invalid syntax")
		}
	})
}
