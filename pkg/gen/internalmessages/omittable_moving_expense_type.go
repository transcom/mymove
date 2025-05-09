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

// OmittableMovingExpenseType Moving Expense Type
//
// swagger:model OmittableMovingExpenseType
type OmittableMovingExpenseType string

func NewOmittableMovingExpenseType(value OmittableMovingExpenseType) *OmittableMovingExpenseType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated OmittableMovingExpenseType.
func (m OmittableMovingExpenseType) Pointer() *OmittableMovingExpenseType {
	return &m
}

const (

	// OmittableMovingExpenseTypeCONTRACTEDEXPENSE captures enum value "CONTRACTED_EXPENSE"
	OmittableMovingExpenseTypeCONTRACTEDEXPENSE OmittableMovingExpenseType = "CONTRACTED_EXPENSE"

	// OmittableMovingExpenseTypeGAS captures enum value "GAS"
	OmittableMovingExpenseTypeGAS OmittableMovingExpenseType = "GAS"

	// OmittableMovingExpenseTypeOIL captures enum value "OIL"
	OmittableMovingExpenseTypeOIL OmittableMovingExpenseType = "OIL"

	// OmittableMovingExpenseTypeOTHER captures enum value "OTHER"
	OmittableMovingExpenseTypeOTHER OmittableMovingExpenseType = "OTHER"

	// OmittableMovingExpenseTypePACKINGMATERIALS captures enum value "PACKING_MATERIALS"
	OmittableMovingExpenseTypePACKINGMATERIALS OmittableMovingExpenseType = "PACKING_MATERIALS"

	// OmittableMovingExpenseTypeRENTALEQUIPMENT captures enum value "RENTAL_EQUIPMENT"
	OmittableMovingExpenseTypeRENTALEQUIPMENT OmittableMovingExpenseType = "RENTAL_EQUIPMENT"

	// OmittableMovingExpenseTypeSTORAGE captures enum value "STORAGE"
	OmittableMovingExpenseTypeSTORAGE OmittableMovingExpenseType = "STORAGE"

	// OmittableMovingExpenseTypeTOLLS captures enum value "TOLLS"
	OmittableMovingExpenseTypeTOLLS OmittableMovingExpenseType = "TOLLS"

	// OmittableMovingExpenseTypeWEIGHINGFEE captures enum value "WEIGHING_FEE"
	OmittableMovingExpenseTypeWEIGHINGFEE OmittableMovingExpenseType = "WEIGHING_FEE"

	// OmittableMovingExpenseTypeSMALLPACKAGE captures enum value "SMALL_PACKAGE"
	OmittableMovingExpenseTypeSMALLPACKAGE OmittableMovingExpenseType = "SMALL_PACKAGE"
)

// for schema
var omittableMovingExpenseTypeEnum []interface{}

func init() {
	var res []OmittableMovingExpenseType
	if err := json.Unmarshal([]byte(`["CONTRACTED_EXPENSE","GAS","OIL","OTHER","PACKING_MATERIALS","RENTAL_EQUIPMENT","STORAGE","TOLLS","WEIGHING_FEE","SMALL_PACKAGE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		omittableMovingExpenseTypeEnum = append(omittableMovingExpenseTypeEnum, v)
	}
}

func (m OmittableMovingExpenseType) validateOmittableMovingExpenseTypeEnum(path, location string, value OmittableMovingExpenseType) error {
	if err := validate.EnumCase(path, location, value, omittableMovingExpenseTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this omittable moving expense type
func (m OmittableMovingExpenseType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateOmittableMovingExpenseTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this omittable moving expense type based on context it is used
func (m OmittableMovingExpenseType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
