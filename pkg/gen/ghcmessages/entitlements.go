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

// Entitlements entitlements
//
// swagger:model Entitlements
type Entitlements struct {

	// Indicates if the move entitlement allows dependents to travel to the new Permanent Duty Station (PDS). This is only present on OCONUS moves.
	// Example: true
	AccompaniedTour *bool `json:"accompaniedTour,omitempty"`

	// authorized weight
	// Example: 2000
	AuthorizedWeight *int64 `json:"authorizedWeight,omitempty"`

	// dependents authorized
	// Example: true
	DependentsAuthorized *bool `json:"dependentsAuthorized,omitempty"`

	// Indicates the number of dependents of the age twelve or older for a move. This is only present on OCONUS moves.
	// Example: 3
	DependentsTwelveAndOver *int64 `json:"dependentsTwelveAndOver,omitempty"`

	// Indicates the number of dependents under the age of twelve for a move. This is only present on OCONUS moves.
	// Example: 5
	DependentsUnderTwelve *int64 `json:"dependentsUnderTwelve,omitempty"`

	// e tag
	ETag string `json:"eTag,omitempty"`

	// gun safe
	// Example: false
	GunSafe bool `json:"gunSafe,omitempty"`

	// id
	// Example: 571008b1-b0de-454d-b843-d71be9f02c04
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// non temporary storage
	// Example: false
	NonTemporaryStorage *bool `json:"nonTemporaryStorage,omitempty"`

	// organizational clothing and individual equipment
	// Example: true
	OrganizationalClothingAndIndividualEquipment bool `json:"organizationalClothingAndIndividualEquipment,omitempty"`

	// privately owned vehicle
	// Example: false
	PrivatelyOwnedVehicle *bool `json:"privatelyOwnedVehicle,omitempty"`

	// pro gear weight
	// Example: 2000
	ProGearWeight int64 `json:"proGearWeight,omitempty"`

	// pro gear weight spouse
	// Example: 500
	ProGearWeightSpouse int64 `json:"proGearWeightSpouse,omitempty"`

	// required medical equipment weight
	// Example: 500
	RequiredMedicalEquipmentWeight int64 `json:"requiredMedicalEquipmentWeight,omitempty"`

	// storage in transit
	// Example: 90
	StorageInTransit *int64 `json:"storageInTransit,omitempty"`

	// total dependents
	// Example: 2
	TotalDependents int64 `json:"totalDependents,omitempty"`

	// total weight
	// Example: 500
	TotalWeight int64 `json:"totalWeight,omitempty"`

	// The amount of weight in pounds that the move is entitled for shipment types of Unaccompanied Baggage.
	// Example: 3
	UbAllowance *int64 `json:"ubAllowance,omitempty"`
}

// Validate validates this entitlements
func (m *Entitlements) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Entitlements) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this entitlements based on context it is used
func (m *Entitlements) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Entitlements) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Entitlements) UnmarshalBinary(b []byte) error {
	var res Entitlements
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}