// Code generated by go-swagger; DO NOT EDIT.

package queues

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetServicesCounselingOriginListOKCode is the HTTP code returned for type GetServicesCounselingOriginListOK
const GetServicesCounselingOriginListOKCode int = 200

/*
GetServicesCounselingOriginListOK Successfully returned all moves matching the criteria

swagger:response getServicesCounselingOriginListOK
*/
type GetServicesCounselingOriginListOK struct {

	/*
	  In: Body
	*/
	Payload ghcmessages.Locations `json:"body,omitempty"`
}

// NewGetServicesCounselingOriginListOK creates GetServicesCounselingOriginListOK with default headers values
func NewGetServicesCounselingOriginListOK() *GetServicesCounselingOriginListOK {

	return &GetServicesCounselingOriginListOK{}
}

// WithPayload adds the payload to the get services counseling origin list o k response
func (o *GetServicesCounselingOriginListOK) WithPayload(payload ghcmessages.Locations) *GetServicesCounselingOriginListOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services counseling origin list o k response
func (o *GetServicesCounselingOriginListOK) SetPayload(payload ghcmessages.Locations) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesCounselingOriginListOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = ghcmessages.Locations{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetServicesCounselingOriginListForbiddenCode is the HTTP code returned for type GetServicesCounselingOriginListForbidden
const GetServicesCounselingOriginListForbiddenCode int = 403

/*
GetServicesCounselingOriginListForbidden The request was denied

swagger:response getServicesCounselingOriginListForbidden
*/
type GetServicesCounselingOriginListForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetServicesCounselingOriginListForbidden creates GetServicesCounselingOriginListForbidden with default headers values
func NewGetServicesCounselingOriginListForbidden() *GetServicesCounselingOriginListForbidden {

	return &GetServicesCounselingOriginListForbidden{}
}

// WithPayload adds the payload to the get services counseling origin list forbidden response
func (o *GetServicesCounselingOriginListForbidden) WithPayload(payload *ghcmessages.Error) *GetServicesCounselingOriginListForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services counseling origin list forbidden response
func (o *GetServicesCounselingOriginListForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesCounselingOriginListForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetServicesCounselingOriginListInternalServerErrorCode is the HTTP code returned for type GetServicesCounselingOriginListInternalServerError
const GetServicesCounselingOriginListInternalServerErrorCode int = 500

/*
GetServicesCounselingOriginListInternalServerError A server error occurred

swagger:response getServicesCounselingOriginListInternalServerError
*/
type GetServicesCounselingOriginListInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetServicesCounselingOriginListInternalServerError creates GetServicesCounselingOriginListInternalServerError with default headers values
func NewGetServicesCounselingOriginListInternalServerError() *GetServicesCounselingOriginListInternalServerError {

	return &GetServicesCounselingOriginListInternalServerError{}
}

// WithPayload adds the payload to the get services counseling origin list internal server error response
func (o *GetServicesCounselingOriginListInternalServerError) WithPayload(payload *ghcmessages.Error) *GetServicesCounselingOriginListInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services counseling origin list internal server error response
func (o *GetServicesCounselingOriginListInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesCounselingOriginListInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
