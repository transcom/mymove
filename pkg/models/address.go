package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

// String is not required by pop and may be deleted
func (a Address) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
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
		if err.Error() != "sql: no rows in result set" {
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

// String is not required by pop and may be deleted
func (a Addresses) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Address) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	stringFields := map[string]string{
		"StreetAddress1": a.StreetAddress1,
		"City":           a.City,
		"State":          a.State,
		"PostalCode":     a.PostalCode,
	}

	for key, field := range stringFields {
		if field == "" {
			verrs.Add(key, fmt.Sprintf("%s must not be blank!", key))
		}
	}

	return verrs, nil
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

// AddressModelFromPayload returns the model struct from payload
func AddressModelFromPayload(rawAddress *internalmessages.Address) *Address {
	if rawAddress == nil {
		return nil
	}
	address := Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		StreetAddress3: rawAddress.StreetAddress3,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		PostalCode:     *rawAddress.PostalCode,
		Country:        rawAddress.Country,
	}
	return &address
}
