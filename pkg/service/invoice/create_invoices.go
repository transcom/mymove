package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// CreateInvoices is a service object to create new invoices from Shipments
type CreateInvoices struct {
	DB        *pop.Connection
	Shipments []models.Shipment
}

// Call creates Invoices and updates their ShipmentLineItem associations
func (c CreateInvoices) Call(clock clock.Clock) ([]models.Invoice, error) {
	currentTime := clock.Now()
	invoices := make(models.Invoices, 0)
	err := c.db.Transaction(func(connection *pop.Connection) error {
		for _, shipment := range c.shipments {
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
			c.db.ValidateAndCreate(&invoice)
			invoices = append(invoices, invoice)
			for index := range shipment.ShipmentLineItems {
				shipment.ShipmentLineItems[index].InvoiceID = &invoice.ID
				shipment.ShipmentLineItems[index].Invoice = invoice
			}
			verrs, err := c.db.ValidateAndSave(&shipment.ShipmentLineItems)
			if err != nil || verrs.HasAny() {
				return err
			}
		}
		return nil
	})
	return invoices, err
}
