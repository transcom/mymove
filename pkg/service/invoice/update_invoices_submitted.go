package invoice

import (
	"errors"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
)

// UpdateInvoicesSubmitted is a service object to invoices into the Submitted state
type UpdateInvoicesSubmitted struct {
	DB *pop.Connection
}

// Call updates the Invoices to InvoiceStatusSUBMITTED and updates their ShipmentLineItem associations
func (u UpdateInvoicesSubmitted) Call(invoices models.Invoices, shipmentLineItems models.ShipmentLineItems) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	var err error
	transactionErr := u.DB.Transaction(func(connection *pop.Connection) error {
		for index := range invoices {
			invoices[index].Status = models.InvoiceStatusSUBMITTED
			// Sample code of what eager creation should like
			// Currently it is attempting to recreate shipment line items
			// and violating pk unique constraints (the docs say it shouldn't)
			// https://gobuffalo.io/en/docs/db/relations#eager-creation
			// verrs, err := c.db.Eager().ValidateAndSave(&invoice)
			verrs, err := u.DB.ValidateAndSave(&invoices[index])
			if err != nil || verrs.HasAny() {
				return errors.New("error saving invoice")
			}
			for liIndex := range shipmentLineItems {
				shipmentLineItems[liIndex].InvoiceID = &invoices[index].ID
				shipmentLineItems[liIndex].Invoice = invoices[index]
			}
			verrs, err = u.DB.ValidateAndSave(&shipmentLineItems)
			if err != nil || verrs.HasAny() {
				return errors.New("error saving shipment line items")
			}
		}
		return nil
	})
	if transactionErr != nil {
		return verrs, err
	}
	return verrs, nil
}
