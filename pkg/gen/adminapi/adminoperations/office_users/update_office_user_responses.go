// Code generated by go-swagger; DO NOT EDIT.

package office_users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
)

// UpdateOfficeUserOKCode is the HTTP code returned for type UpdateOfficeUserOK
const UpdateOfficeUserOKCode int = 200

/*
UpdateOfficeUserOK Successfully updated Office User

swagger:response updateOfficeUserOK
*/
type UpdateOfficeUserOK struct {

	/*
	  In: Body
	*/
	Payload *adminmessages.OfficeUser `json:"body,omitempty"`
}

// NewUpdateOfficeUserOK creates UpdateOfficeUserOK with default headers values
func NewUpdateOfficeUserOK() *UpdateOfficeUserOK {

	return &UpdateOfficeUserOK{}
}

// WithPayload adds the payload to the update office user o k response
func (o *UpdateOfficeUserOK) WithPayload(payload *adminmessages.OfficeUser) *UpdateOfficeUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update office user o k response
func (o *UpdateOfficeUserOK) SetPayload(payload *adminmessages.OfficeUser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateOfficeUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateOfficeUserBadRequestCode is the HTTP code returned for type UpdateOfficeUserBadRequest
const UpdateOfficeUserBadRequestCode int = 400

/*
UpdateOfficeUserBadRequest Invalid Request

swagger:response updateOfficeUserBadRequest
*/
type UpdateOfficeUserBadRequest struct {
}

// NewUpdateOfficeUserBadRequest creates UpdateOfficeUserBadRequest with default headers values
func NewUpdateOfficeUserBadRequest() *UpdateOfficeUserBadRequest {

	return &UpdateOfficeUserBadRequest{}
}

// WriteResponse to the client
func (o *UpdateOfficeUserBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// UpdateOfficeUserUnauthorizedCode is the HTTP code returned for type UpdateOfficeUserUnauthorized
const UpdateOfficeUserUnauthorizedCode int = 401

/*
UpdateOfficeUserUnauthorized Must be authenticated to use this end point

swagger:response updateOfficeUserUnauthorized
*/
type UpdateOfficeUserUnauthorized struct {
}

// NewUpdateOfficeUserUnauthorized creates UpdateOfficeUserUnauthorized with default headers values
func NewUpdateOfficeUserUnauthorized() *UpdateOfficeUserUnauthorized {

	return &UpdateOfficeUserUnauthorized{}
}

// WriteResponse to the client
func (o *UpdateOfficeUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// UpdateOfficeUserForbiddenCode is the HTTP code returned for type UpdateOfficeUserForbidden
const UpdateOfficeUserForbiddenCode int = 403

/*
UpdateOfficeUserForbidden Not authorized to update an Office User

swagger:response updateOfficeUserForbidden
*/
type UpdateOfficeUserForbidden struct {
}

// NewUpdateOfficeUserForbidden creates UpdateOfficeUserForbidden with default headers values
func NewUpdateOfficeUserForbidden() *UpdateOfficeUserForbidden {

	return &UpdateOfficeUserForbidden{}
}

// WriteResponse to the client
func (o *UpdateOfficeUserForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// UpdateOfficeUserNotFoundCode is the HTTP code returned for type UpdateOfficeUserNotFound
const UpdateOfficeUserNotFoundCode int = 404

/*
UpdateOfficeUserNotFound Office User not found

swagger:response updateOfficeUserNotFound
*/
type UpdateOfficeUserNotFound struct {
}

// NewUpdateOfficeUserNotFound creates UpdateOfficeUserNotFound with default headers values
func NewUpdateOfficeUserNotFound() *UpdateOfficeUserNotFound {

	return &UpdateOfficeUserNotFound{}
}

// WriteResponse to the client
func (o *UpdateOfficeUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// UpdateOfficeUserInternalServerErrorCode is the HTTP code returned for type UpdateOfficeUserInternalServerError
const UpdateOfficeUserInternalServerErrorCode int = 500

/*
UpdateOfficeUserInternalServerError Server error

swagger:response updateOfficeUserInternalServerError
*/
type UpdateOfficeUserInternalServerError struct {
}

// NewUpdateOfficeUserInternalServerError creates UpdateOfficeUserInternalServerError with default headers values
func NewUpdateOfficeUserInternalServerError() *UpdateOfficeUserInternalServerError {

	return &UpdateOfficeUserInternalServerError{}
}

// WriteResponse to the client
func (o *UpdateOfficeUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
