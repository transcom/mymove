package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Customer is an object representing data for a customer
type Customer struct {
	ID        uuid.UUID `db:"id"`
	DODID     string    `db:"dod_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	User      User      `belongs_to:"users"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (c *Customer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: c.UserID, Name: "UserID"})
	vs = append(vs, &validators.StringIsPresent{Field: c.DODID, Name: "DODID"})
	return validate.Validate(vs...), nil
}
