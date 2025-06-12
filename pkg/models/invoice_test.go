package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestInvoiceValidation() {
	testCases := map[string]struct {
		invoice      models.Invoice
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			invoice: models.Invoice{
				ApproverID:    uuid.Must(uuid.NewV4()),
				Status:        "APPROVED",
				InvoiceNumber: "12345678",
				InvoicedDate:  time.Now(),
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			invoice: models.Invoice{},
			expectedErrs: map[string][]string{
				"invoiced_date":  {"InvoicedDate can not be blank."},
				"approver_id":    {"ApproverID can not be blank."},
				"status":         {"Status can not be blank."},
				"invoice_number": {"InvoiceNumber not in range(8, 255)"},
			},
		},
		"Other Errors": {
			invoice: models.Invoice{
				ApproverID:    uuid.Must(uuid.NewV4()),
				Status:        "APPROVED",
				InvoiceNumber: "1234567",
				InvoicedDate:  time.Now(),
			},
			expectedErrs: map[string][]string{
				"invoice_number": {"InvoiceNumber not in range(8, 255)"},
			},
		},
	}

	for name, test := range testCases {
		suite.Run(name, func() {
			//nolint:gosec // G601
			suite.verifyValidationErrors(&test.invoice, test.expectedErrs, nil)
		})
	}

}
