// Code generated by go-swagger; DO NOT EDIT.

package requested_office_users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
)

// GetRequestedOfficeUserOKCode is the HTTP code returned for type GetRequestedOfficeUserOK
const GetRequestedOfficeUserOKCode int = 200

/*
GetRequestedOfficeUserOK success

swagger:response getRequestedOfficeUserOK
*/
type GetRequestedOfficeUserOK struct {

	/*
	  In: Body
	*/
	Payload *adminmessages.OfficeUser `json:"body,omitempty"`
}

// NewGetRequestedOfficeUserOK creates GetRequestedOfficeUserOK with default headers values
func NewGetRequestedOfficeUserOK() *GetRequestedOfficeUserOK {

	return &GetRequestedOfficeUserOK{}
}

// WithPayload adds the payload to the get requested office user o k response
func (o *GetRequestedOfficeUserOK) WithPayload(payload *adminmessages.OfficeUser) *GetRequestedOfficeUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get requested office user o k response
func (o *GetRequestedOfficeUserOK) SetPayload(payload *adminmessages.OfficeUser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRequestedOfficeUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRequestedOfficeUserBadRequestCode is the HTTP code returned for type GetRequestedOfficeUserBadRequest
const GetRequestedOfficeUserBadRequestCode int = 400

/*
GetRequestedOfficeUserBadRequest invalid request

swagger:response getRequestedOfficeUserBadRequest
*/
type GetRequestedOfficeUserBadRequest struct {
}

// NewGetRequestedOfficeUserBadRequest creates GetRequestedOfficeUserBadRequest with default headers values
func NewGetRequestedOfficeUserBadRequest() *GetRequestedOfficeUserBadRequest {

	return &GetRequestedOfficeUserBadRequest{}
}

// WriteResponse to the client
func (o *GetRequestedOfficeUserBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetRequestedOfficeUserUnauthorizedCode is the HTTP code returned for type GetRequestedOfficeUserUnauthorized
const GetRequestedOfficeUserUnauthorizedCode int = 401

/*
GetRequestedOfficeUserUnauthorized request requires user authentication

swagger:response getRequestedOfficeUserUnauthorized
*/
type GetRequestedOfficeUserUnauthorized struct {
}

// NewGetRequestedOfficeUserUnauthorized creates GetRequestedOfficeUserUnauthorized with default headers values
func NewGetRequestedOfficeUserUnauthorized() *GetRequestedOfficeUserUnauthorized {

	return &GetRequestedOfficeUserUnauthorized{}
}

// WriteResponse to the client
func (o *GetRequestedOfficeUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetRequestedOfficeUserNotFoundCode is the HTTP code returned for type GetRequestedOfficeUserNotFound
const GetRequestedOfficeUserNotFoundCode int = 404

/*
GetRequestedOfficeUserNotFound Office User not found

swagger:response getRequestedOfficeUserNotFound
*/
type GetRequestedOfficeUserNotFound struct {
}

// NewGetRequestedOfficeUserNotFound creates GetRequestedOfficeUserNotFound with default headers values
func NewGetRequestedOfficeUserNotFound() *GetRequestedOfficeUserNotFound {

	return &GetRequestedOfficeUserNotFound{}
}

// WriteResponse to the client
func (o *GetRequestedOfficeUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetRequestedOfficeUserInternalServerErrorCode is the HTTP code returned for type GetRequestedOfficeUserInternalServerError
const GetRequestedOfficeUserInternalServerErrorCode int = 500

/*
GetRequestedOfficeUserInternalServerError server error

swagger:response getRequestedOfficeUserInternalServerError
*/
type GetRequestedOfficeUserInternalServerError struct {
}

// NewGetRequestedOfficeUserInternalServerError creates GetRequestedOfficeUserInternalServerError with default headers values
func NewGetRequestedOfficeUserInternalServerError() *GetRequestedOfficeUserInternalServerError {

	return &GetRequestedOfficeUserInternalServerError{}
}

// WriteResponse to the client
func (o *GetRequestedOfficeUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}