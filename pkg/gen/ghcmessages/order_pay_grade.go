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

// OrderPayGrade Grade
//
// swagger:model OrderPayGrade
type OrderPayGrade string

func NewOrderPayGrade(value OrderPayGrade) *OrderPayGrade {
	return &value
}

// Pointer returns a pointer to a freshly-allocated OrderPayGrade.
func (m OrderPayGrade) Pointer() *OrderPayGrade {
	return &m
}

const (

	// OrderPayGradeEDash1 captures enum value "E-1"
	OrderPayGradeEDash1 OrderPayGrade = "E-1"

	// OrderPayGradeEDash2 captures enum value "E-2"
	OrderPayGradeEDash2 OrderPayGrade = "E-2"

	// OrderPayGradeEDash3 captures enum value "E-3"
	OrderPayGradeEDash3 OrderPayGrade = "E-3"

	// OrderPayGradeEDash4 captures enum value "E-4"
	OrderPayGradeEDash4 OrderPayGrade = "E-4"

	// OrderPayGradeEDash5 captures enum value "E-5"
	OrderPayGradeEDash5 OrderPayGrade = "E-5"

	// OrderPayGradeEDash6 captures enum value "E-6"
	OrderPayGradeEDash6 OrderPayGrade = "E-6"

	// OrderPayGradeEDash7 captures enum value "E-7"
	OrderPayGradeEDash7 OrderPayGrade = "E-7"

	// OrderPayGradeEDash8 captures enum value "E-8"
	OrderPayGradeEDash8 OrderPayGrade = "E-8"

	// OrderPayGradeEDash9 captures enum value "E-9"
	OrderPayGradeEDash9 OrderPayGrade = "E-9"

	// OrderPayGradeEDash9DashSPECIALDashSENIORDashENLISTED captures enum value "E-9-SPECIAL-SENIOR-ENLISTED"
	OrderPayGradeEDash9DashSPECIALDashSENIORDashENLISTED OrderPayGrade = "E-9-SPECIAL-SENIOR-ENLISTED"

	// OrderPayGradeODash1 captures enum value "O-1"
	OrderPayGradeODash1 OrderPayGrade = "O-1"

	// OrderPayGradeODash2 captures enum value "O-2"
	OrderPayGradeODash2 OrderPayGrade = "O-2"

	// OrderPayGradeODash3 captures enum value "O-3"
	OrderPayGradeODash3 OrderPayGrade = "O-3"

	// OrderPayGradeODash4 captures enum value "O-4"
	OrderPayGradeODash4 OrderPayGrade = "O-4"

	// OrderPayGradeODash5 captures enum value "O-5"
	OrderPayGradeODash5 OrderPayGrade = "O-5"

	// OrderPayGradeODash6 captures enum value "O-6"
	OrderPayGradeODash6 OrderPayGrade = "O-6"

	// OrderPayGradeODash7 captures enum value "O-7"
	OrderPayGradeODash7 OrderPayGrade = "O-7"

	// OrderPayGradeODash8 captures enum value "O-8"
	OrderPayGradeODash8 OrderPayGrade = "O-8"

	// OrderPayGradeODash9 captures enum value "O-9"
	OrderPayGradeODash9 OrderPayGrade = "O-9"

	// OrderPayGradeODash10 captures enum value "O-10"
	OrderPayGradeODash10 OrderPayGrade = "O-10"

	// OrderPayGradeWDash1 captures enum value "W-1"
	OrderPayGradeWDash1 OrderPayGrade = "W-1"

	// OrderPayGradeWDash2 captures enum value "W-2"
	OrderPayGradeWDash2 OrderPayGrade = "W-2"

	// OrderPayGradeWDash3 captures enum value "W-3"
	OrderPayGradeWDash3 OrderPayGrade = "W-3"

	// OrderPayGradeWDash4 captures enum value "W-4"
	OrderPayGradeWDash4 OrderPayGrade = "W-4"

	// OrderPayGradeWDash5 captures enum value "W-5"
	OrderPayGradeWDash5 OrderPayGrade = "W-5"

	// OrderPayGradeAVIATIONCADET captures enum value "AVIATION_CADET"
	OrderPayGradeAVIATIONCADET OrderPayGrade = "AVIATION_CADET"

	// OrderPayGradeCIVILIANEMPLOYEE captures enum value "CIVILIAN_EMPLOYEE"
	OrderPayGradeCIVILIANEMPLOYEE OrderPayGrade = "CIVILIAN_EMPLOYEE"

	// OrderPayGradeACADEMYCADET captures enum value "ACADEMY_CADET"
	OrderPayGradeACADEMYCADET OrderPayGrade = "ACADEMY_CADET"

	// OrderPayGradeMIDSHIPMAN captures enum value "MIDSHIPMAN"
	OrderPayGradeMIDSHIPMAN OrderPayGrade = "MIDSHIPMAN"
)

// for schema
var orderPayGradeEnum []interface{}

func init() {
	var res []OrderPayGrade
	if err := json.Unmarshal([]byte(`["E-1","E-2","E-3","E-4","E-5","E-6","E-7","E-8","E-9","E-9-SPECIAL-SENIOR-ENLISTED","O-1","O-2","O-3","O-4","O-5","O-6","O-7","O-8","O-9","O-10","W-1","W-2","W-3","W-4","W-5","AVIATION_CADET","CIVILIAN_EMPLOYEE","ACADEMY_CADET","MIDSHIPMAN"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		orderPayGradeEnum = append(orderPayGradeEnum, v)
	}
}

func (m OrderPayGrade) validateOrderPayGradeEnum(path, location string, value OrderPayGrade) error {
	if err := validate.EnumCase(path, location, value, orderPayGradeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this order pay grade
func (m OrderPayGrade) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateOrderPayGradeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this order pay grade based on context it is used
func (m OrderPayGrade) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
