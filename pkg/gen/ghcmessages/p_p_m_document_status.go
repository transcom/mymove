// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// PPMDocumentStatus Status of the PPM document.
//
// swagger:model PPMDocumentStatus
type PPMDocumentStatus string

func NewPPMDocumentStatus(value PPMDocumentStatus) *PPMDocumentStatus {
	return &value
}

// Pointer returns a pointer to a freshly-allocated PPMDocumentStatus.
func (m PPMDocumentStatus) Pointer() *PPMDocumentStatus {
	return &m
}

const (

	// PPMDocumentStatusAPPROVED captures enum value "APPROVED"
	PPMDocumentStatusAPPROVED PPMDocumentStatus = "APPROVED"

	// PPMDocumentStatusEXCLUDED captures enum value "EXCLUDED"
	PPMDocumentStatusEXCLUDED PPMDocumentStatus = "EXCLUDED"

	// PPMDocumentStatusREJECTED captures enum value "REJECTED"
	PPMDocumentStatusREJECTED PPMDocumentStatus = "REJECTED"
)

// for schema
var pPMDocumentStatusEnum []interface{}

func init() {
	var res []PPMDocumentStatus
	if err := json.Unmarshal([]byte(`["APPROVED","EXCLUDED","REJECTED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		pPMDocumentStatusEnum = append(pPMDocumentStatusEnum, v)
	}
}

func (m PPMDocumentStatus) validatePPMDocumentStatusEnum(path, location string, value PPMDocumentStatus) error {
	if err := validate.EnumCase(path, location, value, pPMDocumentStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this p p m document status
func (m PPMDocumentStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validatePPMDocumentStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this p p m document status based on context it is used
func (m PPMDocumentStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}