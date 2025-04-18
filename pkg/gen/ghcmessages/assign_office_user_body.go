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

// AssignOfficeUserBody assign office user body
//
// swagger:model AssignOfficeUserBody
type AssignOfficeUserBody struct {

	// office user Id
	// Required: true
	// Format: uuid
	OfficeUserID *strfmt.UUID `json:"officeUserId"`

	// queue type
	// Required: true
	QueueType *string `json:"queueType"`
}

// Validate validates this assign office user body
func (m *AssignOfficeUserBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOfficeUserID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateQueueType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AssignOfficeUserBody) validateOfficeUserID(formats strfmt.Registry) error {

	if err := validate.Required("officeUserId", "body", m.OfficeUserID); err != nil {
		return err
	}

	if err := validate.FormatOf("officeUserId", "body", "uuid", m.OfficeUserID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *AssignOfficeUserBody) validateQueueType(formats strfmt.Registry) error {

	if err := validate.Required("queueType", "body", m.QueueType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this assign office user body based on context it is used
func (m *AssignOfficeUserBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AssignOfficeUserBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AssignOfficeUserBody) UnmarshalBinary(b []byte) error {
	var res AssignOfficeUserBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
