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
	ID                   uuid.UUID  `db:"id"`
	Agency               string     `db:"agency"`
	CurrentAddress       Address    `belongs_to:"address"`
	CurrentAddressID     *uuid.UUID `db:"current_address_id"`
	DODID                string     `db:"dod_id"`
	DestinationAddress   Address    `belongs_to:"address"`
	DestinationAddressID *uuid.UUID `db:"current_address_id"`
	Email                *string    `db:"email"`
	FirstName            string     `db:"first_name"`
	LastName             string     `db:"last_name"`
	PhoneNumber          *string    `db:"phone_number"`
	User                 User       `belongs_to:"users"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (c *Customer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: c.UserID, Name: "UserID"})
	vs = append(vs, &validators.StringIsPresent{Field: c.DODID, Name: "DODID"})
	return validate.Validate(vs...), nil
}
