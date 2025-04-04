// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetPPMDocumentsOKCode is the HTTP code returned for type GetPPMDocumentsOK
const GetPPMDocumentsOKCode int = 200

/*
GetPPMDocumentsOK All PPM documents and associated uploads for the specified PPM shipment.

swagger:response getPPMDocumentsOK
*/
type GetPPMDocumentsOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.PPMDocuments `json:"body,omitempty"`
}

// NewGetPPMDocumentsOK creates GetPPMDocumentsOK with default headers values
func NewGetPPMDocumentsOK() *GetPPMDocumentsOK {

	return &GetPPMDocumentsOK{}
}

// WithPayload adds the payload to the get p p m documents o k response
func (o *GetPPMDocumentsOK) WithPayload(payload *ghcmessages.PPMDocuments) *GetPPMDocumentsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get p p m documents o k response
func (o *GetPPMDocumentsOK) SetPayload(payload *ghcmessages.PPMDocuments) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPPMDocumentsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPPMDocumentsUnauthorizedCode is the HTTP code returned for type GetPPMDocumentsUnauthorized
const GetPPMDocumentsUnauthorizedCode int = 401

/*
GetPPMDocumentsUnauthorized The request was denied

swagger:response getPPMDocumentsUnauthorized
*/
type GetPPMDocumentsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetPPMDocumentsUnauthorized creates GetPPMDocumentsUnauthorized with default headers values
func NewGetPPMDocumentsUnauthorized() *GetPPMDocumentsUnauthorized {

	return &GetPPMDocumentsUnauthorized{}
}

// WithPayload adds the payload to the get p p m documents unauthorized response
func (o *GetPPMDocumentsUnauthorized) WithPayload(payload *ghcmessages.Error) *GetPPMDocumentsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get p p m documents unauthorized response
func (o *GetPPMDocumentsUnauthorized) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPPMDocumentsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPPMDocumentsForbiddenCode is the HTTP code returned for type GetPPMDocumentsForbidden
const GetPPMDocumentsForbiddenCode int = 403

/*
GetPPMDocumentsForbidden The request was denied

swagger:response getPPMDocumentsForbidden
*/
type GetPPMDocumentsForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetPPMDocumentsForbidden creates GetPPMDocumentsForbidden with default headers values
func NewGetPPMDocumentsForbidden() *GetPPMDocumentsForbidden {

	return &GetPPMDocumentsForbidden{}
}

// WithPayload adds the payload to the get p p m documents forbidden response
func (o *GetPPMDocumentsForbidden) WithPayload(payload *ghcmessages.Error) *GetPPMDocumentsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get p p m documents forbidden response
func (o *GetPPMDocumentsForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPPMDocumentsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPPMDocumentsUnprocessableEntityCode is the HTTP code returned for type GetPPMDocumentsUnprocessableEntity
const GetPPMDocumentsUnprocessableEntityCode int = 422

/*
GetPPMDocumentsUnprocessableEntity The payload was unprocessable.

swagger:response getPPMDocumentsUnprocessableEntity
*/
type GetPPMDocumentsUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewGetPPMDocumentsUnprocessableEntity creates GetPPMDocumentsUnprocessableEntity with default headers values
func NewGetPPMDocumentsUnprocessableEntity() *GetPPMDocumentsUnprocessableEntity {

	return &GetPPMDocumentsUnprocessableEntity{}
}

// WithPayload adds the payload to the get p p m documents unprocessable entity response
func (o *GetPPMDocumentsUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *GetPPMDocumentsUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get p p m documents unprocessable entity response
func (o *GetPPMDocumentsUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPPMDocumentsUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPPMDocumentsInternalServerErrorCode is the HTTP code returned for type GetPPMDocumentsInternalServerError
const GetPPMDocumentsInternalServerErrorCode int = 500

/*
GetPPMDocumentsInternalServerError A server error occurred

swagger:response getPPMDocumentsInternalServerError
*/
type GetPPMDocumentsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetPPMDocumentsInternalServerError creates GetPPMDocumentsInternalServerError with default headers values
func NewGetPPMDocumentsInternalServerError() *GetPPMDocumentsInternalServerError {

	return &GetPPMDocumentsInternalServerError{}
}

// WithPayload adds the payload to the get p p m documents internal server error response
func (o *GetPPMDocumentsInternalServerError) WithPayload(payload *ghcmessages.Error) *GetPPMDocumentsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get p p m documents internal server error response
func (o *GetPPMDocumentsInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPPMDocumentsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
