package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	// "github.com/transcom/mymove/pkg/models"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	// InvoiceStatusDRAFT captures enum value "DRAFT"
	InvoiceStatusDRAFT InvoiceStatus = "DRAFT"
	// InvoiceStatusINPROCESS captures enum value "IN_PROCESS"
	InvoiceStatusINPROCESS InvoiceStatus = "IN_PROCESS"
	// InvoiceStatusSUBMITTED captures enum value "SUBMITTED"
	InvoiceStatusSUBMITTED InvoiceStatus = "SUBMITTED"
	// InvoiceStatusSUBMISSIONFAILURE captures enum value "SUBMISSION_FAILURE"
	InvoiceStatusSUBMISSIONFAILURE InvoiceStatus = "SUBMISSION_FAILURE"
)

// Invoice is a collection of line item charges to be sent for payment
type Invoice struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	Status        InvoiceStatus `json:"status" db:"status"`
	InvoiceNumber string        `json:"invoice_number" db:"invoice_number"`
	InvoicedDate  time.Time     `json:"invoiced_date" db:"invoiced_date"`
	ShipmentID    uuid.UUID     `json:"shipment_id" db:"shipment_id"`
	Shipment      Shipment      `belongs_to:"shipments"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Invoice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(i.Status), Name: "Status"},
		&validators.StringIsPresent{Field: i.InvoiceNumber, Name: "InvoiceNumber"},
		&validators.TimeIsPresent{Field: i.InvoicedDate, Name: "InvoicedDate"},
		&validators.UUIDIsPresent{Field: i.ShipmentID, Name: "ShipmentID"},
	), nil
}

// FetchInvoice fetches and validates an invoice model
func FetchInvoice(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Invoice, error) {

	// Fetch invoice via invoice id
	var invoice Invoice
	err := db.Eager().Find(&invoice, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// Check that the TSP user is authorized to get this Invoice
	if session.IsTspUser() {
		_, _, err := FetchShipmentForVerifiedTSPUser(db, session.TspUserID, invoice.ShipmentID)
		if err != nil {
			return nil, err
		}
	} else if !session.IsOfficeUser() {
		// Allow office users to fetch invoices
		return nil, ErrFetchForbidden
	}

	return &invoice, nil
}

// GenerateInvoiceNumber creates an invoice number for a given SCAC/year.
func GenerateInvoiceNumber(db *pop.Connection, scac string, year int) (string, error) {
	if len(scac) == 0 {
		return "", errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return "", errors.Errorf("Year (%d) must be non-negative", year)
	}

	var sequenceNumber int
	sql := `INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
			VALUES ($1, $2, 1)
		ON CONFLICT (standard_carrier_alpha_code, year)
		DO
			UPDATE
				SET sequence_number = trackers.sequence_number + 1
				WHERE trackers.standard_carrier_alpha_code = $1 AND trackers.year = $2
		RETURNING sequence_number
	`

	err := db.RawQuery(sql, scac, year).First(&sequenceNumber)
	if err != nil {
		return "", errors.Wrapf(err, "Error when incrementing invoice sequence number for %s/%d", scac, year)
	}

	if sequenceNumber > 9999 {
		return "", errors.Errorf("All four-digit invoice sequence numbers already used for %s/%d", scac, year)
	}

	return fmt.Sprintf("%s%d%04d", scac, year%100, sequenceNumber), nil
}
