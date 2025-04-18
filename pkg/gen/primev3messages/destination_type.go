// Code generated by go-swagger; DO NOT EDIT.

package primev3messages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// DestinationType Destination Type
// Example: OTHER_THAN_AUTHORIZED
//
// swagger:model DestinationType
type DestinationType string

func NewDestinationType(value DestinationType) *DestinationType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated DestinationType.
func (m DestinationType) Pointer() *DestinationType {
	return &m
}

const (

	// DestinationTypeHOMEOFRECORD captures enum value "HOME_OF_RECORD"
	DestinationTypeHOMEOFRECORD DestinationType = "HOME_OF_RECORD"

	// DestinationTypeHOMEOFSELECTION captures enum value "HOME_OF_SELECTION"
	DestinationTypeHOMEOFSELECTION DestinationType = "HOME_OF_SELECTION"

	// DestinationTypePLACEENTEREDACTIVEDUTY captures enum value "PLACE_ENTERED_ACTIVE_DUTY"
	DestinationTypePLACEENTEREDACTIVEDUTY DestinationType = "PLACE_ENTERED_ACTIVE_DUTY"

	// DestinationTypeOTHERTHANAUTHORIZED captures enum value "OTHER_THAN_AUTHORIZED"
	DestinationTypeOTHERTHANAUTHORIZED DestinationType = "OTHER_THAN_AUTHORIZED"
)

// for schema
var destinationTypeEnum []interface{}

func init() {
	var res []DestinationType
	if err := json.Unmarshal([]byte(`["HOME_OF_RECORD","HOME_OF_SELECTION","PLACE_ENTERED_ACTIVE_DUTY","OTHER_THAN_AUTHORIZED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		destinationTypeEnum = append(destinationTypeEnum, v)
	}
}

func (m DestinationType) validateDestinationTypeEnum(path, location string, value DestinationType) error {
	if err := validate.EnumCase(path, location, value, destinationTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this destination type
func (m DestinationType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateDestinationTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this destination type based on context it is used
func (m DestinationType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
