// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CreateAppeal Appeal status and remarks left for a violation, created by a GSR user.
//
// swagger:model CreateAppeal
type CreateAppeal struct {

	// The status of the appeal set by the GSR user
	// Example: These are my violation appeal remarks
	// Enum: [sustained rejected]
	AppealStatus string `json:"appealStatus,omitempty"`

	// Remarks left by the GSR user
	// Example: These are my violation appeal remarks
	Remarks string `json:"remarks,omitempty"`
}

// Validate validates this create appeal
func (m *CreateAppeal) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAppealStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var createAppealTypeAppealStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["sustained","rejected"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createAppealTypeAppealStatusPropEnum = append(createAppealTypeAppealStatusPropEnum, v)
	}
}

const (

	// CreateAppealAppealStatusSustained captures enum value "sustained"
	CreateAppealAppealStatusSustained string = "sustained"

	// CreateAppealAppealStatusRejected captures enum value "rejected"
	CreateAppealAppealStatusRejected string = "rejected"
)

// prop value enum
func (m *CreateAppeal) validateAppealStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, createAppealTypeAppealStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *CreateAppeal) validateAppealStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.AppealStatus) { // not required
		return nil
	}

	// value enum
	if err := m.validateAppealStatusEnum("appealStatus", "body", m.AppealStatus); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this create appeal based on context it is used
func (m *CreateAppeal) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreateAppeal) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateAppeal) UnmarshalBinary(b []byte) error {
	var res CreateAppeal
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}