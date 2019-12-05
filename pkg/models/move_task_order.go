package models

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type MoveTaskOrder struct {
	ID                 uuid.UUID           `db:"id"`
	MoveOrder          MoveOrder           `belongs_to:"move_orders"`
	MoveOrderID        uuid.UUID           `db:"move_order_id"`
	ReferenceID        *string             `db:"reference_id"`
	Status             MoveTaskOrderStatus `db:"status"`
	IsAvailableToPrime bool                `db:"is_available_to_prime"`
	IsCancelled        bool                `db:"is_cancelled"`
	CreatedAt          time.Time           `db:"created_at"`
	UpdatedAt          time.Time           `db:"updated_at"`
}

type MoveTaskOrderStatus string

const (
	MoveTaskOrderStatusApproved MoveTaskOrderStatus = "APPROVED"
	MoveTaskOrderStatusDraft    MoveTaskOrderStatus = "DRAFT"
	MoveTaskOrderStatusRejected MoveTaskOrderStatus = "REJECTED"
)

func (m *MoveTaskOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
	), nil
}

type MoveTaskOrders []MoveTaskOrder

// GenerateReferenceID creates a random ID for an MTO. Format (xxxx-xxxx) with X being a number 0-9 (ex. 0009-1234. 4321-4444)
func generateReferenceID(tx *pop.Connection) (string, error) {
	min := 0
	max := 9999
	firstNum := rand.Intn(max - min + 1)
	secondNum := rand.Intn(max - min + 1)
	newReferenceID := fmt.Sprintf("%04d-%04d", firstNum, secondNum)
	count, err := tx.Where(`reference_id= $1`, newReferenceID).Count(&MoveTaskOrder{})
	if err != nil || count > 0 {
		return "", err
	}
	return newReferenceID, nil
}

func GenerateReferenceID(tx *pop.Connection) (*string, error) {
	const maxAttempts = 10
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceID(tx)
		if err == nil {
			return &referenceID, nil
		}
	}
	return nil, errors.New("move_task_order: failed to generate reference id")
}
