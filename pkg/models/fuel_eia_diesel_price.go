package models

import (
	"encoding/json"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

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

// TableName overrides the table name used by Pop.
func (f FuelEIADieselPrice) TableName() string {
	return "fuel_eia_diesel_prices"
}

func (f FuelEIADieselPrice) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

type FuelEIADieselPrices []FuelEIADieselPrice

func (f FuelEIADieselPrices) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// FetchMostRecentFuelPrices queries and fetches all fuel_eia_diesel_prices for past specified number of months, including this month
func FetchMostRecentFuelPrices(dbConnection *pop.Connection, clock clock.Clock, numMonths int) ([]FuelEIADieselPrice, error) {
	today := clock.Now().UTC()

	query := dbConnection.Where("pub_date BETWEEN $1 AND $2", today.AddDate(0, -numMonths, 0), today)

	var fuelPrices FuelEIADieselPrices
	err := query.Eager().All(&fuelPrices)

	if err != nil {
		return fuelPrices, errors.Wrap(err, "Fetch line items query failed")
	}
	return fuelPrices, nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (f *FuelEIADieselPrice) Validate(_ *pop.Connection) (*validate.Errors, error) {
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
