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

// User user
//
// swagger:model User
type User struct {

	// active
	// Required: true
	Active *bool `json:"active"`

	// created at
	// Required: true
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt"`

	// current admin session Id
	// Example: WiPgsPj-jPySR1d0dpmvIZ-HvZqemjmaQWxGQ6B8K_w
	// Required: true
	CurrentAdminSessionID *string `json:"currentAdminSessionId"`

	// current mil session Id
	// Example: WiPgsPj-jPySR1d0dpmvIZ-HvZqemjmaQWxGQ6B8K_w
	// Required: true
	CurrentMilSessionID *string `json:"currentMilSessionId"`

	// current office session Id
	// Example: WiPgsPj-jPySR1d0dpmvIZ-HvZqemjmaQWxGQ6B8K_w
	// Required: true
	CurrentOfficeSessionID *string `json:"currentOfficeSessionId"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// okta email
	// Required: true
	// Pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	OktaEmail *string `json:"oktaEmail"`

	// updated at
	// Required: true
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActive(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCurrentAdminSessionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCurrentMilSessionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCurrentOfficeSessionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOktaEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *User) validateActive(formats strfmt.Registry) error {

	if err := validate.Required("active", "body", m.Active); err != nil {
		return err
	}

	return nil
}

func (m *User) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *User) validateCurrentAdminSessionID(formats strfmt.Registry) error {

	if err := validate.Required("currentAdminSessionId", "body", m.CurrentAdminSessionID); err != nil {
		return err
	}

	return nil
}

func (m *User) validateCurrentMilSessionID(formats strfmt.Registry) error {

	if err := validate.Required("currentMilSessionId", "body", m.CurrentMilSessionID); err != nil {
		return err
	}

	return nil
}

func (m *User) validateCurrentOfficeSessionID(formats strfmt.Registry) error {

	if err := validate.Required("currentOfficeSessionId", "body", m.CurrentOfficeSessionID); err != nil {
		return err
	}

	return nil
}

func (m *User) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *User) validateOktaEmail(formats strfmt.Registry) error {

	if err := validate.Required("oktaEmail", "body", m.OktaEmail); err != nil {
		return err
	}

	if err := validate.Pattern("oktaEmail", "body", *m.OktaEmail, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`); err != nil {
		return err
	}

	return nil
}

func (m *User) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this user based on the context it is used
func (m *User) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUpdatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *User) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *User) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *User) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *User) UnmarshalBinary(b []byte) error {
	var res User
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}