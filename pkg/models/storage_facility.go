package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type StorageFacility struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	FacilityName string     `json:"facility_name" db:"facility_name"`
	Address      Address    `belongs_to:"addresses" fk_id:"address_id"`
	AddressID    uuid.UUID  `json:"address_id" db:"address_id"`
	LotNumber    *string    `json:"lot_number" db:"lot_number"`
	Phone        *string    `json:"phone" db:"phone"`
	Email        *string    `json:"email" db:"email"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

type StorageFacilities []StorageFacility

func (f *StorageFacility) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: f.AddressID, Name: "AddressID"},
		&validators.StringIsPresent{Field: f.FacilityName, Name: "FacilityName"},
		&StringIsNilOrNotBlank{Field: f.LotNumber, Name: "LotNumber"},
		&StringIsNilOrNotBlank{Field: f.Phone, Name: "Phone"},
		&StringIsNilOrNotBlank{Field: f.Email, Name: "Email"},
	), nil
}
