package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
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
	ServiceArea   int       `json:"service_area" db:"service_area"`
	RateArea      string    `json:"rate_area" db:"rate_area"`
	Region        int       `json:"region" db:"region"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngZip3) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngZip3s is not required by pop and may be deleted
type Tariff400ngZip3s []Tariff400ngZip3

// String is not required by pop and may be deleted
func (t Tariff400ngZip3s) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngZip3) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Zip3, Name: "Zip3"},
		&validators.StringIsPresent{Field: t.BasepointCity, Name: "BasepointCity"},
		&validators.StringIsPresent{Field: t.State, Name: "State"},
		&validators.IntIsPresent{Field: t.ServiceArea, Name: "ServiceArea"},
		&validators.StringIsPresent{Field: t.RateArea, Name: "RateArea"},
		&validators.IntIsPresent{Field: t.Region, Name: "Region"},
	), nil
}

// FetchRateAreaForZip5 returns the rate area for a specified zip5.
func FetchRateAreaForZip5(db *pop.Connection, zip5 string) (string, error) {
	if len(zip5) != 5 {
		return "", errors.Errorf("zip5 must have a length of 5, got '%s'", zip5)
	}
	zip3 := zip5[0:3]

	tariffZip3 := Tariff400ngZip3{}
	if err := db.Where("zip3 = ?", zip3).First(&tariffZip3); err != nil {
		return "", errors.Wrapf(err, "could not find zip3 for %s", zip3)
	}

	if tariffZip3.RateArea == "ZIP" {
		zip5RateArea := Tariff400ngZip5RateArea{}
		if err := db.Where("zip5 = ?", zip5).First(&zip5RateArea); err != nil {
			return "", errors.Wrapf(err, "could not find zip5_rate_area for %s", zip5)
		}
		return zip5RateArea.RateArea, nil
	}

	return tariffZip3.RateArea, nil
}

// FetchRegionForZip5 returns the region for a specified zip5.
func FetchRegionForZip5(db *pop.Connection, zip5 string) (int, error) {
	if len(zip5) != 5 {
		return 0, errors.Errorf("zip5 must have a length of 5, got '%s'", zip5)
	}
	zip3 := zip5[0:3]

	tariffZip3 := Tariff400ngZip3{}
	if err := db.Where("zip3 = ?", zip3).First(&tariffZip3); err != nil {
		return 0, errors.Wrapf(err, "could not find zip3 for %s", zip3)
	}

	return tariffZip3.Region, nil
}
