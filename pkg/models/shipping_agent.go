package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// AgentType represents the type of agent being recorded
type AgentType string

const (
	// AgentTypeORIGIN capture enum value "ORIGIN"
	AgentTypeORIGIN AgentType = "ORIGIN"
	// AgentTypeDESTINATION capture enum value "DESTINATION"
	AgentTypeDESTINATION AgentType = "DESTINATION"
)

// ShippingAgent represents an assigned agent for a shipment
type ShippingAgent struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ShipmentID  uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment    *Shipment `belongs_to:"shipment"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	AgentType   AgentType `json:"agent_type" db:"agent_type"`
	Name        string    `json:"name" db:"name"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
}

// String is not required by pop and may be deleted
func (s ShippingAgent) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ShippingAgents is not required by pop and may be deleted
type ShippingAgents []ShippingAgent

// String is not required by pop and may be deleted
func (s ShippingAgents) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ShippingAgent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.ShipmentID, Name: "ShipmentID"},
		&validators.StringIsPresent{Field: string(s.AgentType), Name: "AgentType"},
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		&validators.StringIsPresent{Field: s.PhoneNumber, Name: "PhoneNumber"},
		&validators.StringIsPresent{Field: s.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ShippingAgent) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ShippingAgent) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
