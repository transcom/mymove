// Code generated by go-swagger; DO NOT EDIT.

package customer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetCustomerOKCode is the HTTP code returned for type GetCustomerOK
const GetCustomerOKCode int = 200

/*
GetCustomerOK Successfully retrieved information on an individual customer

swagger:response getCustomerOK
*/
type GetCustomerOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Customer `json:"body,omitempty"`
}

// NewGetCustomerOK creates GetCustomerOK with default headers values
func NewGetCustomerOK() *GetCustomerOK {

	return &GetCustomerOK{}
}

// WithPayload adds the payload to the get customer o k response
func (o *GetCustomerOK) WithPayload(payload *ghcmessages.Customer) *GetCustomerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer o k response
func (o *GetCustomerOK) SetPayload(payload *ghcmessages.Customer) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerBadRequestCode is the HTTP code returned for type GetCustomerBadRequest
const GetCustomerBadRequestCode int = 400

/*
GetCustomerBadRequest The request payload is invalid

swagger:response getCustomerBadRequest
*/
type GetCustomerBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerBadRequest creates GetCustomerBadRequest with default headers values
func NewGetCustomerBadRequest() *GetCustomerBadRequest {

	return &GetCustomerBadRequest{}
}

// WithPayload adds the payload to the get customer bad request response
func (o *GetCustomerBadRequest) WithPayload(payload *ghcmessages.Error) *GetCustomerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer bad request response
func (o *GetCustomerBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerUnauthorizedCode is the HTTP code returned for type GetCustomerUnauthorized
const GetCustomerUnauthorizedCode int = 401

/*
GetCustomerUnauthorized The request was denied

swagger:response getCustomerUnauthorized
*/
type GetCustomerUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerUnauthorized creates GetCustomerUnauthorized with default headers values
func NewGetCustomerUnauthorized() *GetCustomerUnauthorized {

	return &GetCustomerUnauthorized{}
}

// WithPayload adds the payload to the get customer unauthorized response
func (o *GetCustomerUnauthorized) WithPayload(payload *ghcmessages.Error) *GetCustomerUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer unauthorized response
func (o *GetCustomerUnauthorized) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerForbiddenCode is the HTTP code returned for type GetCustomerForbidden
const GetCustomerForbiddenCode int = 403

/*
GetCustomerForbidden The request was denied

swagger:response getCustomerForbidden
*/
type GetCustomerForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerForbidden creates GetCustomerForbidden with default headers values
func NewGetCustomerForbidden() *GetCustomerForbidden {

	return &GetCustomerForbidden{}
}

// WithPayload adds the payload to the get customer forbidden response
func (o *GetCustomerForbidden) WithPayload(payload *ghcmessages.Error) *GetCustomerForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer forbidden response
func (o *GetCustomerForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerNotFoundCode is the HTTP code returned for type GetCustomerNotFound
const GetCustomerNotFoundCode int = 404

/*
GetCustomerNotFound The requested resource wasn't found

swagger:response getCustomerNotFound
*/
type GetCustomerNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerNotFound creates GetCustomerNotFound with default headers values
func NewGetCustomerNotFound() *GetCustomerNotFound {

	return &GetCustomerNotFound{}
}

// WithPayload adds the payload to the get customer not found response
func (o *GetCustomerNotFound) WithPayload(payload *ghcmessages.Error) *GetCustomerNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer not found response
func (o *GetCustomerNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerInternalServerErrorCode is the HTTP code returned for type GetCustomerInternalServerError
const GetCustomerInternalServerErrorCode int = 500

/*
GetCustomerInternalServerError A server error occurred

swagger:response getCustomerInternalServerError
*/
type GetCustomerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerInternalServerError creates GetCustomerInternalServerError with default headers values
func NewGetCustomerInternalServerError() *GetCustomerInternalServerError {

	return &GetCustomerInternalServerError{}
}

// WithPayload adds the payload to the get customer internal server error response
func (o *GetCustomerInternalServerError) WithPayload(payload *ghcmessages.Error) *GetCustomerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer internal server error response
func (o *GetCustomerInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
