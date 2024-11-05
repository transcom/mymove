package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Address is an address
type Address struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`
	StreetAddress1     string            `json:"street_address_1" db:"street_address_1"`
	StreetAddress2     *string           `json:"street_address_2" db:"street_address_2"`
	StreetAddress3     *string           `json:"street_address_3" db:"street_address_3"`
	City               string            `json:"city" db:"city"`
	State              string            `json:"state" db:"state"`
	PostalCode         string            `json:"postal_code" db:"postal_code"`
	CountryId          *uuid.UUID        `json:"country_id" db:"country_id"`
	Country            *Country          `belongs_to:"re_countries" fk_id:"country_id"`
	County             string            `json:"county" db:"county"`
	IsOconus           *bool             `json:"is_oconus" db:"is_oconus"`
	UsPostRegionCityId *uuid.UUID        `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	UsPostRegionCity   *UsPostRegionCity `belongs_to:"us_post_region_cities" fk_id:"us_post_region_cities_id"`
}

// TableName overrides the table name used by Pop.
func (a Address) TableName() string {
	return "addresses"
}

// FetchAddressByID returns an address model by ID
func FetchAddressByID(dbConnection *pop.Connection, id *uuid.UUID) *Address {
	if id == nil {
		return nil
	}
	address := Address{}
	var response *Address
	if err := dbConnection.Q().Eager("Country").Find(&address, id); err != nil {
		response = nil
		if err.Error() != RecordNotFoundErrorString {
			// This is an unknown error from the db
			zap.L().Error("DB Insertion error", zap.Error(err))
		}
	} else {
		response = &address
	}
	return response
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (a *Address) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.StreetAddress1, Name: "StreetAddress1"},
		&validators.StringIsPresent{Field: a.City, Name: "City"},
		&validators.StringIsPresent{Field: a.State, Name: "State"},
		&validators.StringIsPresent{Field: a.PostalCode, Name: "PostalCode"},
		&validators.StringIsPresent{Field: a.County, Name: "County"},
	), nil
}

// MarshalLogObject is required to be able to zap.Object log TDLs
func (a *Address) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("street1", a.StreetAddress1)
	if a.StreetAddress2 != nil {
		encoder.AddString("street2", *a.StreetAddress2)
	}
	if a.StreetAddress3 != nil {
		encoder.AddString("street3", *a.StreetAddress3)
	}
	encoder.AddString("city", a.City)
	encoder.AddString("state", a.State)
	encoder.AddString("code", a.PostalCode)
	encoder.AddString("countryId", a.CountryId.String())
	return nil
}

// Format returns the address in default US mailing address format
func (a *Address) Format() string {
	lines := []string{}
	lines = append(lines, a.StreetAddress1)

	if a.StreetAddress2 != nil && len(*a.StreetAddress2) > 0 {
		lines = append(lines, *a.StreetAddress2)
	}
	if a.StreetAddress3 != nil && len(*a.StreetAddress3) > 0 {
		lines = append(lines, *a.StreetAddress3)
	}

	lines = append(lines, fmt.Sprintf("%s, %s %s", a.City, a.State, a.PostalCode))

	return strings.Join(lines, "\n")
}

// LineFormat returns the address as a string, formatted into a single line
func (a *Address) LineFormat() string {
	parts := []string{}
	if len(a.StreetAddress1) > 0 {
		parts = append(parts, a.StreetAddress1)
	}
	if a.StreetAddress2 != nil && len(*a.StreetAddress2) > 0 {
		parts = append(parts, *a.StreetAddress2)
	}
	if a.StreetAddress3 != nil && len(*a.StreetAddress3) > 0 {
		parts = append(parts, *a.StreetAddress3)
	}
	if len(a.City) > 0 {
		parts = append(parts, a.City)
	}
	if len(a.State) > 0 {
		parts = append(parts, a.State)
	}
	if len(a.PostalCode) > 0 {
		parts = append(parts, a.PostalCode)
	}
	if len(*a.CountryId) > 0 {
		parts = append(parts, a.Country.CountryName)
	}

	return strings.Join(parts, ", ")
}

// NotImplementedCountryCode is the default for unimplemented country code lookup
type NotImplementedCountryCode struct {
	message string
}

func (e NotImplementedCountryCode) Error() string {
	return fmt.Sprintf("NotImplementedCountryCode: %s", e.message)
}

// CountryCode returns 2-3 character code for country, returns nil if no Country
func (a *Address) CountryCode() (*string, error) {
	if a.Country != nil {
		return &a.Country.Country, nil
	}
	return nil, nil
}

// Copy returns a pointer that is a copy of the original pointer Address
func (a *Address) Copy() *Address {
	if a != nil {
		address := *a
		return &address
	}
	return nil
}

// Check if an address is CONUS or OCONUS
func IsAddressOconus(db *pop.Connection, address Address) (bool, error) {
	// use the data we have first, if it's not nil
	if address.Country != nil {
		isOconus := EvaluateIsOconus(address)
		return isOconus, nil
	} else if address.CountryId != nil {
		country, err := FetchCountryByID(db, *address.CountryId)
		if err != nil {
			return false, err
		}
		address.Country = &country
		isOconus := EvaluateIsOconus(address)
		return isOconus, nil
	} else {
		if address.State == "HI" || address.State == "AK" {
			return true, nil
		}
		return false, nil
	}
}

// Conditional logic for a CONUS and OCONUS address
func EvaluateIsOconus(address Address) bool {
	if address.Country.Country != "US" || address.Country.Country == "US" && address.State == "AK" || address.Country.Country == "US" && address.State == "HI" {
		return true
	} else {
		return false
	}
}
