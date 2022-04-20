package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// MTOAgentType represents the type label for move task order agent
type MTOAgentType string

//Constants for the MTOAgentType
const (
	MTOAgentReleasing MTOAgentType = "RELEASING_AGENT"
	MTOAgentReceiving MTOAgentType = "RECEIVING_AGENT"
)

// MTOAgent is a struct that represents the mto_agents table.
type MTOAgent struct {
	ID            uuid.UUID    `db:"id"`
	MTOShipment   MTOShipment  `belongs_to:"mto_shipments" fk_id:"mto_shipment_id"`
	MTOShipmentID uuid.UUID    `db:"mto_shipment_id"`
	FirstName     *string      `db:"first_name"`
	LastName      *string      `db:"last_name"`
	Email         *string      `db:"email"`
	Phone         *string      `db:"phone"`
	MTOAgentType  MTOAgentType `db:"agent_type"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     *time.Time   `db:"deleted_at"`
}

//MTOAgents is a collection of MTOAgent
type MTOAgents []MTOAgent

//Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOAgent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MTOShipmentID, Name: "MTOShipmentID"})
	vs = append(vs, &validators.StringInclusion{Field: string(m.MTOAgentType), Name: "MTOAgentType", List: []string{
		string(MTOAgentReceiving),
		string(MTOAgentReleasing),
	}})
	if m.FirstName != nil {
		firstName := *m.FirstName
		vs = append(vs, &validators.StringIsPresent{Field: firstName, Name: "FirstName"})
	}
	if m.LastName != nil {
		lastName := *m.LastName
		vs = append(vs, &validators.StringIsPresent{Field: lastName, Name: "LastName"})
	}
	if m.Email != nil {
		email := *m.Email
		vs = append(vs, &validators.StringIsPresent{Field: email, Name: "Email"})
	}
	if m.Phone != nil {
		phone := *m.Phone
		vs = append(vs, &validators.StringIsPresent{Field: phone, Name: "Phone"})
	}
	return validate.Validate(vs...), nil
}

//TableName overrides the table name used by Pop.
func (m MTOAgent) TableName() string {
	return "mto_agents"
}
