package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
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
	// MoveDocumentTypeWEIGHTTICKET captures enum value "WEIGHT_TICKET"
	MoveDocumentTypeWEIGHTTICKET MoveDocumentType = "WEIGHT_TICKET"
	// MoveDocumentTypeSTORAGEEXPENSE captures enum value "STORAGE_EXPENSE"
	MoveDocumentTypeSTORAGEEXPENSE MoveDocumentType = "STORAGE_EXPENSE"
	// MoveDocumentTypeSHIPMENTSUMMARY captures enum value "SHIPMENT_SUMMARY"
	MoveDocumentTypeSHIPMENTSUMMARY MoveDocumentType = "SHIPMENT_SUMMARY"
	// MoveDocumentTypeEXPENSE captures enum value "EXPENSE"
	MoveDocumentTypeEXPENSE MoveDocumentType = "EXPENSE"
)

// MoveDocumentSaveAction represents actions that can be taken during save
type MoveDocumentSaveAction string

const (
	// MoveDocumentSaveActionDELETEEXPENSEMODEL encodes an action to delete a linked expense model
	MoveDocumentSaveActionDELETEEXPENSEMODEL MoveDocumentSaveAction = "DELETE_EXPENSE_MODEL"
	// MoveDocumentSaveActionSAVEEXPENSEMODEL encodes an action to save a linked expense model
	MoveDocumentSaveActionSAVEEXPENSEMODEL MoveDocumentSaveAction = "SAVE_EXPENSE_MODEL"
)

// MoveDocument is an object representing a move document
type MoveDocument struct {
	ID                       uuid.UUID              `json:"id" db:"id"`
	DocumentID               uuid.UUID              `json:"document_id" db:"document_id"`
	Document                 Document               `belongs_to:"documents"`
	MoveID                   uuid.UUID              `json:"move_id" db:"move_id"`
	Move                     Move                   `belongs_to:"moves"`
	PersonallyProcuredMoveID *uuid.UUID             `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	Title                    string                 `json:"title" db:"title"`
	Status                   MoveDocumentStatus     `json:"status" db:"status"`
	MoveDocumentType         MoveDocumentType       `json:"move_document_type" db:"move_document_type"`
	MovingExpenseDocument    *MovingExpenseDocument `has_one:"moving_expense_document"`
	Notes                    *string                `json:"notes" db:"notes"`
	CreatedAt                time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time              `json:"updated_at" db:"updated_at"`
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
		&validators.StringIsPresent{Field: string(m.MoveDocumentType), Name: "MoveDocumentType"},
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

// FetchMoveDocument fetches a MoveDocument model
func FetchMoveDocument(db *pop.Connection, session *auth.Session, id uuid.UUID) (*MoveDocument, error) {
	var moveDoc MoveDocument
	err := db.Q().Eager("Document.Uploads").Find(&moveDoc, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}

	// Pointer associations are buggy, so we manually load expense document things
	movingExpenseDocument := MovingExpenseDocument{}
	moveDoc.MovingExpenseDocument = nil
	err = db.Where("move_document_id = $1", moveDoc.ID.String()).Eager("Reimbursement").First(&movingExpenseDocument)
	if err != nil {
		if errors.Cause(err).Error() != recordNotFoundErrorString {
			return nil, err
		}
	} else {
		moveDoc.MovingExpenseDocument = &movingExpenseDocument
	}

	// Check that the logged-in service member is associated to the document
	if session.IsMyApp() && moveDoc.Document.ServiceMemberID != session.ServiceMemberID {
		return &MoveDocument{}, ErrFetchForbidden
	}
	// Allow all office users to fetch move doc
	if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
		return &MoveDocument{}, ErrFetchForbidden
	}
	return &moveDoc, nil
}

// SaveMoveDocument saves a move document
func SaveMoveDocument(db *pop.Connection, moveDocument *MoveDocument, saveAction MoveDocumentSaveAction) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if saveAction == MoveDocumentSaveActionSAVEEXPENSEMODEL {
			// Save reimbursement first
			reimbursement := moveDocument.MovingExpenseDocument.Reimbursement
			if verrs, err := db.ValidateAndSave(&reimbursement); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving Moving Expense Reimbursement")
				return transactionError
			}
			// Make sure the moveDocument has the associated ID
			moveDocument.MovingExpenseDocument.ReimbursementID = reimbursement.ID
			moveDocument.MovingExpenseDocument.Reimbursement = reimbursement

			// Then save expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if verrs, err := db.ValidateAndSave(expenseDocument); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Creating Moving Expense Document")
				return transactionError
			}
		} else if saveAction == MoveDocumentSaveActionDELETEEXPENSEMODEL {
			// Destroy reimbursement first
			reimbursement := moveDocument.MovingExpenseDocument.Reimbursement
			if err := db.Destroy(&reimbursement); err != nil {
				responseError = errors.Wrap(err, "Error Deleting Moving Expense Reimbursement")
				return transactionError
			}

			// Then destroy expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if err := db.Destroy(expenseDocument); err != nil {
				responseError = errors.Wrap(err, "Error Deleting Moving Expense Document")
				return transactionError
			}
			moveDocument.MovingExpenseDocument = nil
		}

		// Finally, save the MoveDocument
		if verrs, err := db.ValidateAndSave(moveDocument); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move Document")
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError
}
