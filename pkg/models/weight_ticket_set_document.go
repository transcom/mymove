package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// WeightTicketDocumentsPayload weight ticket documents payload
type WeightTicketSetDocument struct {
	ID                       uuid.UUID    `json:"id" db:"id"`
	MoveDocumentID           uuid.UUID    `json:"move_document_id" db:"move_document_id"`
	MoveDocument             MoveDocument `belongs_to:"move_documents"`
	EmptyWeight              *unit.Pound  `json:"empty_weight,omitempty" db:"empty_weight"`
	EmptyWeightTicketMissing bool         `json:"empty_weight_ticket_missing,omitempty" db:"empty_weight_ticket_missing"`
	FullWeight               *unit.Pound  `json:"full_weight,omitempty" db:"full_weight"`
	FullWeightTicketMissing  bool         `json:"full_weight_ticket_missing,omitempty" db:"full_weight_ticket_missing"`
	VehicleNickname          string       `json:"vehicle_nickname,omitempty" db:"vehicle_nickname"`
	VehicleOptions           string       `json:"vehicle_options,omitempty" db:"vehicle_options"`
	WeightTicketDate         *time.Time   `json:"weight_ticket_date,omitempty" db:"weight_ticket_date"`
	TrailerOwnershipMissing  bool         `json:"trailer_ownership_missing,omitempty" db:"trailer_ownership_missing"`
	CreatedAt                time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time    `json:"updated_at" db:"updated_at"`
}

// WeightTicketSetDocuments slice of WeightTicketSetDocuments
type WeightTicketSetDocuments []WeightTicketSetDocuments

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.MoveDocumentID, Name: "MoveDocumentID"},
		&validators.StringIsPresent{Field: string(m.VehicleNickname), Name: "VehicleNickname"},
		&validators.StringIsPresent{Field: string(m.VehicleOptions), Name: "VehicleOptions"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
