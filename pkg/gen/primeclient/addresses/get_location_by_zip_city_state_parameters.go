// Code generated by go-swagger; DO NOT EDIT.

package addresses

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

// NewGetLocationByZipCityStateParams creates a new GetLocationByZipCityStateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetLocationByZipCityStateParams() *GetLocationByZipCityStateParams {
	return &GetLocationByZipCityStateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetLocationByZipCityStateParamsWithTimeout creates a new GetLocationByZipCityStateParams object
// with the ability to set a timeout on a request.
func NewGetLocationByZipCityStateParamsWithTimeout(timeout time.Duration) *GetLocationByZipCityStateParams {
	return &GetLocationByZipCityStateParams{
		timeout: timeout,
	}
}

// NewGetLocationByZipCityStateParamsWithContext creates a new GetLocationByZipCityStateParams object
// with the ability to set a context for a request.
func NewGetLocationByZipCityStateParamsWithContext(ctx context.Context) *GetLocationByZipCityStateParams {
	return &GetLocationByZipCityStateParams{
		Context: ctx,
	}
}

// NewGetLocationByZipCityStateParamsWithHTTPClient creates a new GetLocationByZipCityStateParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetLocationByZipCityStateParamsWithHTTPClient(client *http.Client) *GetLocationByZipCityStateParams {
	return &GetLocationByZipCityStateParams{
		HTTPClient: client,
	}
}

/*
GetLocationByZipCityStateParams contains all the parameters to send to the API endpoint

	for the get location by zip city state operation.

	Typically these are written to a http.Request.
*/
type GetLocationByZipCityStateParams struct {

	// Search.
	Search string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get location by zip city state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetLocationByZipCityStateParams) WithDefaults() *GetLocationByZipCityStateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get location by zip city state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetLocationByZipCityStateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) WithTimeout(timeout time.Duration) *GetLocationByZipCityStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) WithContext(ctx context.Context) *GetLocationByZipCityStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) WithHTTPClient(client *http.Client) *GetLocationByZipCityStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSearch adds the search to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) WithSearch(search string) *GetLocationByZipCityStateParams {
	o.SetSearch(search)
	return o
}

// SetSearch adds the search to the get location by zip city state params
func (o *GetLocationByZipCityStateParams) SetSearch(search string) {
	o.Search = search
}

// WriteToRequest writes these params to a swagger request
func (o *GetLocationByZipCityStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param search
	if err := r.SetPathParam("search", o.Search); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
