package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// DutyStationName represents a military duty station for a specific affiliation
type DutyStationName struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	Name          string      `json:"name" db:"name"`
	DutyStationID uuid.UUID   `json:"duty_station_id" db:"duty_station_id"`
	DutyStation   DutyStation `belongs_to:"duty_stations" fk_id:"duty_station_id"`
}

// DutyStationNames is not required by pop and may be deleted
type DutyStationNames []DutyStationName
