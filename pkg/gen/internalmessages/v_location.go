// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// VLocation A postal code, city, and state lookup
//
// swagger:model VLocation
type VLocation struct {

	// City
	// Example: Anytown
	City string `json:"city,omitempty"`

	// County
	// Example: LOS ANGELES
	County *string `json:"county,omitempty"`

	// ZIP
	// Example: 90210
	// Pattern: ^(\d{5}?)$
	PostalCode string `json:"postalCode,omitempty"`

	// State
	// Enum: [AL AK AR AZ CA CO CT DC DE FL GA HI IA ID IL IN KS KY LA MA MD ME MI MN MO MS MT NC ND NE NH NJ NM NV NY OH OK OR PA RI SC SD TN TX UT VA VT WA WI WV WY]
	State string `json:"state,omitempty"`

	// us post region cities Id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	UsPostRegionCitiesID strfmt.UUID `json:"usPostRegionCitiesId,omitempty"`
}

// Validate validates this v location
func (m *VLocation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePostalCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsPostRegionCitiesID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VLocation) validatePostalCode(formats strfmt.Registry) error {
	if swag.IsZero(m.PostalCode) { // not required
		return nil
	}

	if err := validate.Pattern("postalCode", "body", m.PostalCode, `^(\d{5}?)$`); err != nil {
		return err
	}

	return nil
}

var vLocationTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["AL","AK","AR","AZ","CA","CO","CT","DC","DE","FL","GA","HI","IA","ID","IL","IN","KS","KY","LA","MA","MD","ME","MI","MN","MO","MS","MT","NC","ND","NE","NH","NJ","NM","NV","NY","OH","OK","OR","PA","RI","SC","SD","TN","TX","UT","VA","VT","WA","WI","WV","WY"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		vLocationTypeStatePropEnum = append(vLocationTypeStatePropEnum, v)
	}
}

const (

	// VLocationStateAL captures enum value "AL"
	VLocationStateAL string = "AL"

	// VLocationStateAK captures enum value "AK"
	VLocationStateAK string = "AK"

	// VLocationStateAR captures enum value "AR"
	VLocationStateAR string = "AR"

	// VLocationStateAZ captures enum value "AZ"
	VLocationStateAZ string = "AZ"

	// VLocationStateCA captures enum value "CA"
	VLocationStateCA string = "CA"

	// VLocationStateCO captures enum value "CO"
	VLocationStateCO string = "CO"

	// VLocationStateCT captures enum value "CT"
	VLocationStateCT string = "CT"

	// VLocationStateDC captures enum value "DC"
	VLocationStateDC string = "DC"

	// VLocationStateDE captures enum value "DE"
	VLocationStateDE string = "DE"

	// VLocationStateFL captures enum value "FL"
	VLocationStateFL string = "FL"

	// VLocationStateGA captures enum value "GA"
	VLocationStateGA string = "GA"

	// VLocationStateHI captures enum value "HI"
	VLocationStateHI string = "HI"

	// VLocationStateIA captures enum value "IA"
	VLocationStateIA string = "IA"

	// VLocationStateID captures enum value "ID"
	VLocationStateID string = "ID"

	// VLocationStateIL captures enum value "IL"
	VLocationStateIL string = "IL"

	// VLocationStateIN captures enum value "IN"
	VLocationStateIN string = "IN"

	// VLocationStateKS captures enum value "KS"
	VLocationStateKS string = "KS"

	// VLocationStateKY captures enum value "KY"
	VLocationStateKY string = "KY"

	// VLocationStateLA captures enum value "LA"
	VLocationStateLA string = "LA"

	// VLocationStateMA captures enum value "MA"
	VLocationStateMA string = "MA"

	// VLocationStateMD captures enum value "MD"
	VLocationStateMD string = "MD"

	// VLocationStateME captures enum value "ME"
	VLocationStateME string = "ME"

	// VLocationStateMI captures enum value "MI"
	VLocationStateMI string = "MI"

	// VLocationStateMN captures enum value "MN"
	VLocationStateMN string = "MN"

	// VLocationStateMO captures enum value "MO"
	VLocationStateMO string = "MO"

	// VLocationStateMS captures enum value "MS"
	VLocationStateMS string = "MS"

	// VLocationStateMT captures enum value "MT"
	VLocationStateMT string = "MT"

	// VLocationStateNC captures enum value "NC"
	VLocationStateNC string = "NC"

	// VLocationStateND captures enum value "ND"
	VLocationStateND string = "ND"

	// VLocationStateNE captures enum value "NE"
	VLocationStateNE string = "NE"

	// VLocationStateNH captures enum value "NH"
	VLocationStateNH string = "NH"

	// VLocationStateNJ captures enum value "NJ"
	VLocationStateNJ string = "NJ"

	// VLocationStateNM captures enum value "NM"
	VLocationStateNM string = "NM"

	// VLocationStateNV captures enum value "NV"
	VLocationStateNV string = "NV"

	// VLocationStateNY captures enum value "NY"
	VLocationStateNY string = "NY"

	// VLocationStateOH captures enum value "OH"
	VLocationStateOH string = "OH"

	// VLocationStateOK captures enum value "OK"
	VLocationStateOK string = "OK"

	// VLocationStateOR captures enum value "OR"
	VLocationStateOR string = "OR"

	// VLocationStatePA captures enum value "PA"
	VLocationStatePA string = "PA"

	// VLocationStateRI captures enum value "RI"
	VLocationStateRI string = "RI"

	// VLocationStateSC captures enum value "SC"
	VLocationStateSC string = "SC"

	// VLocationStateSD captures enum value "SD"
	VLocationStateSD string = "SD"

	// VLocationStateTN captures enum value "TN"
	VLocationStateTN string = "TN"

	// VLocationStateTX captures enum value "TX"
	VLocationStateTX string = "TX"

	// VLocationStateUT captures enum value "UT"
	VLocationStateUT string = "UT"

	// VLocationStateVA captures enum value "VA"
	VLocationStateVA string = "VA"

	// VLocationStateVT captures enum value "VT"
	VLocationStateVT string = "VT"

	// VLocationStateWA captures enum value "WA"
	VLocationStateWA string = "WA"

	// VLocationStateWI captures enum value "WI"
	VLocationStateWI string = "WI"

	// VLocationStateWV captures enum value "WV"
	VLocationStateWV string = "WV"

	// VLocationStateWY captures enum value "WY"
	VLocationStateWY string = "WY"
)

// prop value enum
func (m *VLocation) validateStateEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, vLocationTypeStatePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *VLocation) validateState(formats strfmt.Registry) error {
	if swag.IsZero(m.State) { // not required
		return nil
	}

	// value enum
	if err := m.validateStateEnum("state", "body", m.State); err != nil {
		return err
	}

	return nil
}

func (m *VLocation) validateUsPostRegionCitiesID(formats strfmt.Registry) error {
	if swag.IsZero(m.UsPostRegionCitiesID) { // not required
		return nil
	}

	if err := validate.FormatOf("usPostRegionCitiesId", "body", "uuid", m.UsPostRegionCitiesID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this v location based on context it is used
func (m *VLocation) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *VLocation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VLocation) UnmarshalBinary(b []byte) error {
	var res VLocation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}