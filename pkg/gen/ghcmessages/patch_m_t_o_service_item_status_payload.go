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

// PatchMTOServiceItemStatusPayload patch m t o service item status payload
//
// swagger:model PatchMTOServiceItemStatusPayload
type PatchMTOServiceItemStatusPayload struct {

	// Reason the service item was rejected
	// Example: Insufficent details provided
	RejectionReason *string `json:"rejectionReason,omitempty"`

	// Describes all statuses for a MTOServiceItem
	// Enum: [SUBMITTED APPROVED REJECTED]
	Status string `json:"status,omitempty"`
}

// Validate validates this patch m t o service item status payload
func (m *PatchMTOServiceItemStatusPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var patchMTOServiceItemStatusPayloadTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["SUBMITTED","APPROVED","REJECTED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		patchMTOServiceItemStatusPayloadTypeStatusPropEnum = append(patchMTOServiceItemStatusPayloadTypeStatusPropEnum, v)
	}
}

const (

	// PatchMTOServiceItemStatusPayloadStatusSUBMITTED captures enum value "SUBMITTED"
	PatchMTOServiceItemStatusPayloadStatusSUBMITTED string = "SUBMITTED"

	// PatchMTOServiceItemStatusPayloadStatusAPPROVED captures enum value "APPROVED"
	PatchMTOServiceItemStatusPayloadStatusAPPROVED string = "APPROVED"

	// PatchMTOServiceItemStatusPayloadStatusREJECTED captures enum value "REJECTED"
	PatchMTOServiceItemStatusPayloadStatusREJECTED string = "REJECTED"
)

// prop value enum
func (m *PatchMTOServiceItemStatusPayload) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, patchMTOServiceItemStatusPayloadTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *PatchMTOServiceItemStatusPayload) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this patch m t o service item status payload based on context it is used
func (m *PatchMTOServiceItemStatusPayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PatchMTOServiceItemStatusPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PatchMTOServiceItemStatusPayload) UnmarshalBinary(b []byte) error {
	var res PatchMTOServiceItemStatusPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
