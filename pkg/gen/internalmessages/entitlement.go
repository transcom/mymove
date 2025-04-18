// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Entitlement entitlement
//
// swagger:model Entitlement
type Entitlement struct {

	// Indicates if the move entitlement allows dependents to travel to the new Permanent Duty Station (PDS). This is only present on OCONUS moves.
	// Example: true
	AccompaniedTour *bool `json:"accompanied_tour,omitempty"`

	// Indicates the number of dependents of the age twelve or older for a move. This is only present on OCONUS moves.
	// Example: 3
	DependentsTwelveAndOver *int64 `json:"dependents_twelve_and_over,omitempty"`

	// Indicates the number of dependents under the age of twelve for a move. This is only present on OCONUS moves.
	// Example: 5
	DependentsUnderTwelve *int64 `json:"dependents_under_twelve,omitempty"`

	// Pro-gear weight limit as set by an Office user, distinct from the service member's default weight allotment determined by pay grade
	//
	// Example: 2000
	ProGear *int64 `json:"proGear,omitempty"`

	// Spouse's pro-gear weight limit as set by an Office user, distinct from the service member's default weight allotment determined by pay grade
	//
	// Example: 500
	ProGearSpouse *int64 `json:"proGearSpouse,omitempty"`

	// The amount of weight in pounds that the move is entitled for shipment types of Unaccompanied Baggage.
	// Example: 3
	UbAllowance *int64 `json:"ub_allowance,omitempty"`

	// Indicates the UB weight restricted to a specific location.
	// Example: 1100
	UbWeightRestriction *int64 `json:"ub_weight_restriction,omitempty"`

	// Indicates the weight restricted to a specific location.
	// Example: 1500
	WeightRestriction *int64 `json:"weight_restriction,omitempty"`
}

// Validate validates this entitlement
func (m *Entitlement) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this entitlement based on context it is used
func (m *Entitlement) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Entitlement) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Entitlement) UnmarshalBinary(b []byte) error {
	var res Entitlement
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
