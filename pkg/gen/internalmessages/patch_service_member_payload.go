// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PatchServiceMemberPayload patch service member payload
//
// swagger:model PatchServiceMemberPayload
type PatchServiceMemberPayload struct {

	// affiliation
	Affiliation *Affiliation `json:"affiliation,omitempty"`

	// backup mailing address
	BackupMailingAddress *Address `json:"backup_mailing_address,omitempty"`

	// current location id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	CurrentLocationID *strfmt.UUID `json:"current_location_id,omitempty"`

	// DoD ID number
	// Example: 5789345789
	// Max Length: 10
	// Min Length: 10
	// Pattern: ^\d{10}$
	Edipi *string `json:"edipi,omitempty"`

	// Email
	EmailIsPreferred *bool `json:"email_is_preferred,omitempty"`

	// USCG EMPLID
	// Example: 5789345
	// Max Length: 7
	// Min Length: 7
	// Pattern: ^\d{7}$
	Emplid *string `json:"emplid,omitempty"`

	// First name
	// Example: John
	FirstName *string `json:"first_name,omitempty"`

	// Last name
	// Example: Donut
	LastName *string `json:"last_name,omitempty"`

	// Middle name
	// Example: L.
	MiddleName *string `json:"middle_name,omitempty"`

	// Personal Email
	// Example: john_bob@example.com
	// Pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	PersonalEmail *string `json:"personal_email,omitempty"`

	// Phone
	PhoneIsPreferred *bool `json:"phone_is_preferred,omitempty"`

	// residential address
	ResidentialAddress *Address `json:"residential_address,omitempty"`

	// Alternate Phone
	// Example: 212-555-5555
	// Pattern: ^([2-9]\d{2}-\d{3}-\d{4})?$
	SecondaryTelephone *string `json:"secondary_telephone,omitempty"`

	// Suffix
	// Example: Jr.
	Suffix *string `json:"suffix,omitempty"`

	// Best Contact Phone
	// Example: 212-555-5555
	// Pattern: ^[2-9]\d{2}-\d{3}-\d{4}$
	Telephone *string `json:"telephone,omitempty"`

	// user id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	UserID strfmt.UUID `json:"user_id,omitempty"`
}

// Validate validates this patch service member payload
func (m *PatchServiceMemberPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAffiliation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBackupMailingAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCurrentLocationID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEdipi(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEmplid(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePersonalEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResidentialAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSecondaryTelephone(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTelephone(formats); err != nil {
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

func (m *PatchServiceMemberPayload) validateAffiliation(formats strfmt.Registry) error {
	if swag.IsZero(m.Affiliation) { // not required
		return nil
	}

	if m.Affiliation != nil {
		if err := m.Affiliation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("affiliation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("affiliation")
			}
			return err
		}
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateBackupMailingAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.BackupMailingAddress) { // not required
		return nil
	}

	if m.BackupMailingAddress != nil {
		if err := m.BackupMailingAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("backup_mailing_address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("backup_mailing_address")
			}
			return err
		}
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateCurrentLocationID(formats strfmt.Registry) error {
	if swag.IsZero(m.CurrentLocationID) { // not required
		return nil
	}

	if err := validate.FormatOf("current_location_id", "body", "uuid", m.CurrentLocationID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateEdipi(formats strfmt.Registry) error {
	if swag.IsZero(m.Edipi) { // not required
		return nil
	}

	if err := validate.MinLength("edipi", "body", *m.Edipi, 10); err != nil {
		return err
	}

	if err := validate.MaxLength("edipi", "body", *m.Edipi, 10); err != nil {
		return err
	}

	if err := validate.Pattern("edipi", "body", *m.Edipi, `^\d{10}$`); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateEmplid(formats strfmt.Registry) error {
	if swag.IsZero(m.Emplid) { // not required
		return nil
	}

	if err := validate.MinLength("emplid", "body", *m.Emplid, 7); err != nil {
		return err
	}

	if err := validate.MaxLength("emplid", "body", *m.Emplid, 7); err != nil {
		return err
	}

	if err := validate.Pattern("emplid", "body", *m.Emplid, `^\d{7}$`); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validatePersonalEmail(formats strfmt.Registry) error {
	if swag.IsZero(m.PersonalEmail) { // not required
		return nil
	}

	if err := validate.Pattern("personal_email", "body", *m.PersonalEmail, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateResidentialAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.ResidentialAddress) { // not required
		return nil
	}

	if m.ResidentialAddress != nil {
		if err := m.ResidentialAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("residential_address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("residential_address")
			}
			return err
		}
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateSecondaryTelephone(formats strfmt.Registry) error {
	if swag.IsZero(m.SecondaryTelephone) { // not required
		return nil
	}

	if err := validate.Pattern("secondary_telephone", "body", *m.SecondaryTelephone, `^([2-9]\d{2}-\d{3}-\d{4})?$`); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateTelephone(formats strfmt.Registry) error {
	if swag.IsZero(m.Telephone) { // not required
		return nil
	}

	if err := validate.Pattern("telephone", "body", *m.Telephone, `^[2-9]\d{2}-\d{3}-\d{4}$`); err != nil {
		return err
	}

	return nil
}

func (m *PatchServiceMemberPayload) validateUserID(formats strfmt.Registry) error {
	if swag.IsZero(m.UserID) { // not required
		return nil
	}

	if err := validate.FormatOf("user_id", "body", "uuid", m.UserID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this patch service member payload based on the context it is used
func (m *PatchServiceMemberPayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAffiliation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateBackupMailingAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResidentialAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PatchServiceMemberPayload) contextValidateAffiliation(ctx context.Context, formats strfmt.Registry) error {

	if m.Affiliation != nil {

		if swag.IsZero(m.Affiliation) { // not required
			return nil
		}

		if err := m.Affiliation.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("affiliation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("affiliation")
			}
			return err
		}
	}

	return nil
}

func (m *PatchServiceMemberPayload) contextValidateBackupMailingAddress(ctx context.Context, formats strfmt.Registry) error {

	if m.BackupMailingAddress != nil {

		if swag.IsZero(m.BackupMailingAddress) { // not required
			return nil
		}

		if err := m.BackupMailingAddress.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("backup_mailing_address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("backup_mailing_address")
			}
			return err
		}
	}

	return nil
}

func (m *PatchServiceMemberPayload) contextValidateResidentialAddress(ctx context.Context, formats strfmt.Registry) error {

	if m.ResidentialAddress != nil {

		if swag.IsZero(m.ResidentialAddress) { // not required
			return nil
		}

		if err := m.ResidentialAddress.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("residential_address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("residential_address")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PatchServiceMemberPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PatchServiceMemberPayload) UnmarshalBinary(b []byte) error {
	var res PatchServiceMemberPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}