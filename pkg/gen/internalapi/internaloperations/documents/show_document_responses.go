// Code generated by go-swagger; DO NOT EDIT.

package documents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ShowDocumentOKCode is the HTTP code returned for type ShowDocumentOK
const ShowDocumentOKCode int = 200

/*
ShowDocumentOK the requested document

swagger:response showDocumentOK
*/
type ShowDocumentOK struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.Document `json:"body,omitempty"`
}

// NewShowDocumentOK creates ShowDocumentOK with default headers values
func NewShowDocumentOK() *ShowDocumentOK {

	return &ShowDocumentOK{}
}

// WithPayload adds the payload to the show document o k response
func (o *ShowDocumentOK) WithPayload(payload *internalmessages.Document) *ShowDocumentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the show document o k response
func (o *ShowDocumentOK) SetPayload(payload *internalmessages.Document) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ShowDocumentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ShowDocumentBadRequestCode is the HTTP code returned for type ShowDocumentBadRequest
const ShowDocumentBadRequestCode int = 400

/*
ShowDocumentBadRequest invalid request

swagger:response showDocumentBadRequest
*/
type ShowDocumentBadRequest struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.InvalidRequestResponsePayload `json:"body,omitempty"`
}

// NewShowDocumentBadRequest creates ShowDocumentBadRequest with default headers values
func NewShowDocumentBadRequest() *ShowDocumentBadRequest {

	return &ShowDocumentBadRequest{}
}

// WithPayload adds the payload to the show document bad request response
func (o *ShowDocumentBadRequest) WithPayload(payload *internalmessages.InvalidRequestResponsePayload) *ShowDocumentBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the show document bad request response
func (o *ShowDocumentBadRequest) SetPayload(payload *internalmessages.InvalidRequestResponsePayload) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ShowDocumentBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ShowDocumentForbiddenCode is the HTTP code returned for type ShowDocumentForbidden
const ShowDocumentForbiddenCode int = 403

/*
ShowDocumentForbidden not authorized

swagger:response showDocumentForbidden
*/
type ShowDocumentForbidden struct {
}

// NewShowDocumentForbidden creates ShowDocumentForbidden with default headers values
func NewShowDocumentForbidden() *ShowDocumentForbidden {

	return &ShowDocumentForbidden{}
}

// WriteResponse to the client
func (o *ShowDocumentForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// ShowDocumentNotFoundCode is the HTTP code returned for type ShowDocumentNotFound
const ShowDocumentNotFoundCode int = 404

/*
ShowDocumentNotFound not found

swagger:response showDocumentNotFound
*/
type ShowDocumentNotFound struct {
}

// NewShowDocumentNotFound creates ShowDocumentNotFound with default headers values
func NewShowDocumentNotFound() *ShowDocumentNotFound {

	return &ShowDocumentNotFound{}
}

// WriteResponse to the client
func (o *ShowDocumentNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// ShowDocumentInternalServerErrorCode is the HTTP code returned for type ShowDocumentInternalServerError
const ShowDocumentInternalServerErrorCode int = 500

/*
ShowDocumentInternalServerError server error

swagger:response showDocumentInternalServerError
*/
type ShowDocumentInternalServerError struct {
}

// NewShowDocumentInternalServerError creates ShowDocumentInternalServerError with default headers values
func NewShowDocumentInternalServerError() *ShowDocumentInternalServerError {

	return &ShowDocumentInternalServerError{}
}

// WriteResponse to the client
func (o *ShowDocumentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
