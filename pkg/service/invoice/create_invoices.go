package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
)

// CreateInvoices is a service object to create new invoices from Shipments
type CreateInvoices struct {
	DB    *pop.Connection
	Clock clock.Clock
}

// Call creates Invoices and updates their ShipmentLineItem associations
func (c CreateInvoices) Call(invoices *models.Invoices, shipments models.Shipments) (*validate.Errors, error) {
	for _, shipment := range shipments {
		*invoices = append(*invoices, models.Invoice{
			Status:        models.InvoiceStatusINPROCESS,
			InvoiceNumber: "1", // placeholder
			InvoicedDate:  c.Clock.Now(),
			ShipmentID:    shipment.ID,
			Shipment:      shipment,
		})
	}
	return c.DB.ValidateAndCreate(invoices)
}
