package testdatagen

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

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
		scac := "ABCD"
		year := shipment.CreatedAt.UTC().Year()
		invoiceNumber = fmt.Sprintf("%s%d%04d", scac, year%100, 1+rand.Intn(9999))
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

// ResetInvoiceSequenceNumber resets the invoice sequence number for a given SCAC/year.
func ResetInvoiceSequenceNumber(db *pop.Connection, scac string, year int) error {
	if len(scac) == 0 {
		return errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	sql := `DELETE FROM invoice_number_trackers WHERE standard_carrier_alpha_code = $1 AND year = $2`
	return db.RawQuery(sql, scac, year).Exec()
}

// SetInvoiceSequenceNumber sets the invoice sequence number for a given SCAC/year.
func SetInvoiceSequenceNumber(db *pop.Connection, scac string, year int, sequenceNumber int) error {
	if len(scac) == 0 {
		return errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	sql := `INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
			VALUES ($1, $2, $3)
		ON CONFLICT (standard_carrier_alpha_code, year)
		DO
			UPDATE
				SET sequence_number = $3
				WHERE trackers.standard_carrier_alpha_code = $1 AND trackers.year = $2
	`

	return db.RawQuery(sql, scac, year, sequenceNumber).Exec()
}
