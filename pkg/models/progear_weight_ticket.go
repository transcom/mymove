package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/transcom/mymove/pkg/unit"
)

type ProgearWeightTicket struct {
	ID                          uuid.UUID          `json:"id" db:"id"`
	PPMShipmentID               uuid.UUID          `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment                 PPMShipment        `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	BelongsToSelf               *bool              `json:"belongs_to_self" db:"belongs_to_self"`
	Description                 *string            `json:"description" db:"description"`
	HasWeightTickets            *bool              `json:"has_weight_tickets" db:"has_weight_tickets"`
	EmptyWeight                 *unit.Pound        `json:"empty_weight" db:"empty_weight"`
	EmptyDocumentID             uuid.UUID          `json:"empty_document_id" db:"empty_document_id"`
	EmptyDocument               Document           `belongs_to:"documents" fk_id:"empty_document_id"`
	FullWeight                  *unit.Pound        `json:"full_weight" db:"full_weight"`
	FullDocumentID              uuid.UUID          `json:"full_document_id" db:"full_document_id"`
	FullDocument                Document           `belongs_to:"documents" fk_id:"full_document_id"`
	ConstructedWeight           *unit.Pound        `json:"constructed_weight" db:"constructed_weight"`
	ConstructedWeightDocumentID uuid.UUID          `json:"constructed_weight_document_id" db:"constructed_weight_document_id"`
	ConstructedWeightDocument   Document           `belongs_to:"documents" fk_id:"constructed_weight_document_id"`
	Status                      *PPMDocumentStatus `json:"status" db:"status"`
	Reason                      *string            `json:"reason" db:"reason"`
	CreatedAt                   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt                   *time.Time         `json:"deleted_at" db:"deleted_at"`
}

type ProgearWeightTickets []ProgearWeightTicket

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (p *ProgearWeightTicket) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "PPMShipmentID", Field: p.PPMShipmentID},
		&StringIsNilOrNotBlank{Name: "Description", Field: p.Description},
		&validators.UUIDIsPresent{Name: "EmptyDocumentID", Field: p.EmptyDocumentID},
		&OptionalPoundIsNonNegative{Name: "EmptyWeight", Field: p.EmptyWeight},
		&validators.UUIDIsPresent{Name: "FullDocumentID", Field: p.FullDocumentID},
		&OptionalPoundIsNonNegative{Name: "FullWeight", Field: p.FullWeight},
		&validators.UUIDIsPresent{Name: "ConstructedWeightDocumentID", Field: p.ConstructedWeightDocumentID},
		&OptionalPoundIsNonNegative{Name: "ConstructedWeight", Field: p.ConstructedWeight},
		&OptionalStringInclusion{Name: "Status", Field: (*string)(p.Status), List: AllowedPPMDocumentStatuses},
		&StringIsNilOrNotBlank{Name: "Reason", Field: p.Reason},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: p.DeletedAt},
	), nil
}
