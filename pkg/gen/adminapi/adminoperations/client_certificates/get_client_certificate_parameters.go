// Code generated by go-swagger; DO NOT EDIT.

package client_certificates

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetClientCertificateParams creates a new GetClientCertificateParams object
//
// There are no default values defined in the spec.
func NewGetClientCertificateParams() GetClientCertificateParams {

	return GetClientCertificateParams{}
}

// GetClientCertificateParams contains all the bound params for the get client certificate operation
// typically these are obtained from a http.Request
//
// swagger:parameters getClientCertificate
type GetClientCertificateParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ClientCertificateID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetClientCertificateParams() beforehand.
func (o *GetClientCertificateParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rClientCertificateID, rhkClientCertificateID, _ := route.Params.GetOK("clientCertificateId")
	if err := o.bindClientCertificateID(rClientCertificateID, rhkClientCertificateID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindClientCertificateID binds and validates parameter ClientCertificateID from path.
func (o *GetClientCertificateParams) bindClientCertificateID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("clientCertificateId", "path", "strfmt.UUID", raw)
	}
	o.ClientCertificateID = *(value.(*strfmt.UUID))

	if err := o.validateClientCertificateID(formats); err != nil {
		return err
	}

	return nil
}

// validateClientCertificateID carries on validations for parameter ClientCertificateID
func (o *GetClientCertificateParams) validateClientCertificateID(formats strfmt.Registry) error {

	if err := validate.FormatOf("clientCertificateId", "path", "uuid", o.ClientCertificateID.String(), formats); err != nil {
		return err
	}
	return nil
}
