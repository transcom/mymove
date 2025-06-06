package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestEdiErrors() {
	icnID := uuid.Must(uuid.NewV4())
	testCases := map[string]struct {
		ediError     models.EdiError
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDIType824,
				PaymentRequestID:           uuid.Must(uuid.NewV4()),
				InterchangeControlNumberID: &icnID,
				Code:                       models.StringPointer("B"),
				Description:                models.StringPointer("EDI Error happened to field 99"),
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			ediError: models.EdiError{},
			expectedErrs: map[string][]string{
				"description":        {"Both Description and Code cannot be nil, one must be valid"},
				"code":               {"Both Code and Description cannot be nil, one must be valid"},
				"payment_request_id": {"PaymentRequestID can not be blank."},
				"editype":            {"EDIType is not in the list [810, 824, 858, 997, TPPSPaidInvoiceReport]."},
			},
		},
		"Message Type Invalid": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    "956",
				PaymentRequestID:           uuid.Must(uuid.NewV4()),
				InterchangeControlNumberID: &icnID,
				Code:                       models.StringPointer("C"),
				Description:                models.StringPointer("EDI Error happened to field 123"),
			},
			expectedErrs: map[string][]string{
				"editype": {"EDIType is not in the list [810, 824, 858, 997, TPPSPaidInvoiceReport]."},
			},
		},
		"At least one valid Code or Description": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDIType824,
				PaymentRequestID:           uuid.Must(uuid.NewV4()),
				InterchangeControlNumberID: &icnID,
				Description:                models.StringPointer("EDI Error happened to field 99"),
			},
			expectedErrs: nil,
		},
		"At least one valid Code or Description and no empty string": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDIType824,
				PaymentRequestID:           uuid.Must(uuid.NewV4()),
				InterchangeControlNumberID: &icnID,
				Description:                models.StringPointer(""),
			},
			expectedErrs: map[string][]string{
				"description": {"Description string if present should not be empty"},
			},
		},
	}

	for name, test := range testCases {
		suite.Run(name, func() {
			// nolint:gosec //G402
			suite.verifyValidationErrors(&test.ediError, test.expectedErrs, nil)

		})
	}
}
