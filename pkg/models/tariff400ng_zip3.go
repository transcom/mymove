package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Tariff400ngZip3 is the first 3 numbers of a zip for Tariff400NG calculations
type Tariff400ngZip3 struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Zip3          string    `json:"zip3" db:"zip3"`
	BasepointCity string    `json:"basepoint_city" db:"basepoint_city"`
	State         string    `json:"state" db:"state"`
	ServiceArea   string    `json:"service_area" db:"service_area"`
	RateArea      string    `json:"rate_area" db:"rate_area"`
	Region        string    `json:"region" db:"region"`
}

// Tariff400ngZip3s is not required by pop and may be deleted
type Tariff400ngZip3s []Tariff400ngZip3

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngZip3) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringLengthInRange{Field: t.Zip3, Name: "Zip3", Min: 3, Max: 3},
		&validators.StringIsPresent{Field: t.BasepointCity, Name: "BasepointCity"},
		&validators.StringIsPresent{Field: t.State, Name: "State"},
		&validators.StringIsPresent{Field: t.ServiceArea, Name: "ServiceArea"},
		&validators.RegexMatch{Field: t.ServiceArea, Name: "ServiceArea", Expr: "^[0-9]+$"},
		&validators.StringIsPresent{Field: t.RateArea, Name: "RateArea"},
		&validators.RegexMatch{Field: t.RateArea, Name: "RateArea", Expr: "^(ZIP|US[0-9]+)$"},
		&validators.StringIsPresent{Field: t.Region, Name: "Region"},
		&validators.RegexMatch{Field: t.Region, Name: "Region", Expr: "^[0-9]+$"},
	), nil
}

// FetchRateAreaForZip5 returns the rate area for a specified zip5.
func FetchRateAreaForZip5(db *pop.Connection, zip string) (string, error) {
	if len(zip) < 5 {
		return "", errors.Errorf("zip must have a length of at least 5, got '%s'", zip)
	}
	zip3 := zip[0:3]

	tariffZip3 := Tariff400ngZip3{}
	if err := db.Where("zip3 = $1", zip3).First(&tariffZip3); err != nil {
		return "", errors.Wrapf(err, "could not find zip3 for %s", zip3)
	}

	if tariffZip3.RateArea == "ZIP" {
		zip5 := zip[0:5]
		zip5RateArea := Tariff400ngZip5RateArea{}
		if err := db.Where("zip5 = $1", zip5).First(&zip5RateArea); err != nil {
			return "", errors.Wrapf(err, "could not find zip5_rate_area for %s", zip5)
		}
		return zip5RateArea.RateArea, nil
	}

	return tariffZip3.RateArea, nil
}

// FetchRegionForZip5 returns the region for a specified zip5.
func FetchRegionForZip5(db *pop.Connection, zip string) (string, error) {
	if len(zip) < 5 {
		return "", errors.Errorf("zip must have a length of at least 5, got '%s'", zip)
	}
	zip3 := zip[0:3]

	tariffZip3 := Tariff400ngZip3{}
	if err := db.Where("zip3 = $1", zip3).First(&tariffZip3); err != nil {
		return "", errors.Wrapf(err, "could not find zip3 for %s", zip3)
	}

	return tariffZip3.Region, nil
}
