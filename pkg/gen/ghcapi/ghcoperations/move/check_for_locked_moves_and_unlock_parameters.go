// Code generated by go-swagger; DO NOT EDIT.

package move

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewCheckForLockedMovesAndUnlockParams creates a new CheckForLockedMovesAndUnlockParams object
//
// There are no default values defined in the spec.
func NewCheckForLockedMovesAndUnlockParams() CheckForLockedMovesAndUnlockParams {

	return CheckForLockedMovesAndUnlockParams{}
}

// CheckForLockedMovesAndUnlockParams contains all the bound params for the check for locked moves and unlock operation
// typically these are obtained from a http.Request
//
// swagger:parameters checkForLockedMovesAndUnlock
type CheckForLockedMovesAndUnlockParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of the move's officer
	  Required: true
	  In: path
	*/
	OfficeUserID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCheckForLockedMovesAndUnlockParams() beforehand.
func (o *CheckForLockedMovesAndUnlockParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rOfficeUserID, rhkOfficeUserID, _ := route.Params.GetOK("officeUserID")
	if err := o.bindOfficeUserID(rOfficeUserID, rhkOfficeUserID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindOfficeUserID binds and validates parameter OfficeUserID from path.
func (o *CheckForLockedMovesAndUnlockParams) bindOfficeUserID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("officeUserID", "path", "strfmt.UUID", raw)
	}
	o.OfficeUserID = *(value.(*strfmt.UUID))

	if err := o.validateOfficeUserID(formats); err != nil {
		return err
	}

	return nil
}

// validateOfficeUserID carries on validations for parameter OfficeUserID
func (o *CheckForLockedMovesAndUnlockParams) validateOfficeUserID(formats strfmt.Registry) error {

	if err := validate.FormatOf("officeUserID", "path", "uuid", o.OfficeUserID.String(), formats); err != nil {
		return err
	}
	return nil
}
