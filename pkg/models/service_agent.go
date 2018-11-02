package models

import (
	"time"

	"encoding/json"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	Email            *string   `json:"email" db:"email"`
	PhoneNumber      *string   `json:"phone_number" db:"phone_number"`
	FaxNumber        *string   `json:"fax_number" db:"fax_number"`
	EmailIsPreferred *bool     `json:"email_is_preferred" db:"email_is_preferred"`
	PhoneIsPreferred *bool     `json:"phone_is_preferred" db:"phone_is_preferred"`
	Notes            *string   `json:"notes" db:"notes"`
	Company          string    `json:"company" db:"company"`
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
		&validators.StringIsPresent{Field: s.Company, Name: "Company"},
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

// FetchServiceAgentsByTSP looks up all service agents beloning to a TSP and a shipment
func FetchServiceAgentsByTSP(tx *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID) ([]ServiceAgent, error) {

	serviceAgents := []ServiceAgent{}

	err := tx.
		Where("shipments.id = $1 AND shipment_offers.transportation_service_provider_id = $2", shipmentID, tspID).
		LeftJoin("shipments", "service_agents.shipment_id=shipments.id").
		LeftJoin("shipment_offers", "shipments.id=shipment_offers.shipment_id").
		All(&serviceAgents)
	if err != nil {
		return nil, err
	}

	return serviceAgents, err
}

// FetchServiceAgentByTSP looks up all service agents beloning to a TSP and a shipment
func FetchServiceAgentByTSP(tx *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID, serviceAgentID uuid.UUID) (*ServiceAgent, error) {

	serviceAgents := []ServiceAgent{}

	err := tx.
		Where("service_agents.id = $1 AND shipments.id = $2 AND shipment_offers.transportation_service_provider_id = $3", serviceAgentID, shipmentID, tspID).
		LeftJoin("shipments", "service_agents.shipment_id=shipments.id").
		LeftJoin("shipment_offers", "shipments.id=shipment_offers.shipment_id").
		All(&serviceAgents)

	if err != nil {
		return nil, err
	}

	// Unlikely that we see more than one but to be safe this will error.
	if len(serviceAgents) != 1 {
		return nil, ErrFetchNotFound
	}

	return &serviceAgents[0], err
}

// CreateServiceAgent creates a ServiceAgent model from payload and queried fields.
func CreateServiceAgent(tx *pop.Connection,
	shipmentID uuid.UUID,
	role Role,
	company *string,
	email *string,
	phoneNumber *string,
	emailIsPreferred *bool,
	phoneIsPreferred *bool,
	notes *string) (ServiceAgent, *validate.Errors, error) {

	var stringCompany string
	if company != nil {
		stringCompany = string(*company)
	}
	newServiceAgent := ServiceAgent{
		ShipmentID:       shipmentID,
		Role:             role,
		Email:            email,
		PhoneNumber:      phoneNumber,
		EmailIsPreferred: emailIsPreferred,
		PhoneIsPreferred: phoneIsPreferred,
		Notes:            notes,
		Company:          stringCompany,
	}
	verrs, err := tx.ValidateAndCreate(&newServiceAgent)
	if err != nil {
		zap.L().Error("DB insertion error", zap.Error(err))
		return ServiceAgent{}, verrs, err
	} else if verrs.HasAny() {
		zap.L().Error("Validation errors", zap.Error(verrs))
		return ServiceAgent{}, verrs, errors.New("Validation error on Service Agent")
	}
	return newServiceAgent, verrs, err
}
