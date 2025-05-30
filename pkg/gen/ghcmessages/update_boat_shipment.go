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

// UpdateBoatShipment update boat shipment
//
// swagger:model UpdateBoatShipment
type UpdateBoatShipment struct {

	// Does the boat have a trailer
	HasTrailer *bool `json:"hasTrailer,omitempty"`

	// Height of the Boat in inches
	HeightInInches *int64 `json:"heightInInches,omitempty"`

	// Is the trailer roadworthy
	IsRoadworthy *bool `json:"isRoadworthy,omitempty"`

	// Length of the Boat in inches
	LengthInInches *int64 `json:"lengthInInches,omitempty"`

	// Make of the Boat
	Make *string `json:"make,omitempty"`

	// Model of the Boat
	Model *string `json:"model,omitempty"`

	// type
	// Enum: [HAUL_AWAY TOW_AWAY]
	Type *string `json:"type,omitempty"`

	// Width of the Boat in inches
	WidthInInches *int64 `json:"widthInInches,omitempty"`

	// Year of the Boat
	Year *int64 `json:"year,omitempty"`
}

// Validate validates this update boat shipment
func (m *UpdateBoatShipment) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var updateBoatShipmentTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["HAUL_AWAY","TOW_AWAY"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		updateBoatShipmentTypeTypePropEnum = append(updateBoatShipmentTypeTypePropEnum, v)
	}
}

const (

	// UpdateBoatShipmentTypeHAULAWAY captures enum value "HAUL_AWAY"
	UpdateBoatShipmentTypeHAULAWAY string = "HAUL_AWAY"

	// UpdateBoatShipmentTypeTOWAWAY captures enum value "TOW_AWAY"
	UpdateBoatShipmentTypeTOWAWAY string = "TOW_AWAY"
)

// prop value enum
func (m *UpdateBoatShipment) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, updateBoatShipmentTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *UpdateBoatShipment) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this update boat shipment based on context it is used
func (m *UpdateBoatShipment) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UpdateBoatShipment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateBoatShipment) UnmarshalBinary(b []byte) error {
	var res UpdateBoatShipment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
