package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type PortLocation struct {
	ID           uuid.UUID `json:"id" db:"id"`
	PortID       uuid.UUID `json:"port_id" db:"port_id"`
	City         string    `json:"city" db:"city"`
	County       string    `json:"county" db:"county"`
	State        string    `json:"state" db:"state"`
	Zip5         string    `json:"zip5" db:"zip5"`
	Country      string    `json:"country" db:"country"`
	InactiveFlag string    `json:"inactive_flag" db:"inactive_flag"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (l PortLocation) TableName() string {
	return "port_locations"
}
