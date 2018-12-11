package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeInvoice creates a single invoice record with an associated shipment
func MakeInvoice(db *pop.Connection, assertions Assertions) models.Invoice {
	approver := assertions.Invoice.Approver
	if isZeroUUID(approver.ID) {
		approver = MakeOfficeUser(db, assertions)
	}

	shipment := assertions.Invoice.Shipment
	if isZeroUUID(shipment.ID) {
		shipment = MakeShipment(db, assertions)
	}

	// filled in dummy data
	Invoice := models.Invoice{
		Status:        models.InvoiceStatusINPROCESS,
		ApproverID:    approver.ID,
		Approver:      approver,
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
