// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// MovingExpenseType Moving Expense Type
//
// swagger:model MovingExpenseType
type MovingExpenseType string

func NewMovingExpenseType(value MovingExpenseType) *MovingExpenseType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated MovingExpenseType.
func (m MovingExpenseType) Pointer() *MovingExpenseType {
	return &m
}

const (

	// MovingExpenseTypeCONTRACTEDEXPENSE captures enum value "CONTRACTED_EXPENSE"
	MovingExpenseTypeCONTRACTEDEXPENSE MovingExpenseType = "CONTRACTED_EXPENSE"

	// MovingExpenseTypeGAS captures enum value "GAS"
	MovingExpenseTypeGAS MovingExpenseType = "GAS"

	// MovingExpenseTypeOIL captures enum value "OIL"
	MovingExpenseTypeOIL MovingExpenseType = "OIL"

	// MovingExpenseTypeOTHER captures enum value "OTHER"
	MovingExpenseTypeOTHER MovingExpenseType = "OTHER"

	// MovingExpenseTypePACKINGMATERIALS captures enum value "PACKING_MATERIALS"
	MovingExpenseTypePACKINGMATERIALS MovingExpenseType = "PACKING_MATERIALS"

	// MovingExpenseTypeRENTALEQUIPMENT captures enum value "RENTAL_EQUIPMENT"
	MovingExpenseTypeRENTALEQUIPMENT MovingExpenseType = "RENTAL_EQUIPMENT"

	// MovingExpenseTypeSTORAGE captures enum value "STORAGE"
	MovingExpenseTypeSTORAGE MovingExpenseType = "STORAGE"

	// MovingExpenseTypeTOLLS captures enum value "TOLLS"
	MovingExpenseTypeTOLLS MovingExpenseType = "TOLLS"

	// MovingExpenseTypeWEIGHINGFEE captures enum value "WEIGHING_FEE"
	MovingExpenseTypeWEIGHINGFEE MovingExpenseType = "WEIGHING_FEE"
)

// for schema
var movingExpenseTypeEnum []interface{}

func init() {
	var res []MovingExpenseType
	if err := json.Unmarshal([]byte(`["CONTRACTED_EXPENSE","GAS","OIL","OTHER","PACKING_MATERIALS","RENTAL_EQUIPMENT","STORAGE","TOLLS","WEIGHING_FEE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		movingExpenseTypeEnum = append(movingExpenseTypeEnum, v)
	}
}

func (m MovingExpenseType) validateMovingExpenseTypeEnum(path, location string, value MovingExpenseType) error {
	if err := validate.EnumCase(path, location, value, movingExpenseTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this moving expense type
func (m MovingExpenseType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateMovingExpenseTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this moving expense type based on context it is used
func (m MovingExpenseType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}