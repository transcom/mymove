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

// MoveTaskOrder The Move (MoveTaskOrder)
//
// swagger:model MoveTaskOrder
type MoveTaskOrder struct {

	// approved at
	// Format: date-time
	ApprovedAt *strfmt.DateTime `json:"approvedAt,omitempty"`

	// available to prime at
	// Format: date-time
	AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// destination address
	DestinationAddress *Address `json:"destinationAddress,omitempty"`

	// destination duty location
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	DestinationDutyLocation strfmt.UUID `json:"destinationDutyLocation,omitempty"`

	// e tag
	ETag string `json:"eTag,omitempty"`

	// entitlements
	Entitlements *Entitlements `json:"entitlements,omitempty"`

	// id
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// locator
	// Example: 1K43AR
	Locator string `json:"locator,omitempty"`

	// order ID
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	OrderID strfmt.UUID `json:"orderID,omitempty"`

	// origin duty location
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Format: uuid
	OriginDutyLocation strfmt.UUID `json:"originDutyLocation,omitempty"`

	// pickup address
	PickupAddress *Address `json:"pickupAddress,omitempty"`

	// reference Id
	// Example: 1001-3456
	ReferenceID string `json:"referenceId,omitempty"`

	// requested pickup date
	// Format: date
	RequestedPickupDate strfmt.Date `json:"requestedPickupDate,omitempty"`

	// service counseling completed at
	// Format: date-time
	ServiceCounselingCompletedAt *strfmt.DateTime `json:"serviceCounselingCompletedAt,omitempty"`

	// tio remarks
	// Example: approved additional weight
	TioRemarks *string `json:"tioRemarks,omitempty"`

	// updated at
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
}

// Validate validates this move task order
func (m *MoveTaskOrder) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApprovedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAvailableToPrimeAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDestinationAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDestinationDutyLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEntitlements(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOrderID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOriginDutyLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePickupAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRequestedPickupDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateServiceCounselingCompletedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MoveTaskOrder) validateApprovedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ApprovedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("approvedAt", "body", "date-time", m.ApprovedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateAvailableToPrimeAt(formats strfmt.Registry) error {
	if swag.IsZero(m.AvailableToPrimeAt) { // not required
		return nil
	}

	if err := validate.FormatOf("availableToPrimeAt", "body", "date-time", m.AvailableToPrimeAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateDestinationAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.DestinationAddress) { // not required
		return nil
	}

	if m.DestinationAddress != nil {
		if err := m.DestinationAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destinationAddress")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destinationAddress")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) validateDestinationDutyLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.DestinationDutyLocation) { // not required
		return nil
	}

	if err := validate.FormatOf("destinationDutyLocation", "body", "uuid", m.DestinationDutyLocation.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateEntitlements(formats strfmt.Registry) error {
	if swag.IsZero(m.Entitlements) { // not required
		return nil
	}

	if m.Entitlements != nil {
		if err := m.Entitlements.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("entitlements")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("entitlements")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateOrderID(formats strfmt.Registry) error {
	if swag.IsZero(m.OrderID) { // not required
		return nil
	}

	if err := validate.FormatOf("orderID", "body", "uuid", m.OrderID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateOriginDutyLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.OriginDutyLocation) { // not required
		return nil
	}

	if err := validate.FormatOf("originDutyLocation", "body", "uuid", m.OriginDutyLocation.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validatePickupAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.PickupAddress) { // not required
		return nil
	}

	if m.PickupAddress != nil {
		if err := m.PickupAddress.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pickupAddress")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pickupAddress")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) validateRequestedPickupDate(formats strfmt.Registry) error {
	if swag.IsZero(m.RequestedPickupDate) { // not required
		return nil
	}

	if err := validate.FormatOf("requestedPickupDate", "body", "date", m.RequestedPickupDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateServiceCounselingCompletedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ServiceCounselingCompletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("serviceCounselingCompletedAt", "body", "date-time", m.ServiceCounselingCompletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this move task order based on the context it is used
func (m *MoveTaskOrder) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDestinationAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateEntitlements(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePickupAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MoveTaskOrder) contextValidateDestinationAddress(ctx context.Context, formats strfmt.Registry) error {

	if m.DestinationAddress != nil {

		if swag.IsZero(m.DestinationAddress) { // not required
			return nil
		}

		if err := m.DestinationAddress.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destinationAddress")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destinationAddress")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateEntitlements(ctx context.Context, formats strfmt.Registry) error {

	if m.Entitlements != nil {

		if swag.IsZero(m.Entitlements) { // not required
			return nil
		}

		if err := m.Entitlements.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("entitlements")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("entitlements")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) contextValidatePickupAddress(ctx context.Context, formats strfmt.Registry) error {

	if m.PickupAddress != nil {

		if swag.IsZero(m.PickupAddress) { // not required
			return nil
		}

		if err := m.PickupAddress.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pickupAddress")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pickupAddress")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MoveTaskOrder) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MoveTaskOrder) UnmarshalBinary(b []byte) error {
	var res MoveTaskOrder
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}