package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"time"
)

// QueueMove is a move for a queue
type QueueMove struct {
	ID               uuid.UUID                           `json:"id" db:"id"`
	CreatedAt        time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                           `json:"updated_at" db:"updated_at"`
	Edipi            *string                             `json:"edipi" db:"edipi"`
	Rank             *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	CustomerName     *string                             `json:"customer_name" db:"customer_name"`
	Status           *string                             `json:"status" db:"status"`
	Locator          *string                             `json:"locator" db:"locator"`
	MoveType         *string                             `json:"move_type" db:"move_type"`
	MoveDate         time.Time                           `json:"move_date" db:"move_date"`
	CustomerDeadline time.Time                           `json:"customer_deadline" db:"customer_deadline"`
	LastModifiedDate *time.Time                          `json:"last_modified_date" db:"last_modified_date"`
	LastModifiedName *string                             `json:"last_modified_name" db:"last_modified_name"`
}

// String is not required by pop and may be deleted
func (q QueueMove) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// QueueMoves is not required by pop and may be deleted
type QueueMoves []QueueMove

// String is not required by pop and may be deleted
func (q QueueMoves) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *QueueMove) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *QueueMove) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *QueueMove) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
