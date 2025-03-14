// Code generated by go-swagger; DO NOT EDIT.

package backup_contacts

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

// NewCreateServiceMemberBackupContactParams creates a new CreateServiceMemberBackupContactParams object
//
// There are no default values defined in the spec.
func NewCreateServiceMemberBackupContactParams() CreateServiceMemberBackupContactParams {

	return CreateServiceMemberBackupContactParams{}
}

// CreateServiceMemberBackupContactParams contains all the bound params for the create service member backup contact operation
// typically these are obtained from a http.Request
//
// swagger:parameters createServiceMemberBackupContact
type CreateServiceMemberBackupContactParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	CreateBackupContactPayload *internalmessages.CreateServiceMemberBackupContactPayload
	/*UUID of the service member
	  Required: true
	  In: path
	*/
	ServiceMemberID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateServiceMemberBackupContactParams() beforehand.
func (o *CreateServiceMemberBackupContactParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body internalmessages.CreateServiceMemberBackupContactPayload
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("createBackupContactPayload", "body", ""))
			} else {
				res = append(res, errors.NewParseError("createBackupContactPayload", "body", "", err))
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
				o.CreateBackupContactPayload = &body
			}
		}
	} else {
		res = append(res, errors.Required("createBackupContactPayload", "body", ""))
	}

	rServiceMemberID, rhkServiceMemberID, _ := route.Params.GetOK("serviceMemberId")
	if err := o.bindServiceMemberID(rServiceMemberID, rhkServiceMemberID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindServiceMemberID binds and validates parameter ServiceMemberID from path.
func (o *CreateServiceMemberBackupContactParams) bindServiceMemberID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("serviceMemberId", "path", "strfmt.UUID", raw)
	}
	o.ServiceMemberID = *(value.(*strfmt.UUID))

	if err := o.validateServiceMemberID(formats); err != nil {
		return err
	}

	return nil
}

// validateServiceMemberID carries on validations for parameter ServiceMemberID
func (o *CreateServiceMemberBackupContactParams) validateServiceMemberID(formats strfmt.Registry) error {

	if err := validate.FormatOf("serviceMemberId", "path", "uuid", o.ServiceMemberID.String(), formats); err != nil {
		return err
	}
	return nil
}
