package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// MoveTaskOrder is an object representing the task orders for a move
type MoveTaskOrder struct {
	ID                 uuid.UUID       `db:"id"`
	MoveOrder          MoveOrder       `belongs_to:"move_orders"`
	MTOServiceItems    MTOServiceItems `has_many:"mto_service_items"`
	PaymentRequests    PaymentRequests `has_many:"payment_requests"`
	MTOShipments       MTOShipments    `has_many:"mto_shipments"`
	MoveOrderID        uuid.UUID       `db:"move_order_id"`
	ReferenceID        string          `db:"reference_id"`
	IsAvailableToPrime bool            `db:"is_available_to_prime"`
	IsCanceled         bool            `db:"is_canceled"`
	PPMEstimatedWeight *unit.Pound     `db:"ppm_estimated_weight"`
	PPMType            *string         `db:"ppm_type"`
	ContractorID       uuid.UUID       `db:"contractor_id"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}

// MoveTaskOrders is a list of move task orders
type MoveTaskOrders []MoveTaskOrder

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MoveTaskOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs,
		&validators.UUIDIsPresent{Field: m.MoveOrderID, Name: "MoveOrderID"},
		&validators.StringIsPresent{Field: m.ReferenceID, Name: "ReferenceID"},
		&validators.UUIDIsPresent{Field: m.ContractorID, Name: "ContractorID"})
	return validate.Validate(vs...), nil
}
