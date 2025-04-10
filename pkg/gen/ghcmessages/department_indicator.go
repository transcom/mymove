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

// DepartmentIndicator Department indicator
//
// # Military branch of service indicator for orders
//
// swagger:model DepartmentIndicator
type DepartmentIndicator string

func NewDepartmentIndicator(value DepartmentIndicator) *DepartmentIndicator {
	return &value
}

// Pointer returns a pointer to a freshly-allocated DepartmentIndicator.
func (m DepartmentIndicator) Pointer() *DepartmentIndicator {
	return &m
}

const (

	// DepartmentIndicatorARMY captures enum value "ARMY"
	DepartmentIndicatorARMY DepartmentIndicator = "ARMY"

	// DepartmentIndicatorARMYCORPSOFENGINEERS captures enum value "ARMY_CORPS_OF_ENGINEERS"
	DepartmentIndicatorARMYCORPSOFENGINEERS DepartmentIndicator = "ARMY_CORPS_OF_ENGINEERS"

	// DepartmentIndicatorCOASTGUARD captures enum value "COAST_GUARD"
	DepartmentIndicatorCOASTGUARD DepartmentIndicator = "COAST_GUARD"

	// DepartmentIndicatorNAVYANDMARINES captures enum value "NAVY_AND_MARINES"
	DepartmentIndicatorNAVYANDMARINES DepartmentIndicator = "NAVY_AND_MARINES"

	// DepartmentIndicatorAIRANDSPACEFORCE captures enum value "AIR_AND_SPACE_FORCE"
	DepartmentIndicatorAIRANDSPACEFORCE DepartmentIndicator = "AIR_AND_SPACE_FORCE"

	// DepartmentIndicatorOFFICEOFSECRETARYOFDEFENSE captures enum value "OFFICE_OF_SECRETARY_OF_DEFENSE"
	DepartmentIndicatorOFFICEOFSECRETARYOFDEFENSE DepartmentIndicator = "OFFICE_OF_SECRETARY_OF_DEFENSE"
)

// for schema
var departmentIndicatorEnum []interface{}

func init() {
	var res []DepartmentIndicator
	if err := json.Unmarshal([]byte(`["ARMY","ARMY_CORPS_OF_ENGINEERS","COAST_GUARD","NAVY_AND_MARINES","AIR_AND_SPACE_FORCE","OFFICE_OF_SECRETARY_OF_DEFENSE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		departmentIndicatorEnum = append(departmentIndicatorEnum, v)
	}
}

func (m DepartmentIndicator) validateDepartmentIndicatorEnum(path, location string, value DepartmentIndicator) error {
	if err := validate.EnumCase(path, location, value, departmentIndicatorEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this department indicator
func (m DepartmentIndicator) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateDepartmentIndicatorEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this department indicator based on context it is used
func (m DepartmentIndicator) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
