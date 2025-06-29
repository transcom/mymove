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

// UpdateAllowancePayload update allowance payload
//
// swagger:model UpdateAllowancePayload
type UpdateAllowancePayload struct {

	// Indicates if the move entitlement allows dependents to travel to the new Permanent Duty Station (PDS). This is only present on OCONUS moves.
	// Example: true
	AccompaniedTour *bool `json:"accompaniedTour,omitempty"`

	// agency
	Agency *Affiliation `json:"agency,omitempty"`

	// Indicates the number of dependents of the age twelve or older for a move. This is only present on OCONUS moves.
	// Example: 3
	DependentsTwelveAndOver *int64 `json:"dependentsTwelveAndOver,omitempty"`

	// Indicates the number of dependents under the age of twelve for a move. This is only present on OCONUS moves.
	// Example: 5
	DependentsUnderTwelve *int64 `json:"dependentsUnderTwelve,omitempty"`

	// grade
	Grade *OrderPayGrade `json:"grade,omitempty"`

	// True if user is entitled to move a gun safe (up to 500 lbs) as part of their move without it being charged against their weight allowance.
	GunSafe *bool `json:"gunSafe,omitempty"`

	// unit is in lbs
	// Example: 500
	// Maximum: 500
	// Minimum: 0
	GunSafeWeight *int64 `json:"gunSafeWeight,omitempty"`

	// only for Army
	OrganizationalClothingAndIndividualEquipment *bool `json:"organizationalClothingAndIndividualEquipment,omitempty"`

	// unit is in lbs
	// Example: 2000
	// Maximum: 2000
	// Minimum: 0
	ProGearWeight *int64 `json:"proGearWeight,omitempty"`

	// unit is in lbs
	// Example: 500
	// Maximum: 500
	// Minimum: 0
	ProGearWeightSpouse *int64 `json:"proGearWeightSpouse,omitempty"`

	// unit is in lbs
	// Example: 2000
	// Minimum: 0
	RequiredMedicalEquipmentWeight *int64 `json:"requiredMedicalEquipmentWeight,omitempty"`

	// the number of storage in transit days that the customer is entitled to for a given shipment on their move
	// Minimum: 0
	StorageInTransit *int64 `json:"storageInTransit,omitempty"`

	// ub allowance
	// Example: 500
	UbAllowance *int64 `json:"ubAllowance,omitempty"`

	// Indicates the UB weight restriction for the move to a particular location.
	// Example: 1500
	UbWeightRestriction *int64 `json:"ubWeightRestriction,omitempty"`

	// Indicates the weight restriction for the move to a particular location.
	// Example: 1500
	WeightRestriction *int64 `json:"weightRestriction,omitempty"`
}

// Validate validates this update allowance payload
func (m *UpdateAllowancePayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAgency(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGrade(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGunSafeWeight(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProGearWeight(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProGearWeightSpouse(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRequiredMedicalEquipmentWeight(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStorageInTransit(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateAllowancePayload) validateAgency(formats strfmt.Registry) error {
	if swag.IsZero(m.Agency) { // not required
		return nil
	}

	if m.Agency != nil {
		if err := m.Agency.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("agency")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("agency")
			}
			return err
		}
	}

	return nil
}

func (m *UpdateAllowancePayload) validateGrade(formats strfmt.Registry) error {
	if swag.IsZero(m.Grade) { // not required
		return nil
	}

	if m.Grade != nil {
		if err := m.Grade.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("grade")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("grade")
			}
			return err
		}
	}

	return nil
}

func (m *UpdateAllowancePayload) validateGunSafeWeight(formats strfmt.Registry) error {
	if swag.IsZero(m.GunSafeWeight) { // not required
		return nil
	}

	if err := validate.MinimumInt("gunSafeWeight", "body", *m.GunSafeWeight, 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("gunSafeWeight", "body", *m.GunSafeWeight, 500, false); err != nil {
		return err
	}

	return nil
}

func (m *UpdateAllowancePayload) validateProGearWeight(formats strfmt.Registry) error {
	if swag.IsZero(m.ProGearWeight) { // not required
		return nil
	}

	if err := validate.MinimumInt("proGearWeight", "body", *m.ProGearWeight, 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("proGearWeight", "body", *m.ProGearWeight, 2000, false); err != nil {
		return err
	}

	return nil
}

func (m *UpdateAllowancePayload) validateProGearWeightSpouse(formats strfmt.Registry) error {
	if swag.IsZero(m.ProGearWeightSpouse) { // not required
		return nil
	}

	if err := validate.MinimumInt("proGearWeightSpouse", "body", *m.ProGearWeightSpouse, 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("proGearWeightSpouse", "body", *m.ProGearWeightSpouse, 500, false); err != nil {
		return err
	}

	return nil
}

func (m *UpdateAllowancePayload) validateRequiredMedicalEquipmentWeight(formats strfmt.Registry) error {
	if swag.IsZero(m.RequiredMedicalEquipmentWeight) { // not required
		return nil
	}

	if err := validate.MinimumInt("requiredMedicalEquipmentWeight", "body", *m.RequiredMedicalEquipmentWeight, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *UpdateAllowancePayload) validateStorageInTransit(formats strfmt.Registry) error {
	if swag.IsZero(m.StorageInTransit) { // not required
		return nil
	}

	if err := validate.MinimumInt("storageInTransit", "body", *m.StorageInTransit, 0, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this update allowance payload based on the context it is used
func (m *UpdateAllowancePayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAgency(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateGrade(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateAllowancePayload) contextValidateAgency(ctx context.Context, formats strfmt.Registry) error {

	if m.Agency != nil {

		if swag.IsZero(m.Agency) { // not required
			return nil
		}

		if err := m.Agency.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("agency")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("agency")
			}
			return err
		}
	}

	return nil
}

func (m *UpdateAllowancePayload) contextValidateGrade(ctx context.Context, formats strfmt.Registry) error {

	if m.Grade != nil {

		if swag.IsZero(m.Grade) { // not required
			return nil
		}

		if err := m.Grade.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("grade")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("grade")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdateAllowancePayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateAllowancePayload) UnmarshalBinary(b []byte) error {
	var res UpdateAllowancePayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
