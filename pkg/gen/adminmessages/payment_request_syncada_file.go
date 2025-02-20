// Code generated by go-swagger; DO NOT EDIT.

package adminmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PaymentRequestSyncadaFile payment request syncada file
//
// swagger:model PaymentRequestSyncadaFile
type PaymentRequestSyncadaFile struct {

	// created at
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// edi string
	EdiString string `json:"ediString,omitempty"`

	// file name
	FileName string `json:"fileName,omitempty"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// payment request number
	PaymentRequestNumber string `json:"paymentRequestNumber,omitempty"`
}

// Validate validates this payment request syncada file
func (m *PaymentRequestSyncadaFile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PaymentRequestSyncadaFile) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PaymentRequestSyncadaFile) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this payment request syncada file based on the context it is used
func (m *PaymentRequestSyncadaFile) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PaymentRequestSyncadaFile) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PaymentRequestSyncadaFile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PaymentRequestSyncadaFile) UnmarshalBinary(b []byte) error {
	var res PaymentRequestSyncadaFile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
