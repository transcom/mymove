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
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	StreetAddress1 string    `json:"street_address_1" db:"street_address_1"`
	StreetAddress2 *string   `json:"street_address_2" db:"street_address_2"`
	StreetAddress3 *string   `json:"street_address_3" db:"street_address_3"`
	City           string    `json:"city" db:"city"`
	State          string    `json:"state" db:"state"`
	PostalCode     string    `json:"postal_code" db:"postal_code"`
	Country        *string   `json:"country" db:"country"`
}

// GetAddressID facilitates grabbing the ID from an address that may be nil
func GetAddressID(address *Address) *uuid.UUID {
	var response *uuid.UUID
	if address != nil {
		response = &address.ID
	}
	return response
}

// FetchAddressByID returns an address model by ID
func FetchAddressByID(dbConnection *pop.Connection, id *uuid.UUID) *Address {
	if id == nil {
		return nil
	}
	address := Address{}
	var response *Address
	if err := dbConnection.Find(&address, id); err != nil {
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

// Addresses is not required by pop and may be deleted
type Addresses []Address

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Address) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.StreetAddress1, Name: "StreetAddress1"},
		&validators.StringIsPresent{Field: a.City, Name: "City"},
		&validators.StringIsPresent{Field: a.State, Name: "State"},
		&validators.StringIsPresent{Field: a.PostalCode, Name: "PostalCode"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Address) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Address) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
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
	if a.Country != nil {
		encoder.AddString("country", *a.Country)
	}
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
	if a.Country != nil && len(*a.Country) > 0 {
		parts = append(parts, *a.Country)
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
// TODO: since we only support CONUS at this time this just returns USA and otherwise throws a NotImplementedCountryCode
func (a *Address) CountryCode() (*string, error) {
	if a.Country != nil && len(*a.Country) > 0 {
		result := ""
		switch *a.Country {
		case "United States", "US":
			result = "USA"
		default:
			return nil, NotImplementedCountryCode{message: fmt.Sprintf("Country '%s'", *a.Country)}
		}
		return &result, nil
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
