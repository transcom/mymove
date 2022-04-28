package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// MoveDocumentExtractor is an object representing ANY move document, and thus has all the fields
type MoveDocumentExtractor struct {
	ID                       uuid.UUID            `json:"id" db:"id"`
	DocumentID               uuid.UUID            `json:"document_id" db:"document_id"`
	Document                 Document             `belongs_to:"documents" fk_id:"document_id"`
	MoveID                   uuid.UUID            `json:"move_id" db:"move_id"`
	Move                     Move                 `belongs_to:"moves" fk_id:"move_id"`
	Title                    string               `json:"title" db:"title"`
	Status                   MoveDocumentStatus   `json:"status" db:"status"`
	PersonallyProcuredMoveID *uuid.UUID           `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	MoveDocumentType         MoveDocumentType     `json:"move_document_type" db:"move_document_type"`
	MovingExpenseType        *MovingExpenseType   `json:"moving_expense_type" db:"moving_expense_type"`
	RequestedAmountCents     *unit.Cents          `json:"requested_amount_cents" db:"requested_amount_cents"`
	ReceiptMissing           *bool                `json:"receipt_missing" db:"receipt_missing"`
	EmptyWeight              *unit.Pound          `json:"empty_weight,omitempty" db:"empty_weight"`
	EmptyWeightTicketMissing *bool                `json:"empty_weight_ticket_missing,omitempty" db:"empty_weight_ticket_missing"`
	FullWeight               *unit.Pound          `json:"full_weight,omitempty" db:"full_weight"`
	FullWeightTicketMissing  *bool                `json:"full_weight_ticket_missing,omitempty" db:"full_weight_ticket_missing"`
	VehicleNickname          *string              `json:"vehicle_nickname,omitempty" db:"vehicle_nickname"`
	VehicleMake              *string              `json:"vehicle_make,omitempty" db:"vehicle_make"`
	VehicleModel             *string              `json:"vehicle_model,omitempty" db:"vehicle_model"`
	WeightTicketSetType      *WeightTicketSetType `json:"weight_ticket_set_type,omitempty" db:"weight_ticket_set_type"`
	WeightTicketDate         *time.Time           `json:"weight_ticket_date,omitempty" db:"weight_ticket_date"`
	TrailerOwnershipMissing  *bool                `json:"trailer_ownership_missing,omitempty" db:"trailer_ownership_missing"`
	PaymentMethod            *string              `json:"payment_method" db:"payment_method"`
	Notes                    *string              `json:"notes" db:"notes"`
	CreatedAt                time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time           `json:"deleted_at" db:"deleted_at"`
	StorageStartDate         *time.Time           `json:"storage_start_date" db:"storage_start_date"`
	StorageEndDate           *time.Time           `json:"storage_end_date" db:"storage_end_date"`
}

// MoveDocumentExtractors is not required by pop and may be deleted
type MoveDocumentExtractors []MoveDocumentExtractor

// FetchAllMoveDocumentsForMove fetches all MoveDocument models
func (m *Move) FetchAllMoveDocumentsForMove(db *pop.Connection, includeAllMoveDocuments bool) (MoveDocumentExtractors, error) {
	var moveDocs MoveDocumentExtractors
	query := db.Q().LeftJoin("moving_expense_documents ed", "ed.move_document_id=move_documents.id").
		LeftJoin("weight_ticket_set_documents wt", "wt.move_document_id=move_documents.id").
		Where("move_documents.move_id=$1", m.ID.String())

	if !includeAllMoveDocuments {
		query = query.Where("move_documents.deleted_at is null")
	}

	sql, args := query.ToSQL(&pop.Model{Value: MoveDocument{}},
		`move_documents.id,
	  move_documents.move_id,
	  move_documents.document_id,
	  move_documents.move_document_type,
	  move_documents.status,
	  move_documents.notes,
	  move_documents.created_at,
	  move_documents.updated_at,
	  move_documents.title,
	  move_documents.personally_procured_move_id,
	  move_documents.deleted_at,
	  ed.moving_expense_type,
	  ed.requested_amount_cents,
	  ed.payment_method,
      ed.receipt_missing,
      ed.storage_start_date,
      ed.storage_end_date,
	  wt.empty_weight,
	  wt.empty_weight_ticket_missing,
	  wt.full_weight_ticket_missing,
	  wt.full_weight,
	  wt.vehicle_nickname,
	  wt.vehicle_make,
	  wt.vehicle_model,
	  wt.weight_ticket_set_type,
	  wt.weight_ticket_date,
	  wt.trailer_ownership_missing`)

	err := db.RawQuery(sql, args...).Eager("Document.UserUploads.Upload").All(&moveDocs)
	if err != nil {
		return moveDocs, err
	}
	return moveDocs, nil
}
