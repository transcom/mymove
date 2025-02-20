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

// SITSummary s i t summary
//
// swagger:model SITSummary
type SITSummary struct {

	// days in s i t
	// Minimum: 0
	DaysInSIT *int64 `json:"daysInSIT,omitempty"`

	// first day s i t service item ID
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	FirstDaySITServiceItemID strfmt.UUID `json:"firstDaySITServiceItemID,omitempty"`

	// location
	// Enum: [ORIGIN DESTINATION]
	Location interface{} `json:"location,omitempty"`

	// sit authorized end date
	// Format: date-time
	SitAuthorizedEndDate strfmt.DateTime `json:"sitAuthorizedEndDate,omitempty"`

	// sit customer contacted
	// Format: date-time
	SitCustomerContacted *strfmt.DateTime `json:"sitCustomerContacted,omitempty"`

	// sit departure date
	// Format: date-time
	SitDepartureDate *strfmt.DateTime `json:"sitDepartureDate,omitempty"`

	// sit entry date
	// Format: date-time
	SitEntryDate strfmt.DateTime `json:"sitEntryDate,omitempty"`

	// sit requested delivery
	// Format: date-time
	SitRequestedDelivery *strfmt.DateTime `json:"sitRequestedDelivery,omitempty"`
}

// Validate validates this s i t summary
func (m *SITSummary) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDaysInSIT(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFirstDaySITServiceItemID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitAuthorizedEndDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitCustomerContacted(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitDepartureDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitEntryDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSitRequestedDelivery(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SITSummary) validateDaysInSIT(formats strfmt.Registry) error {
	if swag.IsZero(m.DaysInSIT) { // not required
		return nil
	}

	if err := validate.MinimumInt("daysInSIT", "body", *m.DaysInSIT, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateFirstDaySITServiceItemID(formats strfmt.Registry) error {
	if swag.IsZero(m.FirstDaySITServiceItemID) { // not required
		return nil
	}

	if err := validate.FormatOf("firstDaySITServiceItemID", "body", "uuid", m.FirstDaySITServiceItemID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateSitAuthorizedEndDate(formats strfmt.Registry) error {
	if swag.IsZero(m.SitAuthorizedEndDate) { // not required
		return nil
	}

	if err := validate.FormatOf("sitAuthorizedEndDate", "body", "date-time", m.SitAuthorizedEndDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateSitCustomerContacted(formats strfmt.Registry) error {
	if swag.IsZero(m.SitCustomerContacted) { // not required
		return nil
	}

	if err := validate.FormatOf("sitCustomerContacted", "body", "date-time", m.SitCustomerContacted.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateSitDepartureDate(formats strfmt.Registry) error {
	if swag.IsZero(m.SitDepartureDate) { // not required
		return nil
	}

	if err := validate.FormatOf("sitDepartureDate", "body", "date-time", m.SitDepartureDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateSitEntryDate(formats strfmt.Registry) error {
	if swag.IsZero(m.SitEntryDate) { // not required
		return nil
	}

	if err := validate.FormatOf("sitEntryDate", "body", "date-time", m.SitEntryDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SITSummary) validateSitRequestedDelivery(formats strfmt.Registry) error {
	if swag.IsZero(m.SitRequestedDelivery) { // not required
		return nil
	}

	if err := validate.FormatOf("sitRequestedDelivery", "body", "date-time", m.SitRequestedDelivery.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this s i t summary based on context it is used
func (m *SITSummary) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SITSummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SITSummary) UnmarshalBinary(b []byte) error {
	var res SITSummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
