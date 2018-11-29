package invoice

import (
	"errors"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
)

// CreateInvoices is a service object to create new invoices from Shipments
type CreateInvoices struct {
	DB        *pop.Connection
	Shipments []models.Shipment
}

// Call creates Invoices and updates their ShipmentLineItem associations
func (c CreateInvoices) Call(clock clock.Clock) (*validate.Errors, error) {
	currentTime := clock.Now()
	var invoices models.Invoices
	verrs := validate.NewErrors()
	var err error
	transactionErr := c.DB.Transaction(func(connection *pop.Connection) error {
		for _, shipment := range c.Shipments {
			invoice := models.Invoice{
				Status:            models.InvoiceStatusINPROCESS,
				InvoiceNumber:     "1", // placeholder
				InvoicedDate:      currentTime,
				ShipmentID:        shipment.ID,
				Shipment:          shipment,
				ShipmentLineItems: shipment.ShipmentLineItems,
			}
			// Sample code of what eager creation should like
			// Currently it is attempting to recreate shipment line items
			// and violating pk unique constraints (the docs say it shouldn't)
			// https://gobuffalo.io/en/docs/db/relations#eager-creation
			// verrs, err := c.db.Eager().ValidateAndCreate(&invoices)
			verrs, err := c.DB.ValidateAndCreate(&invoice)
			if err != nil || verrs.HasAny() {
				return errors.New("error saving invoice")
			}
			invoices = append(invoices, invoice)
			for index := range shipment.ShipmentLineItems {
				shipment.ShipmentLineItems[index].InvoiceID = &invoice.ID
				shipment.ShipmentLineItems[index].Invoice = invoice
			}
			verrs, err = c.DB.ValidateAndSave(&shipment.ShipmentLineItems)
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
