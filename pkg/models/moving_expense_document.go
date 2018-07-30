package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// MovingExpenseType represents types of different moving expenses
type MovingExpenseType string

const (
	// MovingExpenseTypeCONTRACTEDEXPENSE captures enum value "CONTRACTED_EXPENSE"
	MovingExpenseTypeCONTRACTEDEXPENSE MovingExpenseType = "CONTRACTED_EXPENSE"
	// MovingExpenseTypeRENTALEQUIPMENT captures enum value "RENTAL_EQUIPMENT"
	MovingExpenseTypeRENTALEQUIPMENT MovingExpenseType = "RENTAL_EQUIPMENT"
	// MovingExpenseTypePACKINGMATERIALS captures enum value "PACKING_MATERIALS"
	MovingExpenseTypePACKINGMATERIALS MovingExpenseType = "PACKING_MATERIALS"
	// MovingExpenseTypeWEIGHINGFEES captures enum value "WEIGHING_FEES"
	MovingExpenseTypeWEIGHINGFEES MovingExpenseType = "WEIGHING_FEES"
	// MovingExpenseTypeGAS captures enum value "GAS"
	MovingExpenseTypeGAS MovingExpenseType = "GAS"
	// MovingExpenseTypeTOLLS captures enum value "TOLLS"
	MovingExpenseTypeTOLLS MovingExpenseType = "TOLLS"
	// MovingExpenseTypeOIL captures enum value "OIL"
	MovingExpenseTypeOIL MovingExpenseType = "OIL"
	// MovingExpenseTypeOTHER captures enum value "OTHER"
	MovingExpenseTypeOTHER MovingExpenseType = "OTHER"
)

// IsExpenseModelDocumentType determines whether a MoveDocumentType is associated with a MovingExpenseDocument
func IsExpenseModelDocumentType(docType MoveDocumentType) bool {
	expenseModelDocumentTypes := []MoveDocumentType{
		MoveDocumentTypeSTORAGEEXPENSE,
		MoveDocumentTypeEXPENSE,
	}

	for _, t := range expenseModelDocumentTypes {
		if t == docType {
			return true
		}
	}

	return false
}

// MovingExpenseDocument is an object representing a move document
type MovingExpenseDocument struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	MoveDocumentID    uuid.UUID         `json:"move_document_id" db:"move_document_id"`
	MoveDocument      MoveDocument      `belongs_to:"move_documents"`
	MovingExpenseType MovingExpenseType `json:"moving_expense_type" db:"moving_expense_type"`
	ReimbursementID   uuid.UUID         `json:"reimbursement_id" db:"reimbursement_id"`
	Reimbursement     Reimbursement     `belongs_to:"reimbursement"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

// MovingExpenseDocuments is not required by pop and may be deleted
type MovingExpenseDocuments []MovingExpenseDocument

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MovingExpenseDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.MoveDocumentID, Name: "MoveDocumentID"},
		&validators.UUIDIsPresent{Field: m.ReimbursementID, Name: "ReimbursementID"},
		&validators.StringIsPresent{Field: string(m.MovingExpenseType), Name: "MovingExpenseType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MovingExpenseDocument) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MovingExpenseDocument) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
