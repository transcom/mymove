package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
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
	ID                uuid.UUID         `json:"id" db:"id"`
	ApproverID        uuid.UUID         `json:"approver_id" db:"approver_id"`
	Approver          OfficeUser        `belongs_to:"office_user"`
	Status            InvoiceStatus     `json:"status" db:"status"`
	InvoiceNumber     string            `json:"invoice_number" db:"invoice_number"`
	InvoicedDate      time.Time         `json:"invoiced_date" db:"invoiced_date"`
	ShipmentID        uuid.UUID         `json:"shipment_id" db:"shipment_id"`
	Shipment          Shipment          `belongs_to:"shipments"`
	ShipmentLineItems ShipmentLineItems `has_many:"shipment_line_items"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
	UploadID          *uuid.UUID        `json:"upload_id" db:"upload_id"`
	Upload            *Upload           `belongs_to:"uploads"`
}

// Invoices is an array of invoices
type Invoices []Invoice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Invoice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(i.Status), Name: "Status"},
		// Note that a SCAC can be 2 to 4 letters long, so the minimum invoice number
		// length should be 2 (SCAC) + 2 (two-digit year) + 4 (sequence number).
		&validators.StringLengthInRange{Field: i.InvoiceNumber, Name: "InvoiceNumber", Min: 8, Max: 255},
		&validators.TimeIsPresent{Field: i.InvoicedDate, Name: "InvoicedDate"},
		&validators.UUIDIsPresent{Field: i.ShipmentID, Name: "ShipmentID"},
		&validators.UUIDIsPresent{Field: i.ApproverID, Name: "ApproverID"},
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

// FetchInvoicesForShipment fetches all the invoices for a given shipment
func FetchInvoicesForShipment(db *pop.Connection, shipmentID uuid.UUID) (Invoices, error) {
	var invoices []Invoice
	err := db.Where("shipment_id = ?", shipmentID).Eager("Approver").All(&invoices)
	return invoices, err
}
