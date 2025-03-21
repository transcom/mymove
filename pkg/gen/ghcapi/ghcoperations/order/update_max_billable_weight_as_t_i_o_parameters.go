// Code generated by go-swagger; DO NOT EDIT.

package order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// NewUpdateMaxBillableWeightAsTIOParams creates a new UpdateMaxBillableWeightAsTIOParams object
//
// There are no default values defined in the spec.
func NewUpdateMaxBillableWeightAsTIOParams() UpdateMaxBillableWeightAsTIOParams {

	return UpdateMaxBillableWeightAsTIOParams{}
}

// UpdateMaxBillableWeightAsTIOParams contains all the bound params for the update max billable weight as t i o operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateMaxBillableWeightAsTIO
type UpdateMaxBillableWeightAsTIOParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Optimistic locking is implemented via the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a `412 Precondition Failed` error.

	  Required: true
	  In: header
	*/
	IfMatch string
	/*
	  Required: true
	  In: body
	*/
	Body *ghcmessages.UpdateMaxBillableWeightAsTIOPayload
	/*ID of order to use
	  Required: true
	  In: path
	*/
	OrderID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateMaxBillableWeightAsTIOParams() beforehand.
func (o *UpdateMaxBillableWeightAsTIOParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindIfMatch(r.Header[http.CanonicalHeaderKey("If-Match")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body ghcmessages.UpdateMaxBillableWeightAsTIOPayload
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body", ""))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(r.Context())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body", ""))
	}

	rOrderID, rhkOrderID, _ := route.Params.GetOK("orderID")
	if err := o.bindOrderID(rOrderID, rhkOrderID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIfMatch binds and validates parameter IfMatch from header.
func (o *UpdateMaxBillableWeightAsTIOParams) bindIfMatch(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindOrderID binds and validates parameter OrderID from path.
func (o *UpdateMaxBillableWeightAsTIOParams) bindOrderID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("orderID", "path", "strfmt.UUID", raw)
	}
	o.OrderID = *(value.(*strfmt.UUID))

	if err := o.validateOrderID(formats); err != nil {
		return err
	}

	return nil
}

// validateOrderID carries on validations for parameter OrderID
func (o *UpdateMaxBillableWeightAsTIOParams) validateOrderID(formats strfmt.Registry) error {

	if err := validate.FormatOf("orderID", "path", "uuid", o.OrderID.String(), formats); err != nil {
		return err
	}
	return nil
}
