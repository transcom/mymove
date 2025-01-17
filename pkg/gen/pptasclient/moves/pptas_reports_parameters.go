// Code generated by go-swagger; DO NOT EDIT.

package moves

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
)

// NewPptasReportsParams creates a new PptasReportsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPptasReportsParams() *PptasReportsParams {
	return &PptasReportsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPptasReportsParamsWithTimeout creates a new PptasReportsParams object
// with the ability to set a timeout on a request.
func NewPptasReportsParamsWithTimeout(timeout time.Duration) *PptasReportsParams {
	return &PptasReportsParams{
		timeout: timeout,
	}
}

// NewPptasReportsParamsWithContext creates a new PptasReportsParams object
// with the ability to set a context for a request.
func NewPptasReportsParamsWithContext(ctx context.Context) *PptasReportsParams {
	return &PptasReportsParams{
		Context: ctx,
	}
}

// NewPptasReportsParamsWithHTTPClient creates a new PptasReportsParams object
// with the ability to set a custom HTTPClient for a request.
func NewPptasReportsParamsWithHTTPClient(client *http.Client) *PptasReportsParams {
	return &PptasReportsParams{
		HTTPClient: client,
	}
}

/*
PptasReportsParams contains all the parameters to send to the API endpoint

	for the pptas reports operation.

	Typically these are written to a http.Request.
*/
type PptasReportsParams struct {

	/* Since.

	   Only return moves updated since this time. Formatted like "2021-07-23T18:30:47.116Z"

	   Format: date-time
	*/
	Since *strfmt.DateTime

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the pptas reports params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PptasReportsParams) WithDefaults() *PptasReportsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the pptas reports params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PptasReportsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the pptas reports params
func (o *PptasReportsParams) WithTimeout(timeout time.Duration) *PptasReportsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the pptas reports params
func (o *PptasReportsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the pptas reports params
func (o *PptasReportsParams) WithContext(ctx context.Context) *PptasReportsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the pptas reports params
func (o *PptasReportsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the pptas reports params
func (o *PptasReportsParams) WithHTTPClient(client *http.Client) *PptasReportsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the pptas reports params
func (o *PptasReportsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSince adds the since to the pptas reports params
func (o *PptasReportsParams) WithSince(since *strfmt.DateTime) *PptasReportsParams {
	o.SetSince(since)
	return o
}

// SetSince adds the since to the pptas reports params
func (o *PptasReportsParams) SetSince(since *strfmt.DateTime) {
	o.Since = since
}

// WriteToRequest writes these params to a swagger request
func (o *PptasReportsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Since != nil {

		// query param since
		var qrSince strfmt.DateTime

		if o.Since != nil {
			qrSince = *o.Since
		}
		qSince := qrSince.String()
		if qSince != "" {

			if err := r.SetQueryParam("since", qSince); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
