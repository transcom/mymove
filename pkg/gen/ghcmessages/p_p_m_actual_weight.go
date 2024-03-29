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

// PPMActualWeight The actual net weight of a single PPM shipment. Used during document review for PPM closeout.
//
// swagger:model PPMActualWeight
type PPMActualWeight struct {

	// actual weight
	// Example: 2000
	// Required: true
	ActualWeight *int64 `json:"actualWeight"`
}

// Validate validates this p p m actual weight
func (m *PPMActualWeight) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActualWeight(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PPMActualWeight) validateActualWeight(formats strfmt.Registry) error {

	if err := validate.Required("actualWeight", "body", m.ActualWeight); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this p p m actual weight based on context it is used
func (m *PPMActualWeight) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PPMActualWeight) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PPMActualWeight) UnmarshalBinary(b []byte) error {
	var res PPMActualWeight
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
