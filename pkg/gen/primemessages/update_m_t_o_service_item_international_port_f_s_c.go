// Code generated by go-swagger; DO NOT EDIT.

package primemessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// UpdateMTOServiceItemInternationalPortFSC Subtype used to provide the port for fuel surcharge. This is not creating a new service item but rather updating an existing service item.
//
// swagger:model UpdateMTOServiceItemInternationalPortFSC
type UpdateMTOServiceItemInternationalPortFSC struct {
	idField strfmt.UUID

	// Port used for the shipment. Relevant for moving (PODFSC & POEFSC) service items.
	// Example: PDX
	PortCode *string `json:"portCode"`

	// Service code allowed for this model type.
	// Enum: [PODFSC POEFSC]
	ReServiceCode string `json:"reServiceCode,omitempty"`
}

// ID gets the id of this subtype
func (m *UpdateMTOServiceItemInternationalPortFSC) ID() strfmt.UUID {
	return m.idField
}

// SetID sets the id of this subtype
func (m *UpdateMTOServiceItemInternationalPortFSC) SetID(val strfmt.UUID) {
	m.idField = val
}

// ModelType gets the model type of this subtype
func (m *UpdateMTOServiceItemInternationalPortFSC) ModelType() UpdateMTOServiceItemModelType {
	return "UpdateMTOServiceItemInternationalPortFSC"
}

// SetModelType sets the model type of this subtype
func (m *UpdateMTOServiceItemInternationalPortFSC) SetModelType(val UpdateMTOServiceItemModelType) {
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *UpdateMTOServiceItemInternationalPortFSC) UnmarshalJSON(raw []byte) error {
	var data struct {

		// Port used for the shipment. Relevant for moving (PODFSC & POEFSC) service items.
		// Example: PDX
		PortCode *string `json:"portCode"`

		// Service code allowed for this model type.
		// Enum: [PODFSC POEFSC]
		ReServiceCode string `json:"reServiceCode,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		ID strfmt.UUID `json:"id,omitempty"`

		ModelType UpdateMTOServiceItemModelType `json:"modelType"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result UpdateMTOServiceItemInternationalPortFSC

	result.idField = base.ID

	if base.ModelType != result.ModelType() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid modelType value: %q", base.ModelType)
	}

	result.PortCode = data.PortCode
	result.ReServiceCode = data.ReServiceCode

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m UpdateMTOServiceItemInternationalPortFSC) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// Port used for the shipment. Relevant for moving (PODFSC & POEFSC) service items.
		// Example: PDX
		PortCode *string `json:"portCode"`

		// Service code allowed for this model type.
		// Enum: [PODFSC POEFSC]
		ReServiceCode string `json:"reServiceCode,omitempty"`
	}{

		PortCode: m.PortCode,

		ReServiceCode: m.ReServiceCode,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		ID strfmt.UUID `json:"id,omitempty"`

		ModelType UpdateMTOServiceItemModelType `json:"modelType"`
	}{

		ID: m.ID(),

		ModelType: m.ModelType(),
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this update m t o service item international port f s c
func (m *UpdateMTOServiceItemInternationalPortFSC) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReServiceCode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateMTOServiceItemInternationalPortFSC) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID()) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID().String(), formats); err != nil {
		return err
	}

	return nil
}

var updateMTOServiceItemInternationalPortFSCTypeReServiceCodePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["PODFSC","POEFSC"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		updateMTOServiceItemInternationalPortFSCTypeReServiceCodePropEnum = append(updateMTOServiceItemInternationalPortFSCTypeReServiceCodePropEnum, v)
	}
}

// property enum
func (m *UpdateMTOServiceItemInternationalPortFSC) validateReServiceCodeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, updateMTOServiceItemInternationalPortFSCTypeReServiceCodePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *UpdateMTOServiceItemInternationalPortFSC) validateReServiceCode(formats strfmt.Registry) error {

	if swag.IsZero(m.ReServiceCode) { // not required
		return nil
	}

	// value enum
	if err := m.validateReServiceCodeEnum("reServiceCode", "body", m.ReServiceCode); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this update m t o service item international port f s c based on the context it is used
func (m *UpdateMTOServiceItemInternationalPortFSC) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdateMTOServiceItemInternationalPortFSC) contextValidateModelType(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ModelType().ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("modelType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("modelType")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdateMTOServiceItemInternationalPortFSC) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateMTOServiceItemInternationalPortFSC) UnmarshalBinary(b []byte) error {
	var res UpdateMTOServiceItemInternationalPortFSC
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
