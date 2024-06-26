// Code generated by go-swagger; DO NOT EDIT.

package duty_locations

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

// NewSearchDutyLocationsParams creates a new SearchDutyLocationsParams object
//
// There are no default values defined in the spec.
func NewSearchDutyLocationsParams() SearchDutyLocationsParams {

	return SearchDutyLocationsParams{}
}

// SearchDutyLocationsParams contains all the bound params for the search duty locations operation
// typically these are obtained from a http.Request
//
// swagger:parameters searchDutyLocations
type SearchDutyLocationsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Search string for duty locations
	  Required: true
	  In: query
	*/
	Search string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSearchDutyLocationsParams() beforehand.
func (o *SearchDutyLocationsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qSearch, qhkSearch, _ := qs.GetOK("search")
	if err := o.bindSearch(qSearch, qhkSearch, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindSearch binds and validates parameter Search from query.
func (o *SearchDutyLocationsParams) bindSearch(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("search", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("search", "query", raw); err != nil {
		return err
	}
	o.Search = raw

	return nil
}
