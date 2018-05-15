package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID                         uuid.UUID                    `json:"id" db:"id"`
	MoveID                     uuid.UUID                    `json:"move_id" db:"move_id"`
	Move                       Move                         `belongs_to:"move"`
	CreatedAt                  time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                    `json:"updated_at" db:"updated_at"`
	Size                       *internalmessages.TShirtSize `json:"size" db:"size"`
	WeightEstimate             *int64                       `json:"weight_estimate" db:"weight_estimate"`
	EstimatedIncentive         *string                      `json:"estimated_incentive" db:"estimated_incentive"`
	PlannedMoveDate            *time.Time                   `json:"planned_move_date" db:"planned_move_date"`
	PickupPostalCode           *string                      `json:"pickup_postal_code" db:"pickup_postal_code"`
	HasAdditionalPostalCode    *bool                        `json:"has_additional_postal_code" db:"has_additional_postal_code"`
	AdditionalPickupPostalCode *string                      `json:"additional_pickup_postal_code" db:"additional_pickup_postal_code"`
	DestinationPostalCode      *string                      `json:"destination_postal_code" db:"destination_postal_code"`
	HasSit                     *bool                        `json:"has_sit" db:"has_sit"`
	DaysInStorage              *int64                       `json:"days_in_storage" db:"days_in_storage"`
}

// PersonallyProcuredMoves is a list of PPMs
type PersonallyProcuredMoves []PersonallyProcuredMove

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchPersonallyProcuredMove Fetches and Validates a PPM model
func FetchPersonallyProcuredMove(db *pop.Connection, authUser User, id uuid.UUID) (*PersonallyProcuredMove, error) {
	var ppm PersonallyProcuredMove
	err := db.Q().Eager("Move.Orders.ServiceMember").Find(&ppm, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// TODO: Handle case where more than one user is authorized to modify ppm
	if ppm.Move.Orders.ServiceMember.UserID != authUser.ID {
		return nil, ErrFetchForbidden
	}

	return &ppm, nil
}
