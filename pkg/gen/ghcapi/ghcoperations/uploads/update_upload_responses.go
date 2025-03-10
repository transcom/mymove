// Code generated by go-swagger; DO NOT EDIT.

package uploads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// UpdateUploadCreatedCode is the HTTP code returned for type UpdateUploadCreated
const UpdateUploadCreatedCode int = 201

/*
UpdateUploadCreated updated upload

swagger:response updateUploadCreated
*/
type UpdateUploadCreated struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Upload `json:"body,omitempty"`
}

// NewUpdateUploadCreated creates UpdateUploadCreated with default headers values
func NewUpdateUploadCreated() *UpdateUploadCreated {

	return &UpdateUploadCreated{}
}

// WithPayload adds the payload to the update upload created response
func (o *UpdateUploadCreated) WithPayload(payload *ghcmessages.Upload) *UpdateUploadCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update upload created response
func (o *UpdateUploadCreated) SetPayload(payload *ghcmessages.Upload) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateUploadCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateUploadBadRequestCode is the HTTP code returned for type UpdateUploadBadRequest
const UpdateUploadBadRequestCode int = 400

/*
UpdateUploadBadRequest invalid request

swagger:response updateUploadBadRequest
*/
type UpdateUploadBadRequest struct {
}

// NewUpdateUploadBadRequest creates UpdateUploadBadRequest with default headers values
func NewUpdateUploadBadRequest() *UpdateUploadBadRequest {

	return &UpdateUploadBadRequest{}
}

// WriteResponse to the client
func (o *UpdateUploadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// UpdateUploadForbiddenCode is the HTTP code returned for type UpdateUploadForbidden
const UpdateUploadForbiddenCode int = 403

/*
UpdateUploadForbidden not authorized

swagger:response updateUploadForbidden
*/
type UpdateUploadForbidden struct {
}

// NewUpdateUploadForbidden creates UpdateUploadForbidden with default headers values
func NewUpdateUploadForbidden() *UpdateUploadForbidden {

	return &UpdateUploadForbidden{}
}

// WriteResponse to the client
func (o *UpdateUploadForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// UpdateUploadNotFoundCode is the HTTP code returned for type UpdateUploadNotFound
const UpdateUploadNotFoundCode int = 404

/*
UpdateUploadNotFound not found

swagger:response updateUploadNotFound
*/
type UpdateUploadNotFound struct {
}

// NewUpdateUploadNotFound creates UpdateUploadNotFound with default headers values
func NewUpdateUploadNotFound() *UpdateUploadNotFound {

	return &UpdateUploadNotFound{}
}

// WriteResponse to the client
func (o *UpdateUploadNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// UpdateUploadRequestEntityTooLargeCode is the HTTP code returned for type UpdateUploadRequestEntityTooLarge
const UpdateUploadRequestEntityTooLargeCode int = 413

/*
UpdateUploadRequestEntityTooLarge payload is too large

swagger:response updateUploadRequestEntityTooLarge
*/
type UpdateUploadRequestEntityTooLarge struct {
}

// NewUpdateUploadRequestEntityTooLarge creates UpdateUploadRequestEntityTooLarge with default headers values
func NewUpdateUploadRequestEntityTooLarge() *UpdateUploadRequestEntityTooLarge {

	return &UpdateUploadRequestEntityTooLarge{}
}

// WriteResponse to the client
func (o *UpdateUploadRequestEntityTooLarge) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(413)
}

// UpdateUploadInternalServerErrorCode is the HTTP code returned for type UpdateUploadInternalServerError
const UpdateUploadInternalServerErrorCode int = 500

/*
UpdateUploadInternalServerError server error

swagger:response updateUploadInternalServerError
*/
type UpdateUploadInternalServerError struct {
}

// NewUpdateUploadInternalServerError creates UpdateUploadInternalServerError with default headers values
func NewUpdateUploadInternalServerError() *UpdateUploadInternalServerError {

	return &UpdateUploadInternalServerError{}
}

// WriteResponse to the client
func (o *UpdateUploadInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
