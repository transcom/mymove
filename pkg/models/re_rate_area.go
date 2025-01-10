package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReRateArea model struct
type ReRateArea struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ContractID uuid.UUID `json:"contract_id" db:"contract_id"`
	IsOconus   bool      `json:"is_oconus" db:"is_oconus"`
	Code       string    `json:"code" db:"code"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
}

// TableName overrides the table name used by Pop.
func (r ReRateArea) TableName() string {
	return "re_rate_areas"
}

type ReRateAreas []ReRateArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReRateArea) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.StringIsPresent{Field: r.Code, Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}

// FetchReRateAreaItem returns an area for a matching code
func FetchReRateAreaItem(tx *pop.Connection, contractID uuid.UUID, code string) (*ReRateArea, error) {
	var area ReRateArea
	err := tx.
		Where("contract_id = $1", contractID).
		Where("code = $2", code).
		First(&area)

	if err != nil {
		return nil, err
	}

	return &area, err
}

// a db stored proc that takes in an address id & a service code to get the rate area id for an address
func FetchRateAreaID(db *pop.Connection, addressID uuid.UUID, serviceID *uuid.UUID, contractID uuid.UUID) (uuid.UUID, error) {
	if addressID != uuid.Nil && contractID != uuid.Nil {
		var rateAreaID uuid.UUID
		err := db.RawQuery("SELECT get_rate_area_id($1, $2, $3)", addressID, serviceID, contractID).First(&rateAreaID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("error fetching rate area id for shipment ID: %s, service ID %s, and contract ID: %s: %s", addressID, serviceID, contractID, err)
		}
		return rateAreaID, nil
	}
	return uuid.Nil, fmt.Errorf("error fetching rate area ID - required parameters not provided")
}
