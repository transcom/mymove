// Code generated by go-swagger; DO NOT EDIT.

package mto_service_item

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

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// NewUpdateMTOServiceItemParams creates a new UpdateMTOServiceItemParams object
//
// There are no default values defined in the spec.
func NewUpdateMTOServiceItemParams() UpdateMTOServiceItemParams {

	return UpdateMTOServiceItemParams{}
}

// UpdateMTOServiceItemParams contains all the bound params for the update m t o service item operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateMTOServiceItem
type UpdateMTOServiceItemParams struct {

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
	Body primemessages.UpdateMTOServiceItem
	/*UUID of service item to update.
	  Required: true
	  In: path
	*/
	MtoServiceItemID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateMTOServiceItemParams() beforehand.
func (o *UpdateMTOServiceItemParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindIfMatch(r.Header[http.CanonicalHeaderKey("If-Match")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		body, err := primemessages.UnmarshalUpdateMTOServiceItem(r.Body, route.Consumer)
		if err != nil {
			if err == io.EOF {
				err = errors.Required("body", "body", "")
			}
			res = append(res, err)
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
				o.Body = body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body", ""))
	}

	rMtoServiceItemID, rhkMtoServiceItemID, _ := route.Params.GetOK("mtoServiceItemID")
	if err := o.bindMtoServiceItemID(rMtoServiceItemID, rhkMtoServiceItemID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIfMatch binds and validates parameter IfMatch from header.
func (o *UpdateMTOServiceItemParams) bindIfMatch(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindMtoServiceItemID binds and validates parameter MtoServiceItemID from path.
func (o *UpdateMTOServiceItemParams) bindMtoServiceItemID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.MtoServiceItemID = raw

	return nil
}