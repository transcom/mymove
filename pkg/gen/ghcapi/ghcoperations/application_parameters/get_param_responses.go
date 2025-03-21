// Code generated by go-swagger; DO NOT EDIT.

package application_parameters

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetParamOKCode is the HTTP code returned for type GetParamOK
const GetParamOKCode int = 200

/*
GetParamOK Application Parameters

swagger:response getParamOK
*/
type GetParamOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ApplicationParameters `json:"body,omitempty"`
}

// NewGetParamOK creates GetParamOK with default headers values
func NewGetParamOK() *GetParamOK {

	return &GetParamOK{}
}

// WithPayload adds the payload to the get param o k response
func (o *GetParamOK) WithPayload(payload *ghcmessages.ApplicationParameters) *GetParamOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get param o k response
func (o *GetParamOK) SetPayload(payload *ghcmessages.ApplicationParameters) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetParamOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetParamBadRequestCode is the HTTP code returned for type GetParamBadRequest
const GetParamBadRequestCode int = 400

/*
GetParamBadRequest invalid request

swagger:response getParamBadRequest
*/
type GetParamBadRequest struct {
}

// NewGetParamBadRequest creates GetParamBadRequest with default headers values
func NewGetParamBadRequest() *GetParamBadRequest {

	return &GetParamBadRequest{}
}

// WriteResponse to the client
func (o *GetParamBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetParamUnauthorizedCode is the HTTP code returned for type GetParamUnauthorized
const GetParamUnauthorizedCode int = 401

/*
GetParamUnauthorized request requires user authentication

swagger:response getParamUnauthorized
*/
type GetParamUnauthorized struct {
}

// NewGetParamUnauthorized creates GetParamUnauthorized with default headers values
func NewGetParamUnauthorized() *GetParamUnauthorized {

	return &GetParamUnauthorized{}
}

// WriteResponse to the client
func (o *GetParamUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetParamInternalServerErrorCode is the HTTP code returned for type GetParamInternalServerError
const GetParamInternalServerErrorCode int = 500

/*
GetParamInternalServerError server error

swagger:response getParamInternalServerError
*/
type GetParamInternalServerError struct {
}

// NewGetParamInternalServerError creates GetParamInternalServerError with default headers values
func NewGetParamInternalServerError() *GetParamInternalServerError {

	return &GetParamInternalServerError{}
}

// WriteResponse to the client
func (o *GetParamInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
