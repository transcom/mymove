// Code generated by go-swagger; DO NOT EDIT.

package move_task_order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewDownloadMoveOrderParams creates a new DownloadMoveOrderParams object
// with the default values initialized.
func NewDownloadMoveOrderParams() DownloadMoveOrderParams {

	var (
		// initialize parameters with default values

		typeVarDefault = string("ALL")
	)

	return DownloadMoveOrderParams{
		Type: &typeVarDefault,
	}
}

// DownloadMoveOrderParams contains all the bound params for the download move order operation
// typically these are obtained from a http.Request
//
// swagger:parameters downloadMoveOrder
type DownloadMoveOrderParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*the locator code for move order to be downloaded
	  Required: true
	  In: path
	*/
	Locator string
	/*upload type
	  In: query
	  Default: "ALL"
	*/
	Type *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDownloadMoveOrderParams() beforehand.
func (o *DownloadMoveOrderParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rLocator, rhkLocator, _ := route.Params.GetOK("locator")
	if err := o.bindLocator(rLocator, rhkLocator, route.Formats); err != nil {
		res = append(res, err)
	}

	qType, qhkType, _ := qs.GetOK("type")
	if err := o.bindType(qType, qhkType, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLocator binds and validates parameter Locator from path.
func (o *DownloadMoveOrderParams) bindLocator(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.Locator = raw

	return nil
}

// bindType binds and validates parameter Type from query.
func (o *DownloadMoveOrderParams) bindType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewDownloadMoveOrderParams()
		return nil
	}
	o.Type = &raw

	if err := o.validateType(formats); err != nil {
		return err
	}

	return nil
}

// validateType carries on validations for parameter Type
func (o *DownloadMoveOrderParams) validateType(formats strfmt.Registry) error {

	if err := validate.EnumCase("type", "query", *o.Type, []interface{}{"ALL", "ORDERS", "AMENDMENTS"}, true); err != nil {
		return err
	}

	return nil
}