// Code generated by go-swagger; DO NOT EDIT.

package mto_shipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// NewUpdateShipmentDestinationAddressParams creates a new UpdateShipmentDestinationAddressParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateShipmentDestinationAddressParams() *UpdateShipmentDestinationAddressParams {
	return &UpdateShipmentDestinationAddressParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateShipmentDestinationAddressParamsWithTimeout creates a new UpdateShipmentDestinationAddressParams object
// with the ability to set a timeout on a request.
func NewUpdateShipmentDestinationAddressParamsWithTimeout(timeout time.Duration) *UpdateShipmentDestinationAddressParams {
	return &UpdateShipmentDestinationAddressParams{
		timeout: timeout,
	}
}

// NewUpdateShipmentDestinationAddressParamsWithContext creates a new UpdateShipmentDestinationAddressParams object
// with the ability to set a context for a request.
func NewUpdateShipmentDestinationAddressParamsWithContext(ctx context.Context) *UpdateShipmentDestinationAddressParams {
	return &UpdateShipmentDestinationAddressParams{
		Context: ctx,
	}
}

// NewUpdateShipmentDestinationAddressParamsWithHTTPClient creates a new UpdateShipmentDestinationAddressParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateShipmentDestinationAddressParamsWithHTTPClient(client *http.Client) *UpdateShipmentDestinationAddressParams {
	return &UpdateShipmentDestinationAddressParams{
		HTTPClient: client,
	}
}

/*
UpdateShipmentDestinationAddressParams contains all the parameters to send to the API endpoint

	for the update shipment destination address operation.

	Typically these are written to a http.Request.
*/
type UpdateShipmentDestinationAddressParams struct {

	/* IfMatch.

	   Needs to be the eTag of the mtoShipment. Optimistic locking is implemented via the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a `412 Precondition Failed` error.

	*/
	IfMatch string

	// Body.
	Body *primemessages.UpdateShipmentDestinationAddress

	/* MtoShipmentID.

	   UUID of the shipment associated with the address

	   Format: uuid
	*/
	MtoShipmentID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update shipment destination address params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateShipmentDestinationAddressParams) WithDefaults() *UpdateShipmentDestinationAddressParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update shipment destination address params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateShipmentDestinationAddressParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithTimeout(timeout time.Duration) *UpdateShipmentDestinationAddressParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithContext(ctx context.Context) *UpdateShipmentDestinationAddressParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithHTTPClient(client *http.Client) *UpdateShipmentDestinationAddressParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIfMatch adds the ifMatch to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithIfMatch(ifMatch string) *UpdateShipmentDestinationAddressParams {
	o.SetIfMatch(ifMatch)
	return o
}

// SetIfMatch adds the ifMatch to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetIfMatch(ifMatch string) {
	o.IfMatch = ifMatch
}

// WithBody adds the body to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithBody(body *primemessages.UpdateShipmentDestinationAddress) *UpdateShipmentDestinationAddressParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetBody(body *primemessages.UpdateShipmentDestinationAddress) {
	o.Body = body
}

// WithMtoShipmentID adds the mtoShipmentID to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) WithMtoShipmentID(mtoShipmentID strfmt.UUID) *UpdateShipmentDestinationAddressParams {
	o.SetMtoShipmentID(mtoShipmentID)
	return o
}

// SetMtoShipmentID adds the mtoShipmentId to the update shipment destination address params
func (o *UpdateShipmentDestinationAddressParams) SetMtoShipmentID(mtoShipmentID strfmt.UUID) {
	o.MtoShipmentID = mtoShipmentID
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateShipmentDestinationAddressParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param If-Match
	if err := r.SetHeaderParam("If-Match", o.IfMatch); err != nil {
		return err
	}
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param mtoShipmentID
	if err := r.SetPathParam("mtoShipmentID", o.MtoShipmentID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}