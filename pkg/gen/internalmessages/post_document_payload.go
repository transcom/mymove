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

// PostDocumentPayload post document payload
//
// swagger:model PostDocumentPayload
type PostDocumentPayload struct {

	// The service member this document belongs to
	// Format: uuid
	ServiceMemberID strfmt.UUID `json:"service_member_id,omitempty"`
}

// Validate validates this post document payload
func (m *PostDocumentPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateServiceMemberID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostDocumentPayload) validateServiceMemberID(formats strfmt.Registry) error {
	if swag.IsZero(m.ServiceMemberID) { // not required
		return nil
	}

	if err := validate.FormatOf("service_member_id", "body", "uuid", m.ServiceMemberID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this post document payload based on context it is used
func (m *PostDocumentPayload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PostDocumentPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostDocumentPayload) UnmarshalBinary(b []byte) error {
	var res PostDocumentPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}