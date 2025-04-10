// Code generated by go-swagger; DO NOT EDIT.

package payment_requests

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

// NewUpdatePaymentRequestStatusParams creates a new UpdatePaymentRequestStatusParams object
//
// There are no default values defined in the spec.
func NewUpdatePaymentRequestStatusParams() UpdatePaymentRequestStatusParams {

	return UpdatePaymentRequestStatusParams{}
}

// UpdatePaymentRequestStatusParams contains all the bound params for the update payment request status operation
// typically these are obtained from a http.Request
//
// swagger:parameters updatePaymentRequestStatus
type UpdatePaymentRequestStatusParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: header
	*/
	IfMatch string
	/*
	  Required: true
	  In: body
	*/
	Body *ghcmessages.UpdatePaymentRequestStatusPayload
	/*UUID of payment request
	  Required: true
	  In: path
	*/
	PaymentRequestID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdatePaymentRequestStatusParams() beforehand.
func (o *UpdatePaymentRequestStatusParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindIfMatch(r.Header[http.CanonicalHeaderKey("If-Match")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body ghcmessages.UpdatePaymentRequestStatusPayload
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

	rPaymentRequestID, rhkPaymentRequestID, _ := route.Params.GetOK("paymentRequestID")
	if err := o.bindPaymentRequestID(rPaymentRequestID, rhkPaymentRequestID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIfMatch binds and validates parameter IfMatch from header.
func (o *UpdatePaymentRequestStatusParams) bindIfMatch(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindPaymentRequestID binds and validates parameter PaymentRequestID from path.
func (o *UpdatePaymentRequestStatusParams) bindPaymentRequestID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("paymentRequestID", "path", "strfmt.UUID", raw)
	}
	o.PaymentRequestID = *(value.(*strfmt.UUID))

	if err := o.validatePaymentRequestID(formats); err != nil {
		return err
	}

	return nil
}

// validatePaymentRequestID carries on validations for parameter PaymentRequestID
func (o *UpdatePaymentRequestStatusParams) validatePaymentRequestID(formats strfmt.Registry) error {

	if err := validate.FormatOf("paymentRequestID", "path", "uuid", o.PaymentRequestID.String(), formats); err != nil {
		return err
	}
	return nil
}
