// Code generated by go-swagger; DO NOT EDIT.

package addresses

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// GetLocationByZipCityStateOKCode is the HTTP code returned for type GetLocationByZipCityStateOK
const GetLocationByZipCityStateOKCode int = 200

/*
GetLocationByZipCityStateOK the requested list of city, state, county, and postal code matches

swagger:response getLocationByZipCityStateOK
*/
type GetLocationByZipCityStateOK struct {

	/*
	  In: Body
	*/
	Payload internalmessages.VLocations `json:"body,omitempty"`
}

// NewGetLocationByZipCityStateOK creates GetLocationByZipCityStateOK with default headers values
func NewGetLocationByZipCityStateOK() *GetLocationByZipCityStateOK {

	return &GetLocationByZipCityStateOK{}
}

// WithPayload adds the payload to the get location by zip city state o k response
func (o *GetLocationByZipCityStateOK) WithPayload(payload internalmessages.VLocations) *GetLocationByZipCityStateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get location by zip city state o k response
func (o *GetLocationByZipCityStateOK) SetPayload(payload internalmessages.VLocations) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLocationByZipCityStateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = internalmessages.VLocations{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetLocationByZipCityStateBadRequestCode is the HTTP code returned for type GetLocationByZipCityStateBadRequest
const GetLocationByZipCityStateBadRequestCode int = 400

/*
GetLocationByZipCityStateBadRequest The request payload is invalid.

swagger:response getLocationByZipCityStateBadRequest
*/
type GetLocationByZipCityStateBadRequest struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewGetLocationByZipCityStateBadRequest creates GetLocationByZipCityStateBadRequest with default headers values
func NewGetLocationByZipCityStateBadRequest() *GetLocationByZipCityStateBadRequest {

	return &GetLocationByZipCityStateBadRequest{}
}

// WithPayload adds the payload to the get location by zip city state bad request response
func (o *GetLocationByZipCityStateBadRequest) WithPayload(payload *internalmessages.ClientError) *GetLocationByZipCityStateBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get location by zip city state bad request response
func (o *GetLocationByZipCityStateBadRequest) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLocationByZipCityStateBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetLocationByZipCityStateForbiddenCode is the HTTP code returned for type GetLocationByZipCityStateForbidden
const GetLocationByZipCityStateForbiddenCode int = 403

/*
GetLocationByZipCityStateForbidden The request was denied.

swagger:response getLocationByZipCityStateForbidden
*/
type GetLocationByZipCityStateForbidden struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewGetLocationByZipCityStateForbidden creates GetLocationByZipCityStateForbidden with default headers values
func NewGetLocationByZipCityStateForbidden() *GetLocationByZipCityStateForbidden {

	return &GetLocationByZipCityStateForbidden{}
}

// WithPayload adds the payload to the get location by zip city state forbidden response
func (o *GetLocationByZipCityStateForbidden) WithPayload(payload *internalmessages.ClientError) *GetLocationByZipCityStateForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get location by zip city state forbidden response
func (o *GetLocationByZipCityStateForbidden) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLocationByZipCityStateForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetLocationByZipCityStateNotFoundCode is the HTTP code returned for type GetLocationByZipCityStateNotFound
const GetLocationByZipCityStateNotFoundCode int = 404

/*
GetLocationByZipCityStateNotFound The requested resource wasn't found.

swagger:response getLocationByZipCityStateNotFound
*/
type GetLocationByZipCityStateNotFound struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewGetLocationByZipCityStateNotFound creates GetLocationByZipCityStateNotFound with default headers values
func NewGetLocationByZipCityStateNotFound() *GetLocationByZipCityStateNotFound {

	return &GetLocationByZipCityStateNotFound{}
}

// WithPayload adds the payload to the get location by zip city state not found response
func (o *GetLocationByZipCityStateNotFound) WithPayload(payload *internalmessages.ClientError) *GetLocationByZipCityStateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get location by zip city state not found response
func (o *GetLocationByZipCityStateNotFound) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLocationByZipCityStateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetLocationByZipCityStateInternalServerErrorCode is the HTTP code returned for type GetLocationByZipCityStateInternalServerError
const GetLocationByZipCityStateInternalServerErrorCode int = 500

/*
GetLocationByZipCityStateInternalServerError A server error occurred.

swagger:response getLocationByZipCityStateInternalServerError
*/
type GetLocationByZipCityStateInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.Error `json:"body,omitempty"`
}

// NewGetLocationByZipCityStateInternalServerError creates GetLocationByZipCityStateInternalServerError with default headers values
func NewGetLocationByZipCityStateInternalServerError() *GetLocationByZipCityStateInternalServerError {

	return &GetLocationByZipCityStateInternalServerError{}
}

// WithPayload adds the payload to the get location by zip city state internal server error response
func (o *GetLocationByZipCityStateInternalServerError) WithPayload(payload *internalmessages.Error) *GetLocationByZipCityStateInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get location by zip city state internal server error response
func (o *GetLocationByZipCityStateInternalServerError) SetPayload(payload *internalmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLocationByZipCityStateInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}