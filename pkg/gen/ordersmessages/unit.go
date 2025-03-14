// Code generated by go-swagger; DO NOT EDIT.

package ordersmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Unit Information about either the losing or gaining Unit. If these are separation orders, the location information for the gaining Unit may be blank.
//
// swagger:model Unit
type Unit struct {

	// May be FPO or APO for OCONUS commands.
	City *string `json:"city,omitempty"`

	// ISO 3166-1 alpha-2 country code. If blank, but city and locality or postalCode are not blank, assume US
	// Pattern: ^[A-Z]{2}$
	Country *string `json:"country,omitempty"`

	// State (US). OCONUS units may not have the equivalent information available.
	Locality *string `json:"locality,omitempty"`

	// Human-readable name of the Unit.
	Name *string `json:"name,omitempty"`

	// In the USA, this is the ZIP Code.
	PostalCode *string `json:"postalCode,omitempty"`

	// Unit Identification Code - a six character alphanumeric code that uniquely identifies each United States Department of Defense entity. Used in Army, Air Force, and Navy orders.
	// Note that the Navy has the habit of omitting the leading character, which is always "N" for them.
	//
	// Pattern: ^[A-Z][A-Z0-9]{5}$
	Uic *string `json:"uic,omitempty"`
}

// Validate validates this unit
func (m *Unit) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCountry(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUic(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Unit) validateCountry(formats strfmt.Registry) error {
	if swag.IsZero(m.Country) { // not required
		return nil
	}

	if err := validate.Pattern("country", "body", *m.Country, `^[A-Z]{2}$`); err != nil {
		return err
	}

	return nil
}

func (m *Unit) validateUic(formats strfmt.Registry) error {
	if swag.IsZero(m.Uic) { // not required
		return nil
	}

	if err := validate.Pattern("uic", "body", *m.Uic, `^[A-Z][A-Z0-9]{5}$`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this unit based on context it is used
func (m *Unit) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Unit) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Unit) UnmarshalBinary(b []byte) error {
	var res Unit
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
