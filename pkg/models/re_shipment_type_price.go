package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReShipmentTypePrice model struct
type ReShipmentTypePrice struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ContractID       uuid.UUID `json:"contract_id" db:"contract_id"`
	ShipmentTypeID   uuid.UUID `json:"shipment_type_id" db:"shipment_type_id"`
	Market           Market    `json:"market" db:"market"`
	FactorHundredths int       `json:"factor_hundredths" db:"factor_hundredths"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	//Associations
	Contract     ReContract     `belongs_to:"re_contract"`
	ShipmentType ReShipmentType `belongs_to:"re_shipment_type"`
}

// ReShipmentTypePrices is not required by pop and may be deleted
type ReShipmentTypePrices []ReShipmentTypePrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReShipmentTypePrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ShipmentTypeID, Name: "ShipmentTypeID"},
		&validators.StringIsPresent{Field: string(r.Market), Name: "Market"},
		&validators.StringInclusion{Field: string(r.Market), Name: "Market", List: new(Market).String()},
		&validators.IntIsPresent{Field: r.FactorHundredths, Name: "FactorHundredths"},
		&validators.IntIsGreaterThan{Field: r.FactorHundredths, Name: "FactorHundredths", Compared: 0},
	), nil
}
