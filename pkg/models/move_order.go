package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type MoveOrder struct {
	ID                       uuid.UUID   `db:"id"`
	CreatedAt                time.Time   `db:"created_at"`
	UpdatedAt                time.Time   `db:"updated_at"`
	Customer                 Customer    `belongs_to:"customers"`
	CustomerID               uuid.UUID   `db:"customer_id"`
	Entitlement              Entitlement `belongs_to:"entitlements"`
	EntitlementID            uuid.UUID   `db:"entitlement_id"`
	DestinationDutyStation   DutyStation `belongs_to:"duty_stations"`
	DestinationDutyStationID uuid.UUID   `db:"destination_duty_station_id"`
	OriginDutyStation        DutyStation `belongs_to:"duty_stations"`
	OriginDutyStationID      uuid.UUID   `db:"origin_duty_station_id"`
}

func (m *MoveOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.CustomerID, Name: "CustomerID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.EntitlementID, Name: "EntitlementID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.DestinationDutyStationID, Name: "DestinationDutyStationID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.OriginDutyStationID, Name: "OriginDutyStationID"})
	return validate.Validate(vs...), nil
}
