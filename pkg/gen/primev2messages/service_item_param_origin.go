// Code generated by go-swagger; DO NOT EDIT.

package primev2messages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ServiceItemParamOrigin service item param origin
//
// swagger:model ServiceItemParamOrigin
type ServiceItemParamOrigin string

func NewServiceItemParamOrigin(value ServiceItemParamOrigin) *ServiceItemParamOrigin {
	return &value
}

// Pointer returns a pointer to a freshly-allocated ServiceItemParamOrigin.
func (m ServiceItemParamOrigin) Pointer() *ServiceItemParamOrigin {
	return &m
}

const (

	// ServiceItemParamOriginPRIME captures enum value "PRIME"
	ServiceItemParamOriginPRIME ServiceItemParamOrigin = "PRIME"

	// ServiceItemParamOriginSYSTEM captures enum value "SYSTEM"
	ServiceItemParamOriginSYSTEM ServiceItemParamOrigin = "SYSTEM"

	// ServiceItemParamOriginPRICER captures enum value "PRICER"
	ServiceItemParamOriginPRICER ServiceItemParamOrigin = "PRICER"

	// ServiceItemParamOriginPAYMENTREQUEST captures enum value "PAYMENT_REQUEST"
	ServiceItemParamOriginPAYMENTREQUEST ServiceItemParamOrigin = "PAYMENT_REQUEST"
)

// for schema
var serviceItemParamOriginEnum []interface{}

func init() {
	var res []ServiceItemParamOrigin
	if err := json.Unmarshal([]byte(`["PRIME","SYSTEM","PRICER","PAYMENT_REQUEST"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		serviceItemParamOriginEnum = append(serviceItemParamOriginEnum, v)
	}
}

func (m ServiceItemParamOrigin) validateServiceItemParamOriginEnum(path, location string, value ServiceItemParamOrigin) error {
	if err := validate.EnumCase(path, location, value, serviceItemParamOriginEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this service item param origin
func (m ServiceItemParamOrigin) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateServiceItemParamOriginEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this service item param origin based on context it is used
func (m ServiceItemParamOrigin) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}