// Code generated by go-swagger; DO NOT EDIT.

package supportmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// MTOAgentType MTO Agent Type
// Example: RELEASING_AGENT
//
// swagger:model MTOAgentType
type MTOAgentType string

func NewMTOAgentType(value MTOAgentType) *MTOAgentType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated MTOAgentType.
func (m MTOAgentType) Pointer() *MTOAgentType {
	return &m
}

const (

	// MTOAgentTypeRELEASINGAGENT captures enum value "RELEASING_AGENT"
	MTOAgentTypeRELEASINGAGENT MTOAgentType = "RELEASING_AGENT"

	// MTOAgentTypeRECEIVINGAGENT captures enum value "RECEIVING_AGENT"
	MTOAgentTypeRECEIVINGAGENT MTOAgentType = "RECEIVING_AGENT"
)

// for schema
var mTOAgentTypeEnum []interface{}

func init() {
	var res []MTOAgentType
	if err := json.Unmarshal([]byte(`["RELEASING_AGENT","RECEIVING_AGENT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mTOAgentTypeEnum = append(mTOAgentTypeEnum, v)
	}
}

func (m MTOAgentType) validateMTOAgentTypeEnum(path, location string, value MTOAgentType) error {
	if err := validate.EnumCase(path, location, value, mTOAgentTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this m t o agent type
func (m MTOAgentType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateMTOAgentTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this m t o agent type based on context it is used
func (m MTOAgentType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
