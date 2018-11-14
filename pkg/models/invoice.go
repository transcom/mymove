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

// Invoice is a collection of line item charges to be sent for payment
type Invoice struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Status        string    `json:"status" db:"status"`
	InvoiceNumber string    `json:"invoice_number" db:"invoice_number"`
	InvoicedDate  time.Time `json:"invoiced_date" db:"invoiced_date"`
	ShipmentID    uuid.UUID `json:"shipment_id" db:"shipment_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Invoice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Status, Name: "Status"},
		&validators.StringIsPresent{Field: i.InvoiceNumber, Name: "InvoiceNumber"},
		&validators.TimeIsPresent{Field: i.InvoicedDate, Name: "InvoicedDate"},
		&validators.UUIDIsPresent{Field: i.ShipmentID, Name: "ShipmentID"},
	), nil
}

// FetchInvoice Fetches and Validates an invoice model
func FetchInvoice(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	err := db.Eager().Find(&invoice, id)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// TODO: Check that the TSP user is authorized to get this Invoice
	// CHeck if office app, then...??
	// Can SMs fetch an invoice?
	// shipment, err := FetchShipmentForTSP(db, session, invoice.ShipmentID)
	// if err != nil {
	// 	return nil, err
	// }
	// if session.IsTspApp() && shipment.ID != invoice.ID{
	// 	return nil, ErrFetchForbidden
	// }

	return &invoice, nil
}
