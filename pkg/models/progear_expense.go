package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

type ProgearExpense struct {
	ID              uuid.UUID          `json:"id" db:"id"`
	PPMShipmentID   uuid.UUID          `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment     PPMShipment        `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	IsOwn           *bool              `json:"is_own" db:"is_own"`
	Description     *string            `json:"description" db:"description"`
	HasWeightTicket *bool              `json:"has_weight_ticket" db:"has_weight_ticket"`
	EmptyWeight     *unit.Pound        `json:"empty_weight" db:"empty_weight"`
	EmptyDocumentID uuid.UUID          `json:"empty_document_id" db:"empty_document_id"`
	EmptyDocument   Document           `belongs_to:"documents" fk_id:"empty_document_id"`
	FullWeight      *unit.Pound        `json:"full_weight" db:"full_weight"`
	FullDocumentID  uuid.UUID          `json:"full_document_id" db:"full_document_id"`
	FullDocument    Document           `belongs_to:"documents" fk_id:"full_document_id"`
	Status          *PPMDocumentStatus `json:"status" db:"status"`
	Reason          *string            `json:"reason" db:"reason"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time         `json:"deleted_at" db:"deleted_at"`
}

type ProgearExpenses []ProgearExpense
