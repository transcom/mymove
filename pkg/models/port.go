package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Port struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PortCode  string    `json:"port_code" db:"port_code"`
	PortType  string    `json:"port_type" db:"port_type"`
	PortName  string    `json:"port_name" db:"port_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (p Port) TableName() string {
	return "ports"
}
