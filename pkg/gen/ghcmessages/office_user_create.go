// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

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

// OfficeUserCreate office user create
//
// swagger:model OfficeUserCreate
type OfficeUserCreate struct {

	// EDIPI
	// Example: 1234567890
	// Max Length: 10
	Edipi *string `json:"edipi,omitempty"`

	// Email
	// Example: user@userdomain.com
	// Required: true
	Email string `json:"email"`

	// First Name
	// Required: true
	FirstName string `json:"firstName"`

	// Last Name
	// Required: true
	LastName string `json:"lastName"`

	// Middle Initials
	// Example: L.
	MiddleInitials *string `json:"middleInitials,omitempty"`

	// Office user identifier when EDIPI is not available
	OtherUniqueID *string `json:"otherUniqueId,omitempty"`

	// roles
	// Required: true
	Roles []*OfficeUserRole `json:"roles"`

	// telephone
	// Example: 212-555-5555
	// Required: true
	// Pattern: ^[2-9]\d{2}-\d{3}-\d{4}$
	Telephone string `json:"telephone"`

	// transportation office Id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Format: uuid
	TransportationOfficeID strfmt.UUID `json:"transportationOfficeId"`
}

// Validate validates this office user create
func (m *OfficeUserCreate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEdipi(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFirstName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTelephone(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOfficeID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserCreate) validateEdipi(formats strfmt.Registry) error {
	if swag.IsZero(m.Edipi) { // not required
		return nil
	}

	if err := validate.MaxLength("edipi", "body", *m.Edipi, 10); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateEmail(formats strfmt.Registry) error {

	if err := validate.RequiredString("email", "body", m.Email); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateFirstName(formats strfmt.Registry) error {

	if err := validate.RequiredString("firstName", "body", m.FirstName); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateLastName(formats strfmt.Registry) error {

	if err := validate.RequiredString("lastName", "body", m.LastName); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateRoles(formats strfmt.Registry) error {

	if err := validate.Required("roles", "body", m.Roles); err != nil {
		return err
	}

	for i := 0; i < len(m.Roles); i++ {
		if swag.IsZero(m.Roles[i]) { // not required
			continue
		}

		if m.Roles[i] != nil {
			if err := m.Roles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("roles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("roles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *OfficeUserCreate) validateTelephone(formats strfmt.Registry) error {

	if err := validate.RequiredString("telephone", "body", m.Telephone); err != nil {
		return err
	}

	if err := validate.Pattern("telephone", "body", m.Telephone, `^[2-9]\d{2}-\d{3}-\d{4}$`); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateTransportationOfficeID(formats strfmt.Registry) error {

	if err := validate.Required("transportationOfficeId", "body", strfmt.UUID(m.TransportationOfficeID)); err != nil {
		return err
	}

	if err := validate.FormatOf("transportationOfficeId", "body", "uuid", m.TransportationOfficeID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this office user create based on the context it is used
func (m *OfficeUserCreate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRoles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserCreate) contextValidateRoles(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Roles); i++ {

		if m.Roles[i] != nil {

			if swag.IsZero(m.Roles[i]) { // not required
				return nil
			}

			if err := m.Roles[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("roles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("roles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OfficeUserCreate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OfficeUserCreate) UnmarshalBinary(b []byte) error {
	var res OfficeUserCreate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
