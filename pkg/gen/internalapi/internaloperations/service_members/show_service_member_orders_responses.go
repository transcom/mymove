// Code generated by go-swagger; DO NOT EDIT.

package service_members

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ShowServiceMemberOrdersOKCode is the HTTP code returned for type ShowServiceMemberOrdersOK
const ShowServiceMemberOrdersOKCode int = 200

/*
ShowServiceMemberOrdersOK the instance of the service member

swagger:response showServiceMemberOrdersOK
*/
type ShowServiceMemberOrdersOK struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.Orders `json:"body,omitempty"`
}

// NewShowServiceMemberOrdersOK creates ShowServiceMemberOrdersOK with default headers values
func NewShowServiceMemberOrdersOK() *ShowServiceMemberOrdersOK {

	return &ShowServiceMemberOrdersOK{}
}

// WithPayload adds the payload to the show service member orders o k response
func (o *ShowServiceMemberOrdersOK) WithPayload(payload *internalmessages.Orders) *ShowServiceMemberOrdersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the show service member orders o k response
func (o *ShowServiceMemberOrdersOK) SetPayload(payload *internalmessages.Orders) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ShowServiceMemberOrdersBadRequestCode is the HTTP code returned for type ShowServiceMemberOrdersBadRequest
const ShowServiceMemberOrdersBadRequestCode int = 400

/*
ShowServiceMemberOrdersBadRequest invalid request

swagger:response showServiceMemberOrdersBadRequest
*/
type ShowServiceMemberOrdersBadRequest struct {
}

// NewShowServiceMemberOrdersBadRequest creates ShowServiceMemberOrdersBadRequest with default headers values
func NewShowServiceMemberOrdersBadRequest() *ShowServiceMemberOrdersBadRequest {

	return &ShowServiceMemberOrdersBadRequest{}
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// ShowServiceMemberOrdersUnauthorizedCode is the HTTP code returned for type ShowServiceMemberOrdersUnauthorized
const ShowServiceMemberOrdersUnauthorizedCode int = 401

/*
ShowServiceMemberOrdersUnauthorized request requires user authentication

swagger:response showServiceMemberOrdersUnauthorized
*/
type ShowServiceMemberOrdersUnauthorized struct {
}

// NewShowServiceMemberOrdersUnauthorized creates ShowServiceMemberOrdersUnauthorized with default headers values
func NewShowServiceMemberOrdersUnauthorized() *ShowServiceMemberOrdersUnauthorized {

	return &ShowServiceMemberOrdersUnauthorized{}
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// ShowServiceMemberOrdersForbiddenCode is the HTTP code returned for type ShowServiceMemberOrdersForbidden
const ShowServiceMemberOrdersForbiddenCode int = 403

/*
ShowServiceMemberOrdersForbidden user is not authorized

swagger:response showServiceMemberOrdersForbidden
*/
type ShowServiceMemberOrdersForbidden struct {
}

// NewShowServiceMemberOrdersForbidden creates ShowServiceMemberOrdersForbidden with default headers values
func NewShowServiceMemberOrdersForbidden() *ShowServiceMemberOrdersForbidden {

	return &ShowServiceMemberOrdersForbidden{}
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// ShowServiceMemberOrdersNotFoundCode is the HTTP code returned for type ShowServiceMemberOrdersNotFound
const ShowServiceMemberOrdersNotFoundCode int = 404

/*
ShowServiceMemberOrdersNotFound service member not found

swagger:response showServiceMemberOrdersNotFound
*/
type ShowServiceMemberOrdersNotFound struct {
}

// NewShowServiceMemberOrdersNotFound creates ShowServiceMemberOrdersNotFound with default headers values
func NewShowServiceMemberOrdersNotFound() *ShowServiceMemberOrdersNotFound {

	return &ShowServiceMemberOrdersNotFound{}
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// ShowServiceMemberOrdersInternalServerErrorCode is the HTTP code returned for type ShowServiceMemberOrdersInternalServerError
const ShowServiceMemberOrdersInternalServerErrorCode int = 500

/*
ShowServiceMemberOrdersInternalServerError internal server error

swagger:response showServiceMemberOrdersInternalServerError
*/
type ShowServiceMemberOrdersInternalServerError struct {
}

// NewShowServiceMemberOrdersInternalServerError creates ShowServiceMemberOrdersInternalServerError with default headers values
func NewShowServiceMemberOrdersInternalServerError() *ShowServiceMemberOrdersInternalServerError {

	return &ShowServiceMemberOrdersInternalServerError{}
}

// WriteResponse to the client
func (o *ShowServiceMemberOrdersInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
