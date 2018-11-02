package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Invoice is a collection of line item charges to be sent for payment
type Invoice struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Status        string    `json:"status" db:"status"`
	InvoiceNumber string    `json:"invoice_number" db:"invoice_number"`
	InvoicedDate  time.Time `json:"invoiced_date" db:"invoiced_date"`
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
	), nil
}
