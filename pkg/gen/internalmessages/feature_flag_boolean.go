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

// FeatureFlagBoolean A feature flag
//
// swagger:model FeatureFlagBoolean
type FeatureFlagBoolean struct {

	// entity
	// Example: 11111111-1111-1111-1111-111111111111
	// Required: true
	Entity *string `json:"entity"`

	// key
	// Example: flag
	// Required: true
	Key *string `json:"key"`

	// match
	// Example: true
	// Required: true
	Match *bool `json:"match"`

	// namespace
	// Example: test
	// Required: true
	Namespace *string `json:"namespace"`
}

// Validate validates this feature flag boolean
func (m *FeatureFlagBoolean) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEntity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKey(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMatch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNamespace(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FeatureFlagBoolean) validateEntity(formats strfmt.Registry) error {

	if err := validate.Required("entity", "body", m.Entity); err != nil {
		return err
	}

	return nil
}

func (m *FeatureFlagBoolean) validateKey(formats strfmt.Registry) error {

	if err := validate.Required("key", "body", m.Key); err != nil {
		return err
	}

	return nil
}

func (m *FeatureFlagBoolean) validateMatch(formats strfmt.Registry) error {

	if err := validate.Required("match", "body", m.Match); err != nil {
		return err
	}

	return nil
}

func (m *FeatureFlagBoolean) validateNamespace(formats strfmt.Registry) error {

	if err := validate.Required("namespace", "body", m.Namespace); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this feature flag boolean based on context it is used
func (m *FeatureFlagBoolean) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FeatureFlagBoolean) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FeatureFlagBoolean) UnmarshalBinary(b []byte) error {
	var res FeatureFlagBoolean
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}