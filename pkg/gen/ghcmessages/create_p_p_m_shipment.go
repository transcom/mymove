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

// CreatePPMShipment A personally procured move is a type of shipment that a service members moves themselves.
//
// swagger:model CreatePPMShipment
type CreatePPMShipment struct {

	// closeout office ID
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	CloseoutOfficeID strfmt.UUID `json:"closeoutOfficeID,omitempty"`

	// destination address
	// Required: true
	DestinationAddress struct {
		PPMDestinationAddress
	} `json:"destinationAddress"`

	// estimated weight
	// Example: 4200
	// Required: true
	EstimatedWeight *int64 `json:"estimatedWeight"`

	// Date the customer expects to move.
	//
	// Required: true
	// Format: date
	ExpectedDepartureDate *strfmt.Date `json:"expectedDepartureDate"`

	// gun safe weight
	GunSafeWeight *int64 `json:"gunSafeWeight,omitempty"`

	// Indicates whether PPM shipment has gun safe.
	//
	// Required: true
	HasGunSafe *bool `json:"hasGunSafe"`

	// Indicates whether PPM shipment has pro-gear.
	//
	// Required: true
	HasProGear *bool `json:"hasProGear"`

	// has secondary destination address
	HasSecondaryDestinationAddress *bool `json:"hasSecondaryDestinationAddress"`

	// has secondary pickup address
	HasSecondaryPickupAddress *bool `json:"hasSecondaryPickupAddress"`

	// has tertiary destination address
	HasTertiaryDestinationAddress *bool `json:"hasTertiaryDestinationAddress"`

	// has tertiary pickup address
	HasTertiaryPickupAddress *bool `json:"hasTertiaryPickupAddress"`

	// Used for PPM shipments only. Denotes if this shipment uses the Actual Expense Reimbursement method.
	// Example: false
	IsActualExpenseReimbursement *bool `json:"isActualExpenseReimbursement"`

	// pickup address
	// Required: true
	PickupAddress struct {
		Address
	} `json:"pickupAddress"`

	// ppm type
	PpmType PPMType `json:"ppmType,omitempty"`

	// pro gear weight
	ProGearWeight *int64 `json:"proGearWeight,omitempty"`

	// secondary destination address
	SecondaryDestinationAddress struct {
		Address
	} `json:"secondaryDestinationAddress,omitempty"`

	// secondary pickup address
	SecondaryPickupAddress struct {
		Address
	} `json:"secondaryPickupAddress,omitempty"`

	// sit estimated departure date
	// Format: date
	SitEstimatedDepartureDate *strfmt.Date `json:"sitEstimatedDepartureDate,omitempty"`

	// sit estimated entry date
	// Format: date
	SitEstimatedEntryDate *strfmt.Date `json:"sitEstimatedEntryDate,omitempty"`

	// sit estimated weight
	// Example: 2000
	SitEstimatedWeight *int64 `json:"sitEstimatedWeight,omitempty"`

	// sit expected
	// Required: true
	SitExpected *bool `json:"sitExpected"`

	// sit location
	SitLocation *SITLocationType `json:"sitLocation,omitempty"`

	// spouse pro gear weight
	SpouseProGearWeight *int64 `json:"spouseProGearWeight,omitempty"`

	// tertiary destination address
	TertiaryDestinationAddress struct {
		Address
	} `json:"tertiaryDestinationAddress,omitempty"`

	// tertiary pickup address
	TertiaryPickupAddress struct {
		Address
	} `json:"tertiaryPickupAddress,omitempty"`
}

// Validate validates this create p p m shipment
func (m *CreatePPMShipment) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCloseoutOfficeID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDestinationAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEstimatedWeight(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExpectedDepartureDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHasGunSafe(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHasProGear(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePickupAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePpmType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSecondaryDestinationAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSecondaryPickupAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitEstimatedDepartureDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitEstimatedEntryDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitExpected(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTertiaryDestinationAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTertiaryPickupAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreatePPMShipment) validateCloseoutOfficeID(formats strfmt.Registry) error {
	if swag.IsZero(m.CloseoutOfficeID) { // not required
		return nil
	}

	if err := validate.FormatOf("closeoutOfficeID", "body", "uuid", m.CloseoutOfficeID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateDestinationAddress(formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) validateEstimatedWeight(formats strfmt.Registry) error {

	if err := validate.Required("estimatedWeight", "body", m.EstimatedWeight); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateExpectedDepartureDate(formats strfmt.Registry) error {

	if err := validate.Required("expectedDepartureDate", "body", m.ExpectedDepartureDate); err != nil {
		return err
	}

	if err := validate.FormatOf("expectedDepartureDate", "body", "date", m.ExpectedDepartureDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateHasGunSafe(formats strfmt.Registry) error {

	if err := validate.Required("hasGunSafe", "body", m.HasGunSafe); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateHasProGear(formats strfmt.Registry) error {

	if err := validate.Required("hasProGear", "body", m.HasProGear); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validatePickupAddress(formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) validatePpmType(formats strfmt.Registry) error {
	if swag.IsZero(m.PpmType) { // not required
		return nil
	}

	if err := m.PpmType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("ppmType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("ppmType")
		}
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateSecondaryDestinationAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.SecondaryDestinationAddress) { // not required
		return nil
	}

	return nil
}

func (m *CreatePPMShipment) validateSecondaryPickupAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.SecondaryPickupAddress) { // not required
		return nil
	}

	return nil
}

func (m *CreatePPMShipment) validateSitEstimatedDepartureDate(formats strfmt.Registry) error {
	if swag.IsZero(m.SitEstimatedDepartureDate) { // not required
		return nil
	}

	if err := validate.FormatOf("sitEstimatedDepartureDate", "body", "date", m.SitEstimatedDepartureDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateSitEstimatedEntryDate(formats strfmt.Registry) error {
	if swag.IsZero(m.SitEstimatedEntryDate) { // not required
		return nil
	}

	if err := validate.FormatOf("sitEstimatedEntryDate", "body", "date", m.SitEstimatedEntryDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateSitExpected(formats strfmt.Registry) error {

	if err := validate.Required("sitExpected", "body", m.SitExpected); err != nil {
		return err
	}

	return nil
}

func (m *CreatePPMShipment) validateSitLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.SitLocation) { // not required
		return nil
	}

	if m.SitLocation != nil {
		if err := m.SitLocation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sitLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sitLocation")
			}
			return err
		}
	}

	return nil
}

func (m *CreatePPMShipment) validateTertiaryDestinationAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.TertiaryDestinationAddress) { // not required
		return nil
	}

	return nil
}

func (m *CreatePPMShipment) validateTertiaryPickupAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.TertiaryPickupAddress) { // not required
		return nil
	}

	return nil
}

// ContextValidate validate this create p p m shipment based on the context it is used
func (m *CreatePPMShipment) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDestinationAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePickupAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePpmType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSecondaryDestinationAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSecondaryPickupAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSitLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTertiaryDestinationAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTertiaryPickupAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreatePPMShipment) contextValidateDestinationAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) contextValidatePickupAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) contextValidatePpmType(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.PpmType) { // not required
		return nil
	}

	if err := m.PpmType.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("ppmType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("ppmType")
		}
		return err
	}

	return nil
}

func (m *CreatePPMShipment) contextValidateSecondaryDestinationAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) contextValidateSecondaryPickupAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) contextValidateSitLocation(ctx context.Context, formats strfmt.Registry) error {

	if m.SitLocation != nil {

		if swag.IsZero(m.SitLocation) { // not required
			return nil
		}

		if err := m.SitLocation.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sitLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sitLocation")
			}
			return err
		}
	}

	return nil
}

func (m *CreatePPMShipment) contextValidateTertiaryDestinationAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *CreatePPMShipment) contextValidateTertiaryPickupAddress(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *CreatePPMShipment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatePPMShipment) UnmarshalBinary(b []byte) error {
	var res CreatePPMShipment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
