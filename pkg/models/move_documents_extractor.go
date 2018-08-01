package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

// MoveDocumentExtractor is an object representing ANY move document, and thus has all the fields
type MoveDocumentExtractor struct {
	ID                       uuid.UUID          `json:"id" db:"id"`
	DocumentID               uuid.UUID          `json:"document_id" db:"document_id"`
	Document                 Document           `belongs_to:"documents"`
	MoveID                   uuid.UUID          `json:"move_id" db:"move_id"`
	Move                     Move               `belongs_to:"moves"`
	Title                    string             `json:"title" db:"title"`
	Status                   MoveDocumentStatus `json:"status" db:"status"`
	PersonallyProcuredMoveID *uuid.UUID         `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	MoveDocumentType         MoveDocumentType   `json:"move_document_type" db:"move_document_type"`
	MovingExpenseType        *MovingExpenseType `json:"moving_expense_type" db:"moving_expense_type"`
	ReimbursementID          *uuid.UUID         `json:"reimbursement_id" db:"reimbursement_id"`
	Reimbursement            Reimbursement      `belongs_to:"reimbursement"`
	Notes                    *string            `json:"notes" db:"notes"`
	CreatedAt                time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at" db:"updated_at"`
}

// MoveDocumentExtractors is not required by pop and may be deleted
type MoveDocumentExtractors []MoveDocumentExtractor

// FetchAllMoveDocumentsForMove fetches all MoveDocument models
func (m *Move) FetchAllMoveDocumentsForMove(db *pop.Connection) (MoveDocumentExtractors, error) {
	var moveDocs MoveDocumentExtractors
	query := db.Q().LeftJoin("moving_expense_documents ed", "ed.move_document_id=move_documents.id").
		Where("move_documents.move_id=$1", m.ID.String())

	sql, args := query.ToSQL(&pop.Model{Value: MoveDocument{}},
		"move_documents.*, ed.moving_expense_type, ed.reimbursement_id")

	err := db.RawQuery(sql, args...).Eager("Document.Uploads", "Reimbursement").All(&moveDocs)
	if err != nil {
		return moveDocs, err
	}

	return moveDocs, nil
}
