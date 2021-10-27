package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type StorageFacility struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	FacilityName *string    `json:"facility_name" db:"facility_name"`
	Address      *Address   `belongs_to:"addresses" fk_id:"address_id"`
	AddressID    *uuid.UUID `json:"address_id" db:"address_id"`
	LotNumber    *string    `json:"lot_number" db:"lot_number"`
	Phone        *string    `json:"phone" db:"phone"`
	Email        *string    `json:"email" db:"email"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type StorageFacilities []StorageFacility

func (r *StorageFacility) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}
