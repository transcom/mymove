package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestInvoiceValidations() {
	invoice := &Invoice{}

	expErrors := map[string][]string{
		"status":         {"Status can not be blank."},
		"invoice_number": {"InvoiceNumber can not be blank."},
		"invoiced_date":  {"InvoicedDate can not be blank."},
		"shipment_id":    {"ShipmentID can not be blank."},
	}

	suite.verifyValidationErrors(invoice, expErrors)
}
