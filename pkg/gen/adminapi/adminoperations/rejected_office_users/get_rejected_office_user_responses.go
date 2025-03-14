// Code generated by go-swagger; DO NOT EDIT.

package rejected_office_users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
)

// GetRejectedOfficeUserOKCode is the HTTP code returned for type GetRejectedOfficeUserOK
const GetRejectedOfficeUserOKCode int = 200

/*
GetRejectedOfficeUserOK success

swagger:response getRejectedOfficeUserOK
*/
type GetRejectedOfficeUserOK struct {

	/*
	  In: Body
	*/
	Payload *adminmessages.OfficeUser `json:"body,omitempty"`
}

// NewGetRejectedOfficeUserOK creates GetRejectedOfficeUserOK with default headers values
func NewGetRejectedOfficeUserOK() *GetRejectedOfficeUserOK {

	return &GetRejectedOfficeUserOK{}
}

// WithPayload adds the payload to the get rejected office user o k response
func (o *GetRejectedOfficeUserOK) WithPayload(payload *adminmessages.OfficeUser) *GetRejectedOfficeUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get rejected office user o k response
func (o *GetRejectedOfficeUserOK) SetPayload(payload *adminmessages.OfficeUser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRejectedOfficeUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRejectedOfficeUserBadRequestCode is the HTTP code returned for type GetRejectedOfficeUserBadRequest
const GetRejectedOfficeUserBadRequestCode int = 400

/*
GetRejectedOfficeUserBadRequest invalid request

swagger:response getRejectedOfficeUserBadRequest
*/
type GetRejectedOfficeUserBadRequest struct {
}

// NewGetRejectedOfficeUserBadRequest creates GetRejectedOfficeUserBadRequest with default headers values
func NewGetRejectedOfficeUserBadRequest() *GetRejectedOfficeUserBadRequest {

	return &GetRejectedOfficeUserBadRequest{}
}

// WriteResponse to the client
func (o *GetRejectedOfficeUserBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetRejectedOfficeUserUnauthorizedCode is the HTTP code returned for type GetRejectedOfficeUserUnauthorized
const GetRejectedOfficeUserUnauthorizedCode int = 401

/*
GetRejectedOfficeUserUnauthorized request requires user authentication

swagger:response getRejectedOfficeUserUnauthorized
*/
type GetRejectedOfficeUserUnauthorized struct {
}

// NewGetRejectedOfficeUserUnauthorized creates GetRejectedOfficeUserUnauthorized with default headers values
func NewGetRejectedOfficeUserUnauthorized() *GetRejectedOfficeUserUnauthorized {

	return &GetRejectedOfficeUserUnauthorized{}
}

// WriteResponse to the client
func (o *GetRejectedOfficeUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetRejectedOfficeUserNotFoundCode is the HTTP code returned for type GetRejectedOfficeUserNotFound
const GetRejectedOfficeUserNotFoundCode int = 404

/*
GetRejectedOfficeUserNotFound Office User not found

swagger:response getRejectedOfficeUserNotFound
*/
type GetRejectedOfficeUserNotFound struct {
}

// NewGetRejectedOfficeUserNotFound creates GetRejectedOfficeUserNotFound with default headers values
func NewGetRejectedOfficeUserNotFound() *GetRejectedOfficeUserNotFound {

	return &GetRejectedOfficeUserNotFound{}
}

// WriteResponse to the client
func (o *GetRejectedOfficeUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetRejectedOfficeUserInternalServerErrorCode is the HTTP code returned for type GetRejectedOfficeUserInternalServerError
const GetRejectedOfficeUserInternalServerErrorCode int = 500

/*
GetRejectedOfficeUserInternalServerError server error

swagger:response getRejectedOfficeUserInternalServerError
*/
type GetRejectedOfficeUserInternalServerError struct {
}

// NewGetRejectedOfficeUserInternalServerError creates GetRejectedOfficeUserInternalServerError with default headers values
func NewGetRejectedOfficeUserInternalServerError() *GetRejectedOfficeUserInternalServerError {

	return &GetRejectedOfficeUserInternalServerError{}
}

// WriteResponse to the client
func (o *GetRejectedOfficeUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
