// Code generated by go-swagger; DO NOT EDIT.

package pptasmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// DimensionType Describes a dimension type for a MTOServiceItemDimension.
//
// swagger:model DimensionType
type DimensionType string

func NewDimensionType(value DimensionType) *DimensionType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated DimensionType.
func (m DimensionType) Pointer() *DimensionType {
	return &m
}

const (

	// DimensionTypeITEM captures enum value "ITEM"
	DimensionTypeITEM DimensionType = "ITEM"

	// DimensionTypeCRATE captures enum value "CRATE"
	DimensionTypeCRATE DimensionType = "CRATE"
)

// for schema
var dimensionTypeEnum []interface{}

func init() {
	var res []DimensionType
	if err := json.Unmarshal([]byte(`["ITEM","CRATE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dimensionTypeEnum = append(dimensionTypeEnum, v)
	}
}

func (m DimensionType) validateDimensionTypeEnum(path, location string, value DimensionType) error {
	if err := validate.EnumCase(path, location, value, dimensionTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this dimension type
func (m DimensionType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateDimensionTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this dimension type based on context it is used
func (m DimensionType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
