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

// GSRAppeal An object associating appeals on violations and serious incidents
//
// swagger:model GSRAppeal
type GSRAppeal struct {

	// appeal status
	AppealStatus GSRAppealStatusType `json:"appealStatus,omitempty"`

	// created at
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// id
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// is serious incident
	// Example: false
	IsSeriousIncident bool `json:"isSeriousIncident,omitempty"`

	// office user
	OfficeUser *EvaluationReportOfficeUser `json:"officeUser,omitempty"`

	// office user ID
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	OfficeUserID strfmt.UUID `json:"officeUserID,omitempty"`

	// remarks
	// Example: Office user remarks
	Remarks string `json:"remarks,omitempty"`

	// report ID
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	ReportID strfmt.UUID `json:"reportID,omitempty"`

	// violation ID
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	ViolationID strfmt.UUID `json:"violationID,omitempty"`
}

// Validate validates this g s r appeal
func (m *GSRAppeal) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAppealStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOfficeUser(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOfficeUserID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReportID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateViolationID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GSRAppeal) validateAppealStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.AppealStatus) { // not required
		return nil
	}

	if err := m.AppealStatus.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("appealStatus")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("appealStatus")
		}
		return err
	}

	return nil
}

func (m *GSRAppeal) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GSRAppeal) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GSRAppeal) validateOfficeUser(formats strfmt.Registry) error {
	if swag.IsZero(m.OfficeUser) { // not required
		return nil
	}

	if m.OfficeUser != nil {
		if err := m.OfficeUser.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("officeUser")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("officeUser")
			}
			return err
		}
	}

	return nil
}

func (m *GSRAppeal) validateOfficeUserID(formats strfmt.Registry) error {
	if swag.IsZero(m.OfficeUserID) { // not required
		return nil
	}

	if err := validate.FormatOf("officeUserID", "body", "uuid", m.OfficeUserID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GSRAppeal) validateReportID(formats strfmt.Registry) error {
	if swag.IsZero(m.ReportID) { // not required
		return nil
	}

	if err := validate.FormatOf("reportID", "body", "uuid", m.ReportID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GSRAppeal) validateViolationID(formats strfmt.Registry) error {
	if swag.IsZero(m.ViolationID) { // not required
		return nil
	}

	if err := validate.FormatOf("violationID", "body", "uuid", m.ViolationID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this g s r appeal based on the context it is used
func (m *GSRAppeal) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAppealStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOfficeUser(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GSRAppeal) contextValidateAppealStatus(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.AppealStatus) { // not required
		return nil
	}

	if err := m.AppealStatus.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("appealStatus")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("appealStatus")
		}
		return err
	}

	return nil
}

func (m *GSRAppeal) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *GSRAppeal) contextValidateOfficeUser(ctx context.Context, formats strfmt.Registry) error {

	if m.OfficeUser != nil {

		if swag.IsZero(m.OfficeUser) { // not required
			return nil
		}

		if err := m.OfficeUser.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("officeUser")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("officeUser")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GSRAppeal) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GSRAppeal) UnmarshalBinary(b []byte) error {
	var res GSRAppeal
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}