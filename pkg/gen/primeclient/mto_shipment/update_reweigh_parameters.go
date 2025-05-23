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

// NewUpdateReweighParams creates a new UpdateReweighParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateReweighParams() *UpdateReweighParams {
	return &UpdateReweighParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateReweighParamsWithTimeout creates a new UpdateReweighParams object
// with the ability to set a timeout on a request.
func NewUpdateReweighParamsWithTimeout(timeout time.Duration) *UpdateReweighParams {
	return &UpdateReweighParams{
		timeout: timeout,
	}
}

// NewUpdateReweighParamsWithContext creates a new UpdateReweighParams object
// with the ability to set a context for a request.
func NewUpdateReweighParamsWithContext(ctx context.Context) *UpdateReweighParams {
	return &UpdateReweighParams{
		Context: ctx,
	}
}

// NewUpdateReweighParamsWithHTTPClient creates a new UpdateReweighParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateReweighParamsWithHTTPClient(client *http.Client) *UpdateReweighParams {
	return &UpdateReweighParams{
		HTTPClient: client,
	}
}

/*
UpdateReweighParams contains all the parameters to send to the API endpoint

	for the update reweigh operation.

	Typically these are written to a http.Request.
*/
type UpdateReweighParams struct {

	/* IfMatch.

	   Optimistic locking is implemented via the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a `412 Precondition Failed` error.

	*/
	IfMatch string

	// Body.
	Body *primemessages.UpdateReweigh

	/* MtoShipmentID.

	   UUID of the shipment associated with the reweigh

	   Format: uuid
	*/
	MtoShipmentID strfmt.UUID

	/* ReweighID.

	   UUID of the reweigh being updated

	   Format: uuid
	*/
	ReweighID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update reweigh params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateReweighParams) WithDefaults() *UpdateReweighParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update reweigh params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateReweighParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update reweigh params
func (o *UpdateReweighParams) WithTimeout(timeout time.Duration) *UpdateReweighParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update reweigh params
func (o *UpdateReweighParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update reweigh params
func (o *UpdateReweighParams) WithContext(ctx context.Context) *UpdateReweighParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update reweigh params
func (o *UpdateReweighParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update reweigh params
func (o *UpdateReweighParams) WithHTTPClient(client *http.Client) *UpdateReweighParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update reweigh params
func (o *UpdateReweighParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIfMatch adds the ifMatch to the update reweigh params
func (o *UpdateReweighParams) WithIfMatch(ifMatch string) *UpdateReweighParams {
	o.SetIfMatch(ifMatch)
	return o
}

// SetIfMatch adds the ifMatch to the update reweigh params
func (o *UpdateReweighParams) SetIfMatch(ifMatch string) {
	o.IfMatch = ifMatch
}

// WithBody adds the body to the update reweigh params
func (o *UpdateReweighParams) WithBody(body *primemessages.UpdateReweigh) *UpdateReweighParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update reweigh params
func (o *UpdateReweighParams) SetBody(body *primemessages.UpdateReweigh) {
	o.Body = body
}

// WithMtoShipmentID adds the mtoShipmentID to the update reweigh params
func (o *UpdateReweighParams) WithMtoShipmentID(mtoShipmentID strfmt.UUID) *UpdateReweighParams {
	o.SetMtoShipmentID(mtoShipmentID)
	return o
}

// SetMtoShipmentID adds the mtoShipmentId to the update reweigh params
func (o *UpdateReweighParams) SetMtoShipmentID(mtoShipmentID strfmt.UUID) {
	o.MtoShipmentID = mtoShipmentID
}

// WithReweighID adds the reweighID to the update reweigh params
func (o *UpdateReweighParams) WithReweighID(reweighID strfmt.UUID) *UpdateReweighParams {
	o.SetReweighID(reweighID)
	return o
}

// SetReweighID adds the reweighId to the update reweigh params
func (o *UpdateReweighParams) SetReweighID(reweighID strfmt.UUID) {
	o.ReweighID = reweighID
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateReweighParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param reweighID
	if err := r.SetPathParam("reweighID", o.ReweighID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
