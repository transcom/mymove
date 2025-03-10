// Code generated by go-swagger; DO NOT EDIT.

package evaluation_reports

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetEvaluationReportParams creates a new GetEvaluationReportParams object
//
// There are no default values defined in the spec.
func NewGetEvaluationReportParams() GetEvaluationReportParams {

	return GetEvaluationReportParams{}
}

// GetEvaluationReportParams contains all the bound params for the get evaluation report operation
// typically these are obtained from a http.Request
//
// swagger:parameters getEvaluationReport
type GetEvaluationReportParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*the evaluation report ID to be modified
	  Required: true
	  In: path
	*/
	ReportID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetEvaluationReportParams() beforehand.
func (o *GetEvaluationReportParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rReportID, rhkReportID, _ := route.Params.GetOK("reportID")
	if err := o.bindReportID(rReportID, rhkReportID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindReportID binds and validates parameter ReportID from path.
func (o *GetEvaluationReportParams) bindReportID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("reportID", "path", "strfmt.UUID", raw)
	}
	o.ReportID = *(value.(*strfmt.UUID))

	if err := o.validateReportID(formats); err != nil {
		return err
	}

	return nil
}

// validateReportID carries on validations for parameter ReportID
func (o *GetEvaluationReportParams) validateReportID(formats strfmt.Registry) error {

	if err := validate.FormatOf("reportID", "path", "uuid", o.ReportID.String(), formats); err != nil {
		return err
	}
	return nil
}
