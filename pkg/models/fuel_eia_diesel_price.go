package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// FuelEIADieselPrice used to hold data from the SDDC Fuel Surcharge information
// found at https://etops.sddc.army.mil/pls/ppcig_camp/fsc.output to calculate a
// shipment's fuel surcharge
type FuelEIADieselPrice struct {
	ID                          uuid.UUID       `json:"id" db:"id"`
	CreatedAt                   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at" db:"updated_at"`
	PubDate                     time.Time       `json:"pub_date" db:"pub_date"`
	RateStartDate               time.Time       `json:"rate_start_date" db:"rate_start_date"`
	RateEndDate                 time.Time       `json:"rate_end_date" db:"rate_end_date"`
	EIAPricePerGallonMillicents unit.Millicents `json:"eia_price_per_gallon_millicents" db:"eia_price_per_gallon_millicents"`
	BaselineRate                int64           `json:"baseline_rate" db:"baseline_rate"`
}

// String is not required by pop and may be deleted
func (f FuelEIADieselPrice) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// FuelEIADieselPrices is not required by pop and may be deleted
type FuelEIADieselPrices []FuelEIADieselPrice

// String is not required by pop and may be deleted
func (f FuelEIADieselPrices) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *FuelEIADieselPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: f.PubDate, Name: "PubDate"},
		&validators.TimeIsPresent{Field: f.RateStartDate, Name: "RateStartDate"},
		&validators.TimeIsPresent{Field: f.RateEndDate, Name: "RateEndDate"},
		&validators.IntIsGreaterThan{Field: f.EIAPricePerGallonMillicents.Int(), Name: "EIAPricePerGallonMillicents", Compared: 0},
		&validators.IntIsGreaterThan{Field: int(f.BaselineRate), Name: "BaselineRate", Compared: -1},
		&validators.IntIsLessThan{Field: int(f.BaselineRate), Name: "BaselineRate", Compared: 101},
		&validators.TimeAfterTime{
			FirstTime: f.RateEndDate, FirstName: "RateEndDate",
			SecondTime: f.RateStartDate, SecondName: "RateStartDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *FuelEIADieselPrice) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *FuelEIADieselPrice) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
