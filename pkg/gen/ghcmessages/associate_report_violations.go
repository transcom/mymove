// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AssociateReportViolations A list of PWS violation string ids to associate with an evaluation report
//
// swagger:model AssociateReportViolations
type AssociateReportViolations struct {

	// violations
	Violations []strfmt.UUID `json:"violations"`
}

// Validate validates this associate report violations
func (m *AssociateReportViolations) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateViolations(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AssociateReportViolations) validateViolations(formats strfmt.Registry) error {
	if swag.IsZero(m.Violations) { // not required
		return nil
	}

	for i := 0; i < len(m.Violations); i++ {

		if err := validate.FormatOf("violations"+"."+strconv.Itoa(i), "body", "uuid", m.Violations[i].String(), formats); err != nil {
			return err
		}

	}

	return nil
}

// ContextValidate validates this associate report violations based on context it is used
func (m *AssociateReportViolations) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AssociateReportViolations) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AssociateReportViolations) UnmarshalBinary(b []byte) error {
	var res AssociateReportViolations
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
