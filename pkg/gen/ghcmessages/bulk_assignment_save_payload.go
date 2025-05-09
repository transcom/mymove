// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BulkAssignmentSavePayload bulk assignment save payload
//
// swagger:model BulkAssignmentSavePayload
type BulkAssignmentSavePayload struct {

	// move data
	MoveData []BulkAssignmentMoveData `json:"moveData"`

	// A string corresponding to the queue type
	// Enum: [COUNSELING CLOSEOUT TASK_ORDER PAYMENT_REQUEST DESTINATION_REQUESTS]
	QueueType string `json:"queueType,omitempty"`

	// user data
	UserData []*BulkAssignmentForUser `json:"userData"`
}

// Validate validates this bulk assignment save payload
func (m *BulkAssignmentSavePayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMoveData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateQueueType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BulkAssignmentSavePayload) validateMoveData(formats strfmt.Registry) error {
	if swag.IsZero(m.MoveData) { // not required
		return nil
	}

	for i := 0; i < len(m.MoveData); i++ {

		if err := m.MoveData[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("moveData" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("moveData" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

var bulkAssignmentSavePayloadTypeQueueTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["COUNSELING","CLOSEOUT","TASK_ORDER","PAYMENT_REQUEST","DESTINATION_REQUESTS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		bulkAssignmentSavePayloadTypeQueueTypePropEnum = append(bulkAssignmentSavePayloadTypeQueueTypePropEnum, v)
	}
}

const (

	// BulkAssignmentSavePayloadQueueTypeCOUNSELING captures enum value "COUNSELING"
	BulkAssignmentSavePayloadQueueTypeCOUNSELING string = "COUNSELING"

	// BulkAssignmentSavePayloadQueueTypeCLOSEOUT captures enum value "CLOSEOUT"
	BulkAssignmentSavePayloadQueueTypeCLOSEOUT string = "CLOSEOUT"

	// BulkAssignmentSavePayloadQueueTypeTASKORDER captures enum value "TASK_ORDER"
	BulkAssignmentSavePayloadQueueTypeTASKORDER string = "TASK_ORDER"

	// BulkAssignmentSavePayloadQueueTypePAYMENTREQUEST captures enum value "PAYMENT_REQUEST"
	BulkAssignmentSavePayloadQueueTypePAYMENTREQUEST string = "PAYMENT_REQUEST"

	// BulkAssignmentSavePayloadQueueTypeDESTINATIONREQUESTS captures enum value "DESTINATION_REQUESTS"
	BulkAssignmentSavePayloadQueueTypeDESTINATIONREQUESTS string = "DESTINATION_REQUESTS"
)

// prop value enum
func (m *BulkAssignmentSavePayload) validateQueueTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, bulkAssignmentSavePayloadTypeQueueTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *BulkAssignmentSavePayload) validateQueueType(formats strfmt.Registry) error {
	if swag.IsZero(m.QueueType) { // not required
		return nil
	}

	// value enum
	if err := m.validateQueueTypeEnum("queueType", "body", m.QueueType); err != nil {
		return err
	}

	return nil
}

func (m *BulkAssignmentSavePayload) validateUserData(formats strfmt.Registry) error {
	if swag.IsZero(m.UserData) { // not required
		return nil
	}

	for i := 0; i < len(m.UserData); i++ {
		if swag.IsZero(m.UserData[i]) { // not required
			continue
		}

		if m.UserData[i] != nil {
			if err := m.UserData[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("userData" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("userData" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this bulk assignment save payload based on the context it is used
func (m *BulkAssignmentSavePayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMoveData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUserData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BulkAssignmentSavePayload) contextValidateMoveData(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.MoveData); i++ {

		if swag.IsZero(m.MoveData[i]) { // not required
			return nil
		}

		if err := m.MoveData[i].ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("moveData" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("moveData" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *BulkAssignmentSavePayload) contextValidateUserData(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.UserData); i++ {

		if m.UserData[i] != nil {

			if swag.IsZero(m.UserData[i]) { // not required
				return nil
			}

			if err := m.UserData[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("userData" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("userData" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BulkAssignmentSavePayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BulkAssignmentSavePayload) UnmarshalBinary(b []byte) error {
	var res BulkAssignmentSavePayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
