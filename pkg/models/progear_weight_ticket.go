package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

type ProgearWeightTicket struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	PPMShipmentID    uuid.UUID          `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment      PPMShipment        `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	BelongsToSelf    *bool              `json:"belongs_to_self" db:"belongs_to_self"`
	Description      *string            `json:"description" db:"description"`
	HasWeightTickets *bool              `json:"has_weight_tickets" db:"has_weight_tickets"`
	Weight           *unit.Pound        `json:"weight" db:"weight"`
	DocumentID       uuid.UUID          `json:"document_id" db:"document_id"`
	Document         Document           `belongs_to:"documents" fk_id:"document_id"`
	Status           *PPMDocumentStatus `json:"status" db:"status"`
	Reason           *string            `json:"reason" db:"reason"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time         `json:"deleted_at" db:"deleted_at"`
}

type ProgearWeightTickets []ProgearWeightTicket

func (e ProgearWeightTickets) FilterDeleted() ProgearWeightTickets {
	if len(e) == 0 {
		return e
	}

	nonDeletedTickets := ProgearWeightTickets{}
	for _, ticket := range e {
		if ticket.DeletedAt == nil {
			nonDeletedTickets = append(nonDeletedTickets, ticket)
		}
	}

	return nonDeletedTickets
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (p *ProgearWeightTicket) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "PPMShipmentID", Field: p.PPMShipmentID},
		&StringIsNilOrNotBlank{Name: "Description", Field: p.Description},
		&validators.UUIDIsPresent{Name: "DocumentID", Field: p.DocumentID},
		&OptionalPoundIsNonNegative{Name: "Weight", Field: p.Weight},
		&OptionalStringInclusion{Name: "Status", Field: (*string)(p.Status), List: AllowedPPMDocumentStatuses},
		&StringIsNilOrNotBlank{Name: "Reason", Field: p.Reason},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: p.DeletedAt},
	), nil
}
