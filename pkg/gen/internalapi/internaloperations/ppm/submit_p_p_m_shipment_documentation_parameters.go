// Code generated by go-swagger; DO NOT EDIT.

package ppm

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

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// NewSubmitPPMShipmentDocumentationParams creates a new SubmitPPMShipmentDocumentationParams object
//
// There are no default values defined in the spec.
func NewSubmitPPMShipmentDocumentationParams() SubmitPPMShipmentDocumentationParams {

	return SubmitPPMShipmentDocumentationParams{}
}

// SubmitPPMShipmentDocumentationParams contains all the bound params for the submit p p m shipment documentation operation
// typically these are obtained from a http.Request
//
// swagger:parameters submitPPMShipmentDocumentation
type SubmitPPMShipmentDocumentationParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*UUID of the PPM shipment
	  Required: true
	  In: path
	*/
	PpmShipmentID strfmt.UUID
	/*
	  Required: true
	  In: body
	*/
	SavePPMShipmentSignedCertificationPayload *internalmessages.SavePPMShipmentSignedCertification
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSubmitPPMShipmentDocumentationParams() beforehand.
func (o *SubmitPPMShipmentDocumentationParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rPpmShipmentID, rhkPpmShipmentID, _ := route.Params.GetOK("ppmShipmentId")
	if err := o.bindPpmShipmentID(rPpmShipmentID, rhkPpmShipmentID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body internalmessages.SavePPMShipmentSignedCertification
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("savePPMShipmentSignedCertificationPayload", "body", ""))
			} else {
				res = append(res, errors.NewParseError("savePPMShipmentSignedCertificationPayload", "body", "", err))
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
				o.SavePPMShipmentSignedCertificationPayload = &body
			}
		}
	} else {
		res = append(res, errors.Required("savePPMShipmentSignedCertificationPayload", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPpmShipmentID binds and validates parameter PpmShipmentID from path.
func (o *SubmitPPMShipmentDocumentationParams) bindPpmShipmentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("ppmShipmentId", "path", "strfmt.UUID", raw)
	}
	o.PpmShipmentID = *(value.(*strfmt.UUID))

	if err := o.validatePpmShipmentID(formats); err != nil {
		return err
	}

	return nil
}

// validatePpmShipmentID carries on validations for parameter PpmShipmentID
func (o *SubmitPPMShipmentDocumentationParams) validatePpmShipmentID(formats strfmt.Registry) error {

	if err := validate.FormatOf("ppmShipmentId", "path", "uuid", o.PpmShipmentID.String(), formats); err != nil {
		return err
	}
	return nil
}