package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// Role represents the type of agent being recorded
type Role string

const (
	// RoleORIGIN capture enum value "ORIGIN"
	RoleORIGIN Role = "ORIGIN"
	// RoleDESTINATION capture enum value "DESTINATION"
	RoleDESTINATION Role = "DESTINATION"
)

// ServiceAgent represents an assigned agent for a shipment
type ServiceAgent struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ShipmentID       uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment         *Shipment `belongs_to:"shipment"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Role             Role      `json:"role" db:"role"`
	PointOfContact   string    `json:"point_of_contact" db:"point_of_contact"`
	Email            *string   `json:"email" db:"email"`
	PhoneNumber      *string   `json:"phone_number" db:"phone_number"`
	FaxNumber        *string   `json:"fax_number" db:"fax_number"`
	EmailIsPreferred *bool     `json:"email_is_preferred" db:"email_is_preferred"`
	PhoneIsPreferred *bool     `json:"phone_is_preferred" db:"phone_is_preferred"`
	Notes            *string   `json:"notes" db:"notes"`
}

// String is not required by pop and may be deleted
func (s ServiceAgent) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ServiceAgents is not required by pop and may be deleted
type ServiceAgents []ServiceAgent

// String is not required by pop and may be deleted
func (s ServiceAgents) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceAgent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.ShipmentID, Name: "ShipmentID"},
		&validators.StringIsPresent{Field: string(s.Role), Name: "Role"},
		&validators.StringIsPresent{Field: s.PointOfContact, Name: "PointOfContact"},
		&validators.StringIsPresent{Field: *s.Email, Name: "Email"},
		&validators.StringIsPresent{Field: *s.PhoneNumber, Name: "PhoneNumber"},
		&validators.StringIsPresent{Field: *s.FaxNumber, Name: "FaxNumber"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceAgent) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceAgent) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
