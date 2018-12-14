package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/models"
)

// CreateInvoice is a service object to create new invoices from Shipment
type CreateInvoice struct {
	DB    *pop.Connection
	Clock clock.Clock
}

// Call creates Invoice and updates its ShipmentLineItem associations
func (c CreateInvoice) Call(invoice *models.Invoice, shipment models.Shipment) (*validate.Errors, error) {
	*invoice = models.Invoice{
		Status:        models.InvoiceStatusINPROCESS,
		InvoiceNumber: "1", // placeholder
		InvoicedDate:  c.Clock.Now().UTC(),
		ShipmentID:    shipment.ID,
		Shipment:      shipment,
	}
	return c.DB.ValidateAndCreate(invoice)
}
