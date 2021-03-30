package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) TestEdiErrorsAcknowledgementCodeError() {
	pr := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	ediError := models.EdiError{
		ID:               uuid.Must(uuid.NewV4()),
		PaymentRequestID: pr.ID,
		PaymentRequest:   pr,
	}
	prICN := models.PaymentRequestToInterchangeControlNumber{
		ID:                       uuid.Must(uuid.NewV4()),
		PaymentRequestID:         pr.ID,
		InterchangeControlNumber: 5,
	}
	suite.MustCreate(suite.DB(), &prICN)
	testCases := map[string]struct {
		ack          models.EdiErrorsAcknowledgementCodeError
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			ack: models.EdiErrorsAcknowledgementCodeError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDI997,
				EdiErrorID:                 ediError.ID,
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: prICN.ID,
				Code:                       "B",
				Description:                "EDI Error happened to field 99",
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			ack: models.EdiErrorsAcknowledgementCodeError{},
			expectedErrs: map[string][]string{
				"edi_error_id":                  {"EdiErrorID can not be blank."},
				"description":                   {"Code or Description must be present"},
				"code":                          {"Code or Description must be present"},
				"payment_request_id":            {"PaymentRequestID can not be blank."},
				"interchange_control_number_id": {"InterchangeControlNumberID can not be blank."},
				"editype":                       {"EDIType is not in the list [997]."},
			},
		},
		"Message Type Invalid": {
			ack: models.EdiErrorsAcknowledgementCodeError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    "EDI956",
				EdiErrorID:                 ediError.ID,
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: prICN.ID,
				Code:                       "C",
				Description:                "EDI Error happened to field 123",
			},
			expectedErrs: map[string][]string{
				"editype": {"EDIType is not in the list [997]."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.ack, test.expectedErrs)
		})
	}
}
