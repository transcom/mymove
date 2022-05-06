package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
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
	// MovingExpenseTypeSTORAGE captures enum value "STORAGE"
	MovingExpenseTypeSTORAGE MovingExpenseType = "STORAGE"
	// MovingExpenseTypeOTHER captures enum value "OTHER"
	MovingExpenseTypeOTHER MovingExpenseType = "OTHER"
)

// IsExpenseModelDocumentType determines whether a MoveDocumentType is associated with a MovingExpenseDocument
func IsExpenseModelDocumentType(docType MoveDocumentType) bool {
	expenseModelDocumentTypes := []MoveDocumentType{
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
	ID                   uuid.UUID         `json:"id" db:"id"`
	MoveDocumentID       uuid.UUID         `json:"move_document_id" db:"move_document_id"`
	MoveDocument         MoveDocument      `belongs_to:"move_documents" fk_id:"move_document_id"`
	MovingExpenseType    MovingExpenseType `json:"moving_expense_type" db:"moving_expense_type"`
	RequestedAmountCents unit.Cents        `json:"requested_amount_cents" db:"requested_amount_cents"`
	PaymentMethod        string            `json:"payment_method" db:"payment_method"`
	ReceiptMissing       bool              `json:"receipt_missing" db:"receipt_missing"`
	StorageStartDate     *time.Time        `json:"storage_start_date" db:"storage_start_date"`
	StorageEndDate       *time.Time        `json:"storage_end_date" db:"storage_end_date"`
	CreatedAt            time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time        `db:"deleted_at"`
}

// MovingExpenseDocuments is not required by pop and may be deleted
type MovingExpenseDocuments []MovingExpenseDocument

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MovingExpenseDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.MoveDocumentID, Name: "MoveDocumentID"},
		&validators.StringIsPresent{Field: string(m.MovingExpenseType), Name: "MovingExpenseType"},
		&validators.StringIsPresent{Field: string(m.PaymentMethod), Name: "PaymentMethod"},
		&validators.IntIsGreaterThan{Field: int(m.RequestedAmountCents), Name: "RequestedAmountCents", Compared: 0},
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

// DaysInStorage calculates the days in storage excluding the entry day
func (m *MovingExpenseDocument) DaysInStorage() (int, error) {
	if m.MovingExpenseType != MovingExpenseTypeSTORAGE {
		return 0, fmt.Errorf("not storage expense")
	}
	if m.StorageStartDate == nil || m.StorageEndDate == nil {
		return 0, fmt.Errorf("storage end date or storage start date is nil")
	}
	if m.StorageEndDate.Before(*m.StorageStartDate) {
		return 0, fmt.Errorf("storage end date before storage start date")
	}
	// don't include the first day
	daysInStorage := int(m.StorageEndDate.Sub(*m.StorageStartDate).Hours() / 24)
	if daysInStorage < 0 {
		return 0, nil
	}
	return daysInStorage, nil
}

//FilterSITExpenses filter MovingExpenseDocuments to only storage expenses
func FilterSITExpenses(movingExpenseDocuments MovingExpenseDocuments) MovingExpenseDocuments {
	var sitExpenses []MovingExpenseDocument
	for _, doc := range movingExpenseDocuments {
		if doc.MovingExpenseType == MovingExpenseTypeSTORAGE {
			sitExpenses = append(sitExpenses, doc)
		}
	}
	return sitExpenses
}

//FilterMovingExpenseDocuments filter MoveDocuments to only moving expense documents
func FilterMovingExpenseDocuments(moveDocuments MoveDocuments) MovingExpenseDocuments {
	var movingExpenses []MovingExpenseDocument
	for _, moveDocument := range moveDocuments {
		if moveDocument.MovingExpenseDocument != nil {
			movingExpenses = append(movingExpenses, *moveDocument.MovingExpenseDocument)
		}
	}
	return movingExpenses
}
