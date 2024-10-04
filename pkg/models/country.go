package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Country is a model representing a country
type Country struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Country     string    `json:"country" db:"country"`
	CountryName string    `json:"country_name" db:"country_name"`
}

// TableName overrides the table name used by Pop.
func (c Country) TableName() string {
	return "re_countries"
}

// fetches countries by the two digit code
func FetchCountryByCode(db *pop.Connection, code string) (Country, error) {
	var country Country
	err := db.Where("country = ?", code).First(&country)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Country{}, errors.Wrap(ErrFetchNotFound, "the country code provided in the request was not found")
		}
		return Country{}, err
	}

	return country, nil
}

// fetches countries by the two digit code
func FetchCountryByID(db *pop.Connection, id uuid.UUID) (Country, error) {
	var country Country
	err := db.Q().Find(&country, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Country{}, ErrFetchNotFound
		}
		return Country{}, err
	}

	return country, nil
}
