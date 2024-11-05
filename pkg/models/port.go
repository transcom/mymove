package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PortType represents the type of port
type PortType string

// String is a string PortType
func (p PortType) String() string {
	return string(p)
}

const (
	PortTypeAir     PortType = "A"
	PortTypeSurface PortType = "S"
	PortTypeBoth    PortType = "B"
)

type Port struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	PortCode  string    `json:"port_code" db:"port_code" rw:"r"`
	PortType  PortType  `json:"port_type" db:"port_type" rw:"r"`
	PortName  string    `json:"port_name" db:"port_name" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

func (p Port) TableName() string {
	return "ports"
}
