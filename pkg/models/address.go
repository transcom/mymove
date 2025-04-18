package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/utils"
)

const STREET_ADDRESS_1_NOT_PROVIDED string = "n/a"

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
	County             *string           `json:"county" db:"county"`
	IsOconus           *bool             `json:"is_oconus" db:"is_oconus"`
	UsPostRegionCityID *uuid.UUID        `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	UsPostRegionCity   *UsPostRegionCity `belongs_to:"us_post_region_cities" fk_id:"us_post_region_cities_id"`
	DestinationGbloc   *string           `db:"-"` // this tells Pop not to look in the db for this value
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
	if err := dbConnection.Q().Eager("Country", "UsPostRegionCity").Find(&address, id); err != nil {
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
func (a *Address) Validate(dbConnection *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringIsPresent{Field: a.StreetAddress1, Name: "StreetAddress1"})
	vs = append(vs, &validators.StringIsPresent{Field: a.City, Name: "City"})
	vs = append(vs, &validators.StringIsPresent{Field: a.State, Name: "State"})
	vs = append(vs, &validators.StringIsPresent{Field: a.PostalCode, Name: "PostalCode"})

	if a.IsOconus != nil && !*a.IsOconus {
		vs = append(vs, &validators.UUIDIsPresent{Field: *a.UsPostRegionCityID, Name: "UsPostRegionCityID"})
	}

	var validPostalCode bool
	if dbConnection != nil {
		var err error
		validPostalCode, err = ValidPostalCode(dbConnection, a.PostalCode)
		if err != nil {
			return nil, err
		}
	}

	if dbConnection != nil && a.UsPostRegionCityID != nil && *a.UsPostRegionCityID != uuid.Nil && validPostalCode {
		validUSPRC, err := ValidateUsPostRegionCityID(dbConnection, *a)
		if err != nil {
			return nil, err
		}

		if !validUSPRC {
			vs = append(vs, &validators.StringsMatch{Field: strconv.FormatBool(validUSPRC), Field2: "true", Name: "UsPostRegionCityID", Message: "UsPostRegionCityID is invalid."})
		}
	}

	return validate.Validate(vs...), nil
}

// Validate an addresses USPRC assignment
func ValidateUsPostRegionCityID(db *pop.Connection, address Address) (bool, error) {

	if address.UsPostRegionCityID != nil && strings.TrimSpace(address.City) != "" && strings.TrimSpace(address.PostalCode) != "" {
		expectedUSPRC, err := FindByZipCodeAndCity(db, address.PostalCode, address.City)
		if err != nil {
			return false, err
		}

		if expectedUSPRC.ID == *address.UsPostRegionCityID {
			return true, nil
		}
	}

	return false, nil
}

func ValidPostalCode(db *pop.Connection, postalCode string) (bool, error) {

	zipCount, err := db.Where("uspr_zip_id = $1", postalCode).CountByField(&UsPostRegionCity{}, "uspr_zip_id")
	if err != nil {
		return false, err
	}

	if zipCount == 0 {
		return false, nil
	}

	if len(strings.TrimSpace(postalCode)) != 5 {
		return false, nil
	}

	return true, nil
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

// LineDisplayFormat returns the address in a single line representation of the US mailing address format
func (a *Address) LineDisplayFormat() string {
	optionalStreetAddress2 := ""
	if a.StreetAddress2 != nil && len(*a.StreetAddress2) > 0 {
		optionalStreetAddress2 = " " + *a.StreetAddress2
	}
	optionalStreetAddress3 := ""
	if a.StreetAddress3 != nil && len(*a.StreetAddress3) > 0 {
		optionalStreetAddress3 = " " + *a.StreetAddress3
	}

	return fmt.Sprintf("%s%s%s, %s, %s %s", a.StreetAddress1, optionalStreetAddress2, optionalStreetAddress3, a.City, a.State, a.PostalCode)
}

func (a *Address) IsAddressAlaska() (bool, error) {
	if a == nil {
		return false, errors.New("address is nil")
	}
	return a.State == "AK", nil
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
func IsAddressEmpty(a *Address) bool {
	return a == nil || ((utils.IsNullOrWhiteSpace(&a.StreetAddress1) || IsDefaultAddressValue(a.StreetAddress1)) &&
		(utils.IsNullOrWhiteSpace(&a.City) || IsDefaultAddressValue(a.City)) &&
		(utils.IsNullOrWhiteSpace(&a.State) || IsDefaultAddressValue(a.State)) &&
		(utils.IsNullOrWhiteSpace(&a.PostalCode) || IsDefaultAddressValue(a.PostalCode)))
}
func IsDefaultAddressValue(s string) bool {
	return s == "n/a"
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

// Fetches the GBLOC for a specific Address (for now this will be used for OCONUS)
func FetchAddressGbloc(db *pop.Connection, address Address, serviceMember ServiceMember) (*string, error) {
	var gbloc *string

	err := db.RawQuery("SELECT * FROM get_address_gbloc($1, $2)", address.ID, serviceMember.Affiliation.String()).
		First(&gbloc)

	if err != nil {
		return nil, err
	}

	return gbloc, nil
}
