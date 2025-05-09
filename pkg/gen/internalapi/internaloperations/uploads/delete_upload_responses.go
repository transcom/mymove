// Code generated by go-swagger; DO NOT EDIT.

package uploads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DeleteUploadNoContentCode is the HTTP code returned for type DeleteUploadNoContent
const DeleteUploadNoContentCode int = 204

/*
DeleteUploadNoContent deleted

swagger:response deleteUploadNoContent
*/
type DeleteUploadNoContent struct {
}

// NewDeleteUploadNoContent creates DeleteUploadNoContent with default headers values
func NewDeleteUploadNoContent() *DeleteUploadNoContent {

	return &DeleteUploadNoContent{}
}

// WriteResponse to the client
func (o *DeleteUploadNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteUploadBadRequestCode is the HTTP code returned for type DeleteUploadBadRequest
const DeleteUploadBadRequestCode int = 400

/*
DeleteUploadBadRequest invalid request

swagger:response deleteUploadBadRequest
*/
type DeleteUploadBadRequest struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.InvalidRequestResponsePayload `json:"body,omitempty"`
}

// NewDeleteUploadBadRequest creates DeleteUploadBadRequest with default headers values
func NewDeleteUploadBadRequest() *DeleteUploadBadRequest {

	return &DeleteUploadBadRequest{}
}

// WithPayload adds the payload to the delete upload bad request response
func (o *DeleteUploadBadRequest) WithPayload(payload *internalmessages.InvalidRequestResponsePayload) *DeleteUploadBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete upload bad request response
func (o *DeleteUploadBadRequest) SetPayload(payload *internalmessages.InvalidRequestResponsePayload) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteUploadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteUploadForbiddenCode is the HTTP code returned for type DeleteUploadForbidden
const DeleteUploadForbiddenCode int = 403

/*
DeleteUploadForbidden not authorized

swagger:response deleteUploadForbidden
*/
type DeleteUploadForbidden struct {
}

// NewDeleteUploadForbidden creates DeleteUploadForbidden with default headers values
func NewDeleteUploadForbidden() *DeleteUploadForbidden {

	return &DeleteUploadForbidden{}
}

// WriteResponse to the client
func (o *DeleteUploadForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// DeleteUploadNotFoundCode is the HTTP code returned for type DeleteUploadNotFound
const DeleteUploadNotFoundCode int = 404

/*
DeleteUploadNotFound not found

swagger:response deleteUploadNotFound
*/
type DeleteUploadNotFound struct {
}

// NewDeleteUploadNotFound creates DeleteUploadNotFound with default headers values
func NewDeleteUploadNotFound() *DeleteUploadNotFound {

	return &DeleteUploadNotFound{}
}

// WriteResponse to the client
func (o *DeleteUploadNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteUploadInternalServerErrorCode is the HTTP code returned for type DeleteUploadInternalServerError
const DeleteUploadInternalServerErrorCode int = 500

/*
DeleteUploadInternalServerError server error

swagger:response deleteUploadInternalServerError
*/
type DeleteUploadInternalServerError struct {
}

// NewDeleteUploadInternalServerError creates DeleteUploadInternalServerError with default headers values
func NewDeleteUploadInternalServerError() *DeleteUploadInternalServerError {

	return &DeleteUploadInternalServerError{}
}

// WriteResponse to the client
func (o *DeleteUploadInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
