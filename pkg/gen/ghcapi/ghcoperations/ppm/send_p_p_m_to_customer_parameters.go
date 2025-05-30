// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewSendPPMToCustomerParams creates a new SendPPMToCustomerParams object
//
// There are no default values defined in the spec.
func NewSendPPMToCustomerParams() SendPPMToCustomerParams {

	return SendPPMToCustomerParams{}
}

// SendPPMToCustomerParams contains all the bound params for the send p p m to customer operation
// typically these are obtained from a http.Request
//
// swagger:parameters sendPPMToCustomer
type SendPPMToCustomerParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: header
	*/
	IfMatch string
	/*UUID of the PPM shipment
	  Required: true
	  In: path
	*/
	PpmShipmentID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSendPPMToCustomerParams() beforehand.
func (o *SendPPMToCustomerParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindIfMatch(r.Header[http.CanonicalHeaderKey("If-Match")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	rPpmShipmentID, rhkPpmShipmentID, _ := route.Params.GetOK("ppmShipmentId")
	if err := o.bindPpmShipmentID(rPpmShipmentID, rhkPpmShipmentID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIfMatch binds and validates parameter IfMatch from header.
func (o *SendPPMToCustomerParams) bindIfMatch(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("If-Match", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("If-Match", "header", raw); err != nil {
		return err
	}
	o.IfMatch = raw

	return nil
}

// bindPpmShipmentID binds and validates parameter PpmShipmentID from path.
func (o *SendPPMToCustomerParams) bindPpmShipmentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("ppmShipmentId", "path", "strfmt.UUID", raw)
	}
	o.PpmShipmentID = *(value.(*strfmt.UUID))

	if err := o.validatePpmShipmentID(formats); err != nil {
		return err
	}

	return nil
}

// validatePpmShipmentID carries on validations for parameter PpmShipmentID
func (o *SendPPMToCustomerParams) validatePpmShipmentID(formats strfmt.Registry) error {

	if err := validate.FormatOf("ppmShipmentId", "path", "uuid", o.PpmShipmentID.String(), formats); err != nil {
		return err
	}
	return nil
}
