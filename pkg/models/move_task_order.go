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

// MoveTaskOrder is an object representing the task orders for a move
type MoveTaskOrder struct {
	ID                 uuid.UUID `db:"id"`
	MoveOrder          MoveOrder `belongs_to:"move_orders"`
	MoveOrderID        uuid.UUID `db:"move_order_id"`
	ReferenceID        *string   `db:"reference_id"`
	IsAvailableToPrime bool      `db:"is_available_to_prime"`
	IsCanceled         bool      `db:"is_canceled"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Customer		   Customer	 `belongs_to:customer`
	CustomerID		   uuid.UUID `db:"customer_id"`
	DestinationAddress       Address             `belongs_to:"addresses"`
	DestinationAddressID     uuid.UUID           `db:"destination_address_id"`
	DestinationDutyStation   DutyStation         `belongs_to:"duty_stations"`
	DestinationDutyStationID uuid.UUID           `db:"destination_duty_station_id"`
	Entitlements             GHCEntitlement      `has_one:"entitlements"`
	OriginDutyStation        DutyStation         `belongs_to:"duty_stations"`
	OriginDutyStationID      uuid.UUID           `db:"origin_duty_station_id"`
	PickupAddress            Address             `belongs_to:"addresses"`
	PickupAddressID          uuid.UUID           `db:"pickup_address_id"`
	RequestedPickupDate      time.Time           `db:"requested_pickup_date"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MoveTaskOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveOrderID, Name: "MoveOrderID"})
	return validate.Validate(vs...), nil
}

// MoveTaskOrders is a list of move task orders
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

func GenerateReferenceID(tx *pop.Connection) (string, error) {
	const maxAttempts = 10
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceID(tx)
		if err == nil {
			return referenceID, nil
		}
	}
	return "", errors.New("move_task_order: failed to generate reference id")
}
