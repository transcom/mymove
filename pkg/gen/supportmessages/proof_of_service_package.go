// Code generated by go-swagger; DO NOT EDIT.

package supportmessages

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

// ProofOfServicePackage proof of service package
//
// swagger:model ProofOfServicePackage
type ProofOfServicePackage struct {

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// uploads
	Uploads []*UploadWithOmissions `json:"uploads"`
}

// Validate validates this proof of service package
func (m *ProofOfServicePackage) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUploads(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProofOfServicePackage) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ProofOfServicePackage) validateUploads(formats strfmt.Registry) error {
	if swag.IsZero(m.Uploads) { // not required
		return nil
	}

	for i := 0; i < len(m.Uploads); i++ {
		if swag.IsZero(m.Uploads[i]) { // not required
			continue
		}

		if m.Uploads[i] != nil {
			if err := m.Uploads[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("uploads" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("uploads" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this proof of service package based on the context it is used
func (m *ProofOfServicePackage) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateUploads(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProofOfServicePackage) contextValidateUploads(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Uploads); i++ {

		if m.Uploads[i] != nil {

			if swag.IsZero(m.Uploads[i]) { // not required
				return nil
			}

			if err := m.Uploads[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("uploads" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("uploads" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProofOfServicePackage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProofOfServicePackage) UnmarshalBinary(b []byte) error {
	var res ProofOfServicePackage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}