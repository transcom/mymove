// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RateEnginePostalCodePayload rate engine postal code payload
//
// swagger:model RateEnginePostalCodePayload
type RateEnginePostalCodePayload struct {

	// ZIP
	//
	// zip code, international allowed
	// Example: '90210' or 'N15 3NL'
	// Required: true
	PostalCode *string `json:"postal_code"`

	// postal code type
	// Required: true
	// Enum: [origin destination]
	PostalCodeType *string `json:"postal_code_type"`

	// valid
	// Example: false
	// Required: true
	Valid *bool `json:"valid"`
}

// Validate validates this rate engine postal code payload
func (m *RateEnginePostalCodePayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePostalCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePostalCodeType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValid(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RateEnginePostalCodePayload) validatePostalCode(formats strfmt.Registry) error {

	if err := validate.Required("postal_code", "body", m.PostalCode); err != nil {
		return err
	}

	return nil
}

var rateEnginePostalCodePayloadTypePostalCodeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["origin","destination"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		rateEnginePostalCodePayloadTypePostalCodeTypePropEnum = append(rateEnginePostalCodePayloadTypePostalCodeTypePropEnum, v)
	}
}

const (

	// RateEnginePostalCodePayloadPostalCodeTypeOrigin captures enum value "origin"
	RateEnginePostalCodePayloadPostalCodeTypeOrigin string = "origin"

	// RateEnginePostalCodePayloadPostalCodeTypeDestination captures enum value "destination"
	RateEnginePostalCodePayloadPostalCodeTypeDestination string = "destination"
)

// prop value enum
func (m *RateEnginePostalCodePayload) validatePostalCodeTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, rateEnginePostalCodePayloadTypePostalCodeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *RateEnginePostalCodePayload) validatePostalCodeType(formats strfmt.Registry) error {

	if err := validate.Required("postal_code_type", "body", m.PostalCodeType); err != nil {
		return err
	}

	// value enum
	if err := m.validatePostalCodeTypeEnum("postal_code_type", "body", *m.PostalCodeType); err != nil {
		return err
	}

	return nil
}

func (m *RateEnginePostalCodePayload) validateValid(formats strfmt.Registry) error {

	if err := validate.Required("valid", "body", m.Valid); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this rate engine postal code payload based on context it is used
func (m *RateEnginePostalCodePayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RateEnginePostalCodePayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RateEnginePostalCodePayload) UnmarshalBinary(b []byte) error {
	var res RateEnginePostalCodePayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
