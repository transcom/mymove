package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type MoveTaskOrder struct {
	ID                               uuid.UUID           `db:"id"`
	ActualWeight                     *unit.Pound         `json:"actual_weight" db:"actual_weight"`
	CreatedAt                        time.Time           `db:"created_at"`
	Customer                         ServiceMember       `belongs_to:"service_members"`
	CustomerID                       uuid.UUID           `db:"customer_id"`
	CustomerRemarks                  string              `db:"customer_remarks"`
	DestinationAddress               Address             `belongs_to:"addresses"`
	DestinationAddressID             uuid.UUID           `db:"destination_address_id"`
	DestinationDutyStation           DutyStation         `belongs_to:"duty_stations"`
	DestinationDutyStationID         uuid.UUID           `db:"destination_duty_station_id"`
	Entitlements                     GHCEntitlement      `has_one:"entitlements"`
	Move                             Move                `belongs_to:"moves"`
	MoveID                           uuid.UUID           `db:"move_id"`
	OriginDutyStation                DutyStation         `belongs_to:"duty_stations"`
	OriginDutyStationID              uuid.UUID           `db:"origin_duty_station_id"`
	PickupAddress                    Address             `belongs_to:"addresses"`
	PickupAddressID                  uuid.UUID           `db:"pickup_address_id"`
	PrimeEstimatedWeight             *unit.Pound         `db:"prime_estimated_weight"`
	PrimeEstimatedWeightRecordedDate *time.Time          `db:"prime_estimated_weight_recorded_date"`
	ReferenceID                      string              `db:"reference_id"`
	RequestedPickupDate              time.Time           `db:"requested_pickup_date"`
	Status                           MoveTaskOrderStatus `db:"status"`
	ServiceItems                     ServiceItems        `has_many:"service_items"`
	UpdatedAt                        time.Time           `db:"updated_at"`
}

type MoveTaskOrderStatus string

const (
	MoveTaskOrderStatusApproved  MoveTaskOrderStatus = "APPROVED"
	MoveTaskOrderStatusSubmitted MoveTaskOrderStatus = "SUBMITTED"
	MoveTaskOrderStatusRejected  MoveTaskOrderStatus = "REJECTED"
	MoveTaskOrderStatusDraft     MoveTaskOrderStatus = "DRAFT"
)

func (m *MoveTaskOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringIsPresent{Field: string(m.Status), Name: "Status"})
	if m.PrimeEstimatedWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeEstimatedWeight.Int(), Compared: -1, Name: "PrimeEstimatedWeight"})
	}
	return validate.Validate(vs...), nil
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

const maxAttempts = 10

func GenerateReferenceID(tx *pop.Connection) (string, error) {
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceID(tx)
		if err == nil {
			return referenceID, nil
		}
		if i >= maxAttempts {
			break
		}
	}
	return "", err
}
