package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Organization represents an organization and their contact information
type Organization struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	PocEmail  *string   `json:"poc_email" db:"poc_email"`
	PocPhone  *string   `json:"poc_phone" db:"poc_phone"`
}

// TableName overrides the table name used by Pop.
func (o Organization) TableName() string {
	return "organizations"
}

func (o Organization) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

type Organizations []Organization

func (o Organizations) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *Organization) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Name, Name: "Name"},
	), nil
}
