// Code generated by go-swagger; DO NOT EDIT.

package adminmessages

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

	// Email
	// Example: user@userdomain.com
	Email string `json:"email,omitempty"`

	// First Name
	FirstName string `json:"firstName,omitempty"`

	// Last Name
	LastName string `json:"lastName,omitempty"`

	// Middle Initials
	// Example: L.
	MiddleInitials *string `json:"middleInitials,omitempty"`

	// privileges
	Privileges []*OfficeUserPrivilege `json:"privileges"`

	// roles
	Roles []*OfficeUserRole `json:"roles"`

	// telephone
	// Example: 212-555-5555
	// Pattern: ^[2-9]\d{2}-\d{3}-\d{4}$
	Telephone string `json:"telephone,omitempty"`

	// transportation office assignments
	TransportationOfficeAssignments []*OfficeUserTransportationOfficeAssignment `json:"transportationOfficeAssignments"`
}

// Validate validates this office user create
func (m *OfficeUserCreate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePrivileges(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTelephone(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOfficeAssignments(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserCreate) validatePrivileges(formats strfmt.Registry) error {
	if swag.IsZero(m.Privileges) { // not required
		return nil
	}

	for i := 0; i < len(m.Privileges); i++ {
		if swag.IsZero(m.Privileges[i]) { // not required
			continue
		}

		if m.Privileges[i] != nil {
			if err := m.Privileges[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("privileges" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("privileges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *OfficeUserCreate) validateRoles(formats strfmt.Registry) error {
	if swag.IsZero(m.Roles) { // not required
		return nil
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
	if swag.IsZero(m.Telephone) { // not required
		return nil
	}

	if err := validate.Pattern("telephone", "body", m.Telephone, `^[2-9]\d{2}-\d{3}-\d{4}$`); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserCreate) validateTransportationOfficeAssignments(formats strfmt.Registry) error {
	if swag.IsZero(m.TransportationOfficeAssignments) { // not required
		return nil
	}

	for i := 0; i < len(m.TransportationOfficeAssignments); i++ {
		if swag.IsZero(m.TransportationOfficeAssignments[i]) { // not required
			continue
		}

		if m.TransportationOfficeAssignments[i] != nil {
			if err := m.TransportationOfficeAssignments[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("transportationOfficeAssignments" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("transportationOfficeAssignments" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this office user create based on the context it is used
func (m *OfficeUserCreate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePrivileges(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRoles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTransportationOfficeAssignments(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserCreate) contextValidatePrivileges(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Privileges); i++ {

		if m.Privileges[i] != nil {

			if swag.IsZero(m.Privileges[i]) { // not required
				return nil
			}

			if err := m.Privileges[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("privileges" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("privileges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

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

func (m *OfficeUserCreate) contextValidateTransportationOfficeAssignments(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.TransportationOfficeAssignments); i++ {

		if m.TransportationOfficeAssignments[i] != nil {

			if swag.IsZero(m.TransportationOfficeAssignments[i]) { // not required
				return nil
			}

			if err := m.TransportationOfficeAssignments[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("transportationOfficeAssignments" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("transportationOfficeAssignments" + "." + strconv.Itoa(i))
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
