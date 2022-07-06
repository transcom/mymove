package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// WeightTicket represents the weight tickets and related data for a single trip of a PPM Shipment. Each trip should be
// its own record.
type WeightTicket struct {
	ID                                uuid.UUID   `json:"id" db:"id"`
	PPMShipmentID                     uuid.UUID   `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment                       PPMShipment `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	CreatedAt                         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt                         time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt                         *time.Time  `json:"deleted_at" db:"deleted_at"`
	VehicleDescription                *string     `json:"vehicle_description" db:"vehicle_description"`
	EmptyWeight                       *unit.Pound `json:"empty_weight" db:"empty_weight"`
	MissingEmptyWeightTicket          *bool       `json:"missing_empty_weight_ticket" db:"missing_empty_weight_ticket"`
	EmptyDocumentID                   uuid.UUID   `json:"empty_document_id" db:"empty_document_id"`
	EmptyDocument                     Document    `belongs_to:"documents" fk_id:"empty_document_id"`
	FullWeight                        *unit.Pound `json:"full_weight" db:"full_weight"`
	MissingFullWeightTicket           *bool       `json:"missing_full_weight_ticket" db:"missing_full_weight_ticket"`
	FullDocumentID                    uuid.UUID   `json:"full_document_id" db:"full_document_id"`
	FullDocument                      Document    `belongs_to:"documents" fk_id:"full_document_id"`
	OwnsTrailer                       *bool       `json:"owns_trailer" db:"owns_trailer"`
	TrailerMeetsCriteria              *bool       `json:"trailer_meets_criteria" db:"trailer_meets_criteria"`
	ProofOfTrailerOwnershipDocumentID uuid.UUID   `json:"proof_of_trailer_ownership_document_id" db:"proof_of_trailer_ownership_document_id"`
	ProofOfTrailerOwnershipDocument   Document    `belongs_to:"documents" fk_id:"proof_of_trailer_ownership_document_id"`
}

type WeightTickets []WeightTicket

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (w *WeightTicket) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: w.DeletedAt},
		&StringIsNilOrNotBlank{Name: "VehicleDescription", Field: w.VehicleDescription},
		&OptionalPoundIsNonNegative{Name: "EmptyWeight", Field: w.EmptyWeight},
		&validators.UUIDIsPresent{Name: "EmptyDocumentID", Field: w.EmptyDocumentID},
		&OptionalPoundIsNonNegative{Name: "FullWeight", Field: w.FullWeight},
		&validators.UUIDIsPresent{Name: "FullDocumentID", Field: w.FullDocumentID},
		&validators.UUIDIsPresent{Name: "ProofOfTrailerOwnershipDocumentID", Field: w.ProofOfTrailerOwnershipDocumentID},
	), nil
}
