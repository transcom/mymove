package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeInvoice creates a single invoice record with an associated shipment
func MakeInvoice(db *pop.Connection, assertions Assertions) models.Invoice {
	shipment := assertions.Invoice.Shipment
	if isZeroUUID(shipment.ID) {
		shipment = MakeShipment(db, assertions)
	}

	// filled in dummy data
	Invoice := models.Invoice{
		Status:        models.InvoiceStatusINPROCESS,
		InvoiceNumber: "1234",
		InvoicedDate:  time.Now(),
		ShipmentID:    shipment.ID,
		Shipment:      shipment,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&Invoice, assertions.Invoice)

	mustCreate(db, &Invoice)

	return Invoice
}

// MakeDefaultInvoice makes a invoice with default values
func MakeDefaultInvoice(db *pop.Connection) models.Invoice {
	return MakeInvoice(db, Assertions{})
}
