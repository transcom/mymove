package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) TestEdiErrorsSendToSyncadaError() {
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
		sendError    models.EdiErrorsSendToSyncadaError
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			sendError: models.EdiErrorsSendToSyncadaError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDI858,
				EdiErrorID:                 ediError.ID,
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: &prICN.ID,
				Description:                "EDI Error happened to field 99",
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			sendError: models.EdiErrorsSendToSyncadaError{},
			expectedErrs: map[string][]string{
				"edi_error_id":       {"EdiErrorID can not be blank."},
				"description":        {"Description cannot be empty"},
				"payment_request_id": {"PaymentRequestID can not be blank."},
				"editype":            {"EDIType is not in the list [858]."},
			},
		},
		"Message Type Invalid": {
			sendError: models.EdiErrorsSendToSyncadaError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    "EDI956",
				EdiErrorID:                 ediError.ID,
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: &prICN.ID,
				Description:                "EDI Error happened to field 123",
			},
			expectedErrs: map[string][]string{
				"editype": {"EDIType is not in the list [858]."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.sendError, test.expectedErrs)
		})
	}
}
