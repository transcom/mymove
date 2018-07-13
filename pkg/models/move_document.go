package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// MoveDocumentStatus represents the status of a move document record's lifecycle
type MoveDocumentStatus string

const (
	// MoveDocumentStatusAWAITINGREVIEW captures enum value "AWAITING_REVIEW"
	MoveDocumentStatusAWAITINGREVIEW MoveDocumentStatus = "AWAITING_REVIEW"
	// MoveDocumentStatusOK captures enum value "OK"
	MoveDocumentStatusOK MoveDocumentStatus = "OK"
	// MoveDocumentStatusHASISSUE captures enum value "HAS_ISSUE"
	MoveDocumentStatusHASISSUE MoveDocumentStatus = "HAS_ISSUE"
)

// MoveDocumentType represents types of different move documents
type MoveDocumentType string

const (
	// MoveDocumentTypeOTHER captures enum value "OTHER"
	MoveDocumentTypeOTHER MoveDocumentType = "OTHER"
)

// MoveDocument is an object representing a move document
type MoveDocument struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	DocumentID       uuid.UUID          `json:"document_id" db:"document_id"`
	Document         Document           `belongs_to:"documents"`
	MoveID           uuid.UUID          `json:"move_id" db:"move_id"`
	Move             Move               `belongs_to:"moves"`
	Title            string             `json:"title" db:"title"`
	Status           MoveDocumentStatus `json:"status" db:"status"`
	MoveDocumentType MoveDocumentType   `json:"move_document_type" db:"move_document_type"`
	Notes            *string            `json:"notes" db:"notes"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
}

// MoveDocuments is not required by pop and may be deleted
type MoveDocuments []MoveDocument

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MoveDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.DocumentID, Name: "DocumentID"},
		&validators.UUIDIsPresent{Field: m.MoveID, Name: "MoveID"},
		&validators.StringIsPresent{Field: string(m.Title), Name: "Title"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "MoveDocumentType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MoveDocument) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MoveDocument) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
