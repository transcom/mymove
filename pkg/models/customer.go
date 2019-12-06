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
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	DodID     *string   `db:"dod_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (c *Customer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringIsPresent{Field: c.FirstName, Name: "FirstName"})
	vs = append(vs, &validators.StringIsPresent{Field: c.LastName, Name: "LastName"})
	vs = append(vs, &validators.StringIsPresent{Field: c.Phone, Name: "Phone"})
	vs = append(vs, &validators.EmailIsPresent{Field: c.Email, Name: "Email"})
	return validate.Validate(vs...), nil
}
