// Code generated by go-swagger; DO NOT EDIT.

package transportation_office

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetTransportationOfficesOKCode is the HTTP code returned for type GetTransportationOfficesOK
const GetTransportationOfficesOKCode int = 200

/*
GetTransportationOfficesOK Successfully retrieved transportation offices

swagger:response getTransportationOfficesOK
*/
type GetTransportationOfficesOK struct {

	/*
	  In: Body
	*/
	Payload ghcmessages.TransportationOffices `json:"body,omitempty"`
}

// NewGetTransportationOfficesOK creates GetTransportationOfficesOK with default headers values
func NewGetTransportationOfficesOK() *GetTransportationOfficesOK {

	return &GetTransportationOfficesOK{}
}

// WithPayload adds the payload to the get transportation offices o k response
func (o *GetTransportationOfficesOK) WithPayload(payload ghcmessages.TransportationOffices) *GetTransportationOfficesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices o k response
func (o *GetTransportationOfficesOK) SetPayload(payload ghcmessages.TransportationOffices) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = ghcmessages.TransportationOffices{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetTransportationOfficesBadRequestCode is the HTTP code returned for type GetTransportationOfficesBadRequest
const GetTransportationOfficesBadRequestCode int = 400

/*
GetTransportationOfficesBadRequest The request payload is invalid

swagger:response getTransportationOfficesBadRequest
*/
type GetTransportationOfficesBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetTransportationOfficesBadRequest creates GetTransportationOfficesBadRequest with default headers values
func NewGetTransportationOfficesBadRequest() *GetTransportationOfficesBadRequest {

	return &GetTransportationOfficesBadRequest{}
}

// WithPayload adds the payload to the get transportation offices bad request response
func (o *GetTransportationOfficesBadRequest) WithPayload(payload *ghcmessages.Error) *GetTransportationOfficesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices bad request response
func (o *GetTransportationOfficesBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTransportationOfficesUnauthorizedCode is the HTTP code returned for type GetTransportationOfficesUnauthorized
const GetTransportationOfficesUnauthorizedCode int = 401

/*
GetTransportationOfficesUnauthorized The request was denied

swagger:response getTransportationOfficesUnauthorized
*/
type GetTransportationOfficesUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetTransportationOfficesUnauthorized creates GetTransportationOfficesUnauthorized with default headers values
func NewGetTransportationOfficesUnauthorized() *GetTransportationOfficesUnauthorized {

	return &GetTransportationOfficesUnauthorized{}
}

// WithPayload adds the payload to the get transportation offices unauthorized response
func (o *GetTransportationOfficesUnauthorized) WithPayload(payload *ghcmessages.Error) *GetTransportationOfficesUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices unauthorized response
func (o *GetTransportationOfficesUnauthorized) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTransportationOfficesForbiddenCode is the HTTP code returned for type GetTransportationOfficesForbidden
const GetTransportationOfficesForbiddenCode int = 403

/*
GetTransportationOfficesForbidden The request was denied

swagger:response getTransportationOfficesForbidden
*/
type GetTransportationOfficesForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetTransportationOfficesForbidden creates GetTransportationOfficesForbidden with default headers values
func NewGetTransportationOfficesForbidden() *GetTransportationOfficesForbidden {

	return &GetTransportationOfficesForbidden{}
}

// WithPayload adds the payload to the get transportation offices forbidden response
func (o *GetTransportationOfficesForbidden) WithPayload(payload *ghcmessages.Error) *GetTransportationOfficesForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices forbidden response
func (o *GetTransportationOfficesForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTransportationOfficesNotFoundCode is the HTTP code returned for type GetTransportationOfficesNotFound
const GetTransportationOfficesNotFoundCode int = 404

/*
GetTransportationOfficesNotFound The requested resource wasn't found

swagger:response getTransportationOfficesNotFound
*/
type GetTransportationOfficesNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetTransportationOfficesNotFound creates GetTransportationOfficesNotFound with default headers values
func NewGetTransportationOfficesNotFound() *GetTransportationOfficesNotFound {

	return &GetTransportationOfficesNotFound{}
}

// WithPayload adds the payload to the get transportation offices not found response
func (o *GetTransportationOfficesNotFound) WithPayload(payload *ghcmessages.Error) *GetTransportationOfficesNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices not found response
func (o *GetTransportationOfficesNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTransportationOfficesInternalServerErrorCode is the HTTP code returned for type GetTransportationOfficesInternalServerError
const GetTransportationOfficesInternalServerErrorCode int = 500

/*
GetTransportationOfficesInternalServerError A server error occurred

swagger:response getTransportationOfficesInternalServerError
*/
type GetTransportationOfficesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetTransportationOfficesInternalServerError creates GetTransportationOfficesInternalServerError with default headers values
func NewGetTransportationOfficesInternalServerError() *GetTransportationOfficesInternalServerError {

	return &GetTransportationOfficesInternalServerError{}
}

// WithPayload adds the payload to the get transportation offices internal server error response
func (o *GetTransportationOfficesInternalServerError) WithPayload(payload *ghcmessages.Error) *GetTransportationOfficesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get transportation offices internal server error response
func (o *GetTransportationOfficesInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTransportationOfficesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
