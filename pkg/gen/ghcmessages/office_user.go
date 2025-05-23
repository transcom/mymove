// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// OfficeUser office user
//
// swagger:model OfficeUser
type OfficeUser struct {

	// active
	// Required: true
	Active *bool `json:"active"`

	// created at
	// Required: true
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt"`

	// edipi
	// Required: true
	Edipi *string `json:"edipi"`

	// email
	// Required: true
	// Pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	Email *string `json:"email"`

	// first name
	// Required: true
	FirstName *string `json:"firstName"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Format: uuid
	ID *strfmt.UUID `json:"id"`

	// last name
	// Required: true
	LastName *string `json:"lastName"`

	// middle initials
	// Required: true
	MiddleInitials *string `json:"middleInitials"`

	// other unique Id
	// Required: true
	OtherUniqueID *string `json:"otherUniqueId"`

	// rejection reason
	// Required: true
	RejectionReason *string `json:"rejectionReason"`

	// roles
	// Required: true
	Roles []*Role `json:"roles"`

	// status
	// Required: true
	// Enum: [APPROVED REQUESTED REJECTED]
	Status *string `json:"status"`

	// telephone
	// Required: true
	// Pattern: ^[2-9]\d{2}-\d{3}-\d{4}$
	Telephone *string `json:"telephone"`

	// transportation office
	TransportationOffice *TransportationOffice `json:"transportationOffice,omitempty"`

	// transportation office assignments
	TransportationOfficeAssignments []*TransportationOfficeAssignment `json:"transportationOfficeAssignments"`

	// transportation office Id
	// Required: true
	// Format: uuid
	TransportationOfficeID *strfmt.UUID `json:"transportationOfficeId"`

	// updated at
	// Required: true
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt"`

	// user Id
	// Format: uuid
	UserID strfmt.UUID `json:"userId,omitempty"`
}

// Validate validates this office user
func (m *OfficeUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActive(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEdipi(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFirstName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMiddleInitials(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOtherUniqueID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRejectionReason(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTelephone(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOffice(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOfficeAssignments(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOfficeID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUser) validateActive(formats strfmt.Registry) error {

	if err := validate.Required("active", "body", m.Active); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateEdipi(formats strfmt.Registry) error {

	if err := validate.Required("edipi", "body", m.Edipi); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateEmail(formats strfmt.Registry) error {

	if err := validate.Required("email", "body", m.Email); err != nil {
		return err
	}

	if err := validate.Pattern("email", "body", *m.Email, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateFirstName(formats strfmt.Registry) error {

	if err := validate.Required("firstName", "body", m.FirstName); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateLastName(formats strfmt.Registry) error {

	if err := validate.Required("lastName", "body", m.LastName); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateMiddleInitials(formats strfmt.Registry) error {

	if err := validate.Required("middleInitials", "body", m.MiddleInitials); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateOtherUniqueID(formats strfmt.Registry) error {

	if err := validate.Required("otherUniqueId", "body", m.OtherUniqueID); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateRejectionReason(formats strfmt.Registry) error {

	if err := validate.Required("rejectionReason", "body", m.RejectionReason); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateRoles(formats strfmt.Registry) error {

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

var officeUserTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["APPROVED","REQUESTED","REJECTED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		officeUserTypeStatusPropEnum = append(officeUserTypeStatusPropEnum, v)
	}
}

const (

	// OfficeUserStatusAPPROVED captures enum value "APPROVED"
	OfficeUserStatusAPPROVED string = "APPROVED"

	// OfficeUserStatusREQUESTED captures enum value "REQUESTED"
	OfficeUserStatusREQUESTED string = "REQUESTED"

	// OfficeUserStatusREJECTED captures enum value "REJECTED"
	OfficeUserStatusREJECTED string = "REJECTED"
)

// prop value enum
func (m *OfficeUser) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, officeUserTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *OfficeUser) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", *m.Status); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateTelephone(formats strfmt.Registry) error {

	if err := validate.Required("telephone", "body", m.Telephone); err != nil {
		return err
	}

	if err := validate.Pattern("telephone", "body", *m.Telephone, `^[2-9]\d{2}-\d{3}-\d{4}$`); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateTransportationOffice(formats strfmt.Registry) error {
	if swag.IsZero(m.TransportationOffice) { // not required
		return nil
	}

	if m.TransportationOffice != nil {
		if err := m.TransportationOffice.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("transportationOffice")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("transportationOffice")
			}
			return err
		}
	}

	return nil
}

func (m *OfficeUser) validateTransportationOfficeAssignments(formats strfmt.Registry) error {
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

func (m *OfficeUser) validateTransportationOfficeID(formats strfmt.Registry) error {

	if err := validate.Required("transportationOfficeId", "body", m.TransportationOfficeID); err != nil {
		return err
	}

	if err := validate.FormatOf("transportationOfficeId", "body", "uuid", m.TransportationOfficeID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) validateUserID(formats strfmt.Registry) error {
	if swag.IsZero(m.UserID) { // not required
		return nil
	}

	if err := validate.FormatOf("userId", "body", "uuid", m.UserID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this office user based on the context it is used
func (m *OfficeUser) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRoles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTransportationOffice(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTransportationOfficeAssignments(ctx, formats); err != nil {
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

func (m *OfficeUser) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUser) contextValidateRoles(ctx context.Context, formats strfmt.Registry) error {

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

func (m *OfficeUser) contextValidateTransportationOffice(ctx context.Context, formats strfmt.Registry) error {

	if m.TransportationOffice != nil {

		if swag.IsZero(m.TransportationOffice) { // not required
			return nil
		}

		if err := m.TransportationOffice.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("transportationOffice")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("transportationOffice")
			}
			return err
		}
	}

	return nil
}

func (m *OfficeUser) contextValidateTransportationOfficeAssignments(ctx context.Context, formats strfmt.Registry) error {

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

func (m *OfficeUser) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OfficeUser) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OfficeUser) UnmarshalBinary(b []byte) error {
	var res OfficeUser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
