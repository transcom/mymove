// Code generated by go-swagger; DO NOT EDIT.

package move_docs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// CreateGenericMoveDocumentOKCode is the HTTP code returned for type CreateGenericMoveDocumentOK
const CreateGenericMoveDocumentOKCode int = 200

/*
CreateGenericMoveDocumentOK returns new move document object

swagger:response createGenericMoveDocumentOK
*/
type CreateGenericMoveDocumentOK struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.MoveDocumentPayload `json:"body,omitempty"`
}

// NewCreateGenericMoveDocumentOK creates CreateGenericMoveDocumentOK with default headers values
func NewCreateGenericMoveDocumentOK() *CreateGenericMoveDocumentOK {

	return &CreateGenericMoveDocumentOK{}
}

// WithPayload adds the payload to the create generic move document o k response
func (o *CreateGenericMoveDocumentOK) WithPayload(payload *internalmessages.MoveDocumentPayload) *CreateGenericMoveDocumentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create generic move document o k response
func (o *CreateGenericMoveDocumentOK) SetPayload(payload *internalmessages.MoveDocumentPayload) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateGenericMoveDocumentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateGenericMoveDocumentBadRequestCode is the HTTP code returned for type CreateGenericMoveDocumentBadRequest
const CreateGenericMoveDocumentBadRequestCode int = 400

/*
CreateGenericMoveDocumentBadRequest invalid request

swagger:response createGenericMoveDocumentBadRequest
*/
type CreateGenericMoveDocumentBadRequest struct {
}

// NewCreateGenericMoveDocumentBadRequest creates CreateGenericMoveDocumentBadRequest with default headers values
func NewCreateGenericMoveDocumentBadRequest() *CreateGenericMoveDocumentBadRequest {

	return &CreateGenericMoveDocumentBadRequest{}
}

// WriteResponse to the client
func (o *CreateGenericMoveDocumentBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// CreateGenericMoveDocumentUnauthorizedCode is the HTTP code returned for type CreateGenericMoveDocumentUnauthorized
const CreateGenericMoveDocumentUnauthorizedCode int = 401

/*
CreateGenericMoveDocumentUnauthorized must be authenticated to use this endpoint

swagger:response createGenericMoveDocumentUnauthorized
*/
type CreateGenericMoveDocumentUnauthorized struct {
}

// NewCreateGenericMoveDocumentUnauthorized creates CreateGenericMoveDocumentUnauthorized with default headers values
func NewCreateGenericMoveDocumentUnauthorized() *CreateGenericMoveDocumentUnauthorized {

	return &CreateGenericMoveDocumentUnauthorized{}
}

// WriteResponse to the client
func (o *CreateGenericMoveDocumentUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// CreateGenericMoveDocumentForbiddenCode is the HTTP code returned for type CreateGenericMoveDocumentForbidden
const CreateGenericMoveDocumentForbiddenCode int = 403

/*
CreateGenericMoveDocumentForbidden not authorized to modify this move

swagger:response createGenericMoveDocumentForbidden
*/
type CreateGenericMoveDocumentForbidden struct {
}

// NewCreateGenericMoveDocumentForbidden creates CreateGenericMoveDocumentForbidden with default headers values
func NewCreateGenericMoveDocumentForbidden() *CreateGenericMoveDocumentForbidden {

	return &CreateGenericMoveDocumentForbidden{}
}

// WriteResponse to the client
func (o *CreateGenericMoveDocumentForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// CreateGenericMoveDocumentInternalServerErrorCode is the HTTP code returned for type CreateGenericMoveDocumentInternalServerError
const CreateGenericMoveDocumentInternalServerErrorCode int = 500

/*
CreateGenericMoveDocumentInternalServerError server error

swagger:response createGenericMoveDocumentInternalServerError
*/
type CreateGenericMoveDocumentInternalServerError struct {
}

// NewCreateGenericMoveDocumentInternalServerError creates CreateGenericMoveDocumentInternalServerError with default headers values
func NewCreateGenericMoveDocumentInternalServerError() *CreateGenericMoveDocumentInternalServerError {

	return &CreateGenericMoveDocumentInternalServerError{}
}

// WriteResponse to the client
func (o *CreateGenericMoveDocumentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}