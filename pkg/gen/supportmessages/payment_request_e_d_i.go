// Code generated by go-swagger; DO NOT EDIT.

package supportmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PaymentRequestEDI payment request e d i
//
// swagger:model PaymentRequestEDI
type PaymentRequestEDI struct {

	// edi
	// Read Only: true
	Edi string `json:"edi,omitempty"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`
}

// Validate validates this payment request e d i
func (m *PaymentRequestEDI) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PaymentRequestEDI) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this payment request e d i based on the context it is used
func (m *PaymentRequestEDI) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEdi(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PaymentRequestEDI) contextValidateEdi(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "edi", "body", string(m.Edi)); err != nil {
		return err
	}

	return nil
}

func (m *PaymentRequestEDI) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PaymentRequestEDI) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PaymentRequestEDI) UnmarshalBinary(b []byte) error {
	var res PaymentRequestEDI
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}