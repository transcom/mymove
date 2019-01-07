package invoice

import (
	"errors"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
)

// UpdateInvoiceSubmitted is a service object to invoices into the Submitted state
type UpdateInvoiceSubmitted struct {
	DB *pop.Connection
}

// Call updates the Invoice to InvoiceStatusSUBMITTED and updates its ShipmentLineItem associations
func (u UpdateInvoiceSubmitted) Call(invoice *models.Invoice, shipmentLineItems models.ShipmentLineItems) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	var err error
	transactionErr := u.DB.Transaction(func(connection *pop.Connection) error {
		invoice.Status = models.InvoiceStatusSUBMITTED
		// Sample code of what eager creation should like
		// verrs, err := c.db.Eager().ValidateAndSave(&invoice)
		// Currently, this is only supported with `ValidateAndCreate`
		// We might want to consider adding this functionality to pop
		verrs, err = connection.ValidateAndSave(invoice)
		if err != nil || verrs.HasAny() {
			return errors.New("error saving invoice")
		}
		for liIndex := range shipmentLineItems {
			shipmentLineItems[liIndex].InvoiceID = &invoice.ID
			shipmentLineItems[liIndex].Invoice = *invoice
		}
		verrs, err = connection.ValidateAndSave(&shipmentLineItems)
		if err != nil || verrs.HasAny() {
			return errors.New("error saving shipment line items")
		}

		return nil
	})
	if transactionErr != nil {
		return verrs, err
	}
	return verrs, nil
}
