// Code generated by go-swagger; DO NOT EDIT.

package primemessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ServiceItem service item
//
// swagger:model ServiceItem
type ServiceItem struct {

	// e tag
	// Read Only: true
	ETag string `json:"eTag,omitempty"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// This should be populated for the following service items:
	//   * DOASIT(Domestic origin Additional day SIT)
	//   * DDASIT(Domestic destination Additional day SIT)
	//
	// Both take in the following param keys:
	//   * `SITPaymentRequestStart`
	//   * `SITPaymentRequestEnd`
	//
	// The value of each is a date string in the format "YYYY-MM-DD" (e.g. "2023-01-15")
	//
	Params []*ServiceItemParamsItems0 `json:"params"`
}

// Validate validates this service item
func (m *ServiceItem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParams(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServiceItem) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ServiceItem) validateParams(formats strfmt.Registry) error {
	if swag.IsZero(m.Params) { // not required
		return nil
	}

	for i := 0; i < len(m.Params); i++ {
		if swag.IsZero(m.Params[i]) { // not required
			continue
		}

		if m.Params[i] != nil {
			if err := m.Params[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("params" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("params" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this service item based on the context it is used
func (m *ServiceItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateETag(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateParams(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServiceItem) contextValidateETag(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "eTag", "body", string(m.ETag)); err != nil {
		return err
	}

	return nil
}

func (m *ServiceItem) contextValidateParams(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Params); i++ {

		if m.Params[i] != nil {

			if swag.IsZero(m.Params[i]) { // not required
				return nil
			}

			if err := m.Params[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("params" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("params" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ServiceItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServiceItem) UnmarshalBinary(b []byte) error {
	var res ServiceItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ServiceItemParamsItems0 service item params items0
//
// swagger:model ServiceItemParamsItems0
type ServiceItemParamsItems0 struct {

	// key
	// Example: Service Item Parameter Name
	Key string `json:"key,omitempty"`

	// value
	// Example: Service Item Parameter Value
	Value string `json:"value,omitempty"`
}

// Validate validates this service item params items0
func (m *ServiceItemParamsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this service item params items0 based on context it is used
func (m *ServiceItemParamsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ServiceItemParamsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServiceItemParamsItems0) UnmarshalBinary(b []byte) error {
	var res ServiceItemParamsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}