// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RequestDiversion request diversion
//
// swagger:model RequestDiversion
type RequestDiversion struct {

	// diversion reason
	// Example: Shipment route needs to change
	// Required: true
	DiversionReason *string `json:"diversionReason"`
}

// Validate validates this request diversion
func (m *RequestDiversion) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDiversionReason(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RequestDiversion) validateDiversionReason(formats strfmt.Registry) error {

	if err := validate.Required("diversionReason", "body", m.DiversionReason); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this request diversion based on context it is used
func (m *RequestDiversion) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RequestDiversion) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RequestDiversion) UnmarshalBinary(b []byte) error {
	var res RequestDiversion
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}