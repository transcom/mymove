// Code generated by go-swagger; DO NOT EDIT.

package move

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DeleteAssignedOfficeUserHandlerFunc turns a function with the right signature into a delete assigned office user handler
type DeleteAssignedOfficeUserHandlerFunc func(DeleteAssignedOfficeUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteAssignedOfficeUserHandlerFunc) Handle(params DeleteAssignedOfficeUserParams) middleware.Responder {
	return fn(params)
}

// DeleteAssignedOfficeUserHandler interface for that can handle valid delete assigned office user params
type DeleteAssignedOfficeUserHandler interface {
	Handle(DeleteAssignedOfficeUserParams) middleware.Responder
}

// NewDeleteAssignedOfficeUser creates a new http.Handler for the delete assigned office user operation
func NewDeleteAssignedOfficeUser(ctx *middleware.Context, handler DeleteAssignedOfficeUserHandler) *DeleteAssignedOfficeUser {
	return &DeleteAssignedOfficeUser{Context: ctx, Handler: handler}
}

/*
	DeleteAssignedOfficeUser swagger:route PATCH /moves/{moveID}/unassignOfficeUser move deleteAssignedOfficeUser

unassigns either a services counselor, task ordering officer, or task invoicing officer from the move
*/
type DeleteAssignedOfficeUser struct {
	Context *middleware.Context
	Handler DeleteAssignedOfficeUserHandler
}

func (o *DeleteAssignedOfficeUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteAssignedOfficeUserParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// DeleteAssignedOfficeUserBody delete assigned office user body
//
// swagger:model DeleteAssignedOfficeUserBody
type DeleteAssignedOfficeUserBody struct {

	// role type
	// Required: true
	RoleType *string `json:"roleType"`
}

// Validate validates this delete assigned office user body
func (o *DeleteAssignedOfficeUserBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateRoleType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteAssignedOfficeUserBody) validateRoleType(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"roleType", "body", o.RoleType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this delete assigned office user body based on context it is used
func (o *DeleteAssignedOfficeUserBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DeleteAssignedOfficeUserBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteAssignedOfficeUserBody) UnmarshalBinary(b []byte) error {
	var res DeleteAssignedOfficeUserBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}