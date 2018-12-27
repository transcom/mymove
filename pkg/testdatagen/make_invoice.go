package testdatagen

import (
	"github.com/pkg/errors"
	"math/rand"
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

	invoiceNumber := assertions.Invoice.InvoiceNumber
	if invoiceNumber == "" {
		invoiceNumber = string(10000 + rand.Intn(89999))
	}

	// filled in dummy data
	Invoice := models.Invoice{
		Status:        models.InvoiceStatusINPROCESS,
		ApproverID:    approver.ID,
		Approver:      approver,
		InvoiceNumber: invoiceNumber,
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

// ResetInvoiceNumber resets the invoice number for a given SCAC/year.
func ResetInvoiceNumber(db *pop.Connection, scac string, year int) error {
	if len(scac) == 0 {
		return errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	sql := `DELETE FROM invoice_number_trackers WHERE standard_carrier_alpha_code = $1 AND year = $2`
	return db.RawQuery(sql, scac, year).Exec()
}
