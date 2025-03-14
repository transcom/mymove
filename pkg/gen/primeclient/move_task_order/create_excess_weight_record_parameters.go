// Code generated by go-swagger; DO NOT EDIT.

package move_task_order

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

// NewCreateExcessWeightRecordParams creates a new CreateExcessWeightRecordParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateExcessWeightRecordParams() *CreateExcessWeightRecordParams {
	return &CreateExcessWeightRecordParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateExcessWeightRecordParamsWithTimeout creates a new CreateExcessWeightRecordParams object
// with the ability to set a timeout on a request.
func NewCreateExcessWeightRecordParamsWithTimeout(timeout time.Duration) *CreateExcessWeightRecordParams {
	return &CreateExcessWeightRecordParams{
		timeout: timeout,
	}
}

// NewCreateExcessWeightRecordParamsWithContext creates a new CreateExcessWeightRecordParams object
// with the ability to set a context for a request.
func NewCreateExcessWeightRecordParamsWithContext(ctx context.Context) *CreateExcessWeightRecordParams {
	return &CreateExcessWeightRecordParams{
		Context: ctx,
	}
}

// NewCreateExcessWeightRecordParamsWithHTTPClient creates a new CreateExcessWeightRecordParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateExcessWeightRecordParamsWithHTTPClient(client *http.Client) *CreateExcessWeightRecordParams {
	return &CreateExcessWeightRecordParams{
		HTTPClient: client,
	}
}

/*
CreateExcessWeightRecordParams contains all the parameters to send to the API endpoint

	for the create excess weight record operation.

	Typically these are written to a http.Request.
*/
type CreateExcessWeightRecordParams struct {

	/* File.

	   The file to upload.
	*/
	File runtime.NamedReadCloser

	/* MoveTaskOrderID.

	   UUID of the move being updated.

	   Format: uuid
	*/
	MoveTaskOrderID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create excess weight record params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateExcessWeightRecordParams) WithDefaults() *CreateExcessWeightRecordParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create excess weight record params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateExcessWeightRecordParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create excess weight record params
func (o *CreateExcessWeightRecordParams) WithTimeout(timeout time.Duration) *CreateExcessWeightRecordParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create excess weight record params
func (o *CreateExcessWeightRecordParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create excess weight record params
func (o *CreateExcessWeightRecordParams) WithContext(ctx context.Context) *CreateExcessWeightRecordParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create excess weight record params
func (o *CreateExcessWeightRecordParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create excess weight record params
func (o *CreateExcessWeightRecordParams) WithHTTPClient(client *http.Client) *CreateExcessWeightRecordParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create excess weight record params
func (o *CreateExcessWeightRecordParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFile adds the file to the create excess weight record params
func (o *CreateExcessWeightRecordParams) WithFile(file runtime.NamedReadCloser) *CreateExcessWeightRecordParams {
	o.SetFile(file)
	return o
}

// SetFile adds the file to the create excess weight record params
func (o *CreateExcessWeightRecordParams) SetFile(file runtime.NamedReadCloser) {
	o.File = file
}

// WithMoveTaskOrderID adds the moveTaskOrderID to the create excess weight record params
func (o *CreateExcessWeightRecordParams) WithMoveTaskOrderID(moveTaskOrderID strfmt.UUID) *CreateExcessWeightRecordParams {
	o.SetMoveTaskOrderID(moveTaskOrderID)
	return o
}

// SetMoveTaskOrderID adds the moveTaskOrderId to the create excess weight record params
func (o *CreateExcessWeightRecordParams) SetMoveTaskOrderID(moveTaskOrderID strfmt.UUID) {
	o.MoveTaskOrderID = moveTaskOrderID
}

// WriteToRequest writes these params to a swagger request
func (o *CreateExcessWeightRecordParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	// form file param file
	if err := r.SetFileParam("file", o.File); err != nil {
		return err
	}

	// path param moveTaskOrderID
	if err := r.SetPathParam("moveTaskOrderID", o.MoveTaskOrderID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
