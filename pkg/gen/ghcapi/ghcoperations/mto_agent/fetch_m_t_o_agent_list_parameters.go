// Code generated by go-swagger; DO NOT EDIT.

package mto_agent

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewFetchMTOAgentListParams creates a new FetchMTOAgentListParams object
//
// There are no default values defined in the spec.
func NewFetchMTOAgentListParams() FetchMTOAgentListParams {

	return FetchMTOAgentListParams{}
}

// FetchMTOAgentListParams contains all the bound params for the fetch m t o agent list operation
// typically these are obtained from a http.Request
//
// swagger:parameters fetchMTOAgentList
type FetchMTOAgentListParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of move task order
	  Required: true
	  In: path
	*/
	MoveTaskOrderID strfmt.UUID
	/*ID of the shipment
	  Required: true
	  In: path
	*/
	ShipmentID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewFetchMTOAgentListParams() beforehand.
func (o *FetchMTOAgentListParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rMoveTaskOrderID, rhkMoveTaskOrderID, _ := route.Params.GetOK("moveTaskOrderID")
	if err := o.bindMoveTaskOrderID(rMoveTaskOrderID, rhkMoveTaskOrderID, route.Formats); err != nil {
		res = append(res, err)
	}

	rShipmentID, rhkShipmentID, _ := route.Params.GetOK("shipmentID")
	if err := o.bindShipmentID(rShipmentID, rhkShipmentID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindMoveTaskOrderID binds and validates parameter MoveTaskOrderID from path.
func (o *FetchMTOAgentListParams) bindMoveTaskOrderID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("moveTaskOrderID", "path", "strfmt.UUID", raw)
	}
	o.MoveTaskOrderID = *(value.(*strfmt.UUID))

	if err := o.validateMoveTaskOrderID(formats); err != nil {
		return err
	}

	return nil
}

// validateMoveTaskOrderID carries on validations for parameter MoveTaskOrderID
func (o *FetchMTOAgentListParams) validateMoveTaskOrderID(formats strfmt.Registry) error {

	if err := validate.FormatOf("moveTaskOrderID", "path", "uuid", o.MoveTaskOrderID.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindShipmentID binds and validates parameter ShipmentID from path.
func (o *FetchMTOAgentListParams) bindShipmentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("shipmentID", "path", "strfmt.UUID", raw)
	}
	o.ShipmentID = *(value.(*strfmt.UUID))

	if err := o.validateShipmentID(formats); err != nil {
		return err
	}

	return nil
}

// validateShipmentID carries on validations for parameter ShipmentID
func (o *FetchMTOAgentListParams) validateShipmentID(formats strfmt.Registry) error {

	if err := validate.FormatOf("shipmentID", "path", "uuid", o.ShipmentID.String(), formats); err != nil {
		return err
	}
	return nil
}
