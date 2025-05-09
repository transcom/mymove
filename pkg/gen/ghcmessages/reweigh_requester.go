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

// ReweighRequester reweigh requester
//
// swagger:model ReweighRequester
type ReweighRequester string

func NewReweighRequester(value ReweighRequester) *ReweighRequester {
	return &value
}

// Pointer returns a pointer to a freshly-allocated ReweighRequester.
func (m ReweighRequester) Pointer() *ReweighRequester {
	return &m
}

const (

	// ReweighRequesterCUSTOMER captures enum value "CUSTOMER"
	ReweighRequesterCUSTOMER ReweighRequester = "CUSTOMER"

	// ReweighRequesterPRIME captures enum value "PRIME"
	ReweighRequesterPRIME ReweighRequester = "PRIME"

	// ReweighRequesterSYSTEM captures enum value "SYSTEM"
	ReweighRequesterSYSTEM ReweighRequester = "SYSTEM"

	// ReweighRequesterTOO captures enum value "TOO"
	ReweighRequesterTOO ReweighRequester = "TOO"
)

// for schema
var reweighRequesterEnum []interface{}

func init() {
	var res []ReweighRequester
	if err := json.Unmarshal([]byte(`["CUSTOMER","PRIME","SYSTEM","TOO"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		reweighRequesterEnum = append(reweighRequesterEnum, v)
	}
}

func (m ReweighRequester) validateReweighRequesterEnum(path, location string, value ReweighRequester) error {
	if err := validate.EnumCase(path, location, value, reweighRequesterEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this reweigh requester
func (m ReweighRequester) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateReweighRequesterEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this reweigh requester based on context it is used
func (m ReweighRequester) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
