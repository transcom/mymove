// Code generated by go-swagger; DO NOT EDIT.

package re_service_items

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetAllReServiceItemsOKCode is the HTTP code returned for type GetAllReServiceItemsOK
const GetAllReServiceItemsOKCode int = 200

/*
GetAllReServiceItemsOK Successfully retrieved all ReServiceItems.

swagger:response getAllReServiceItemsOK
*/
type GetAllReServiceItemsOK struct {

	/*
	  In: Body
	*/
	Payload ghcmessages.ReServiceItems `json:"body,omitempty"`
}

// NewGetAllReServiceItemsOK creates GetAllReServiceItemsOK with default headers values
func NewGetAllReServiceItemsOK() *GetAllReServiceItemsOK {

	return &GetAllReServiceItemsOK{}
}

// WithPayload adds the payload to the get all re service items o k response
func (o *GetAllReServiceItemsOK) WithPayload(payload ghcmessages.ReServiceItems) *GetAllReServiceItemsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all re service items o k response
func (o *GetAllReServiceItemsOK) SetPayload(payload ghcmessages.ReServiceItems) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReServiceItemsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = ghcmessages.ReServiceItems{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetAllReServiceItemsBadRequestCode is the HTTP code returned for type GetAllReServiceItemsBadRequest
const GetAllReServiceItemsBadRequestCode int = 400

/*
GetAllReServiceItemsBadRequest The request payload is invalid

swagger:response getAllReServiceItemsBadRequest
*/
type GetAllReServiceItemsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetAllReServiceItemsBadRequest creates GetAllReServiceItemsBadRequest with default headers values
func NewGetAllReServiceItemsBadRequest() *GetAllReServiceItemsBadRequest {

	return &GetAllReServiceItemsBadRequest{}
}

// WithPayload adds the payload to the get all re service items bad request response
func (o *GetAllReServiceItemsBadRequest) WithPayload(payload *ghcmessages.Error) *GetAllReServiceItemsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all re service items bad request response
func (o *GetAllReServiceItemsBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReServiceItemsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAllReServiceItemsUnauthorizedCode is the HTTP code returned for type GetAllReServiceItemsUnauthorized
const GetAllReServiceItemsUnauthorizedCode int = 401

/*
GetAllReServiceItemsUnauthorized The request was denied

swagger:response getAllReServiceItemsUnauthorized
*/
type GetAllReServiceItemsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetAllReServiceItemsUnauthorized creates GetAllReServiceItemsUnauthorized with default headers values
func NewGetAllReServiceItemsUnauthorized() *GetAllReServiceItemsUnauthorized {

	return &GetAllReServiceItemsUnauthorized{}
}

// WithPayload adds the payload to the get all re service items unauthorized response
func (o *GetAllReServiceItemsUnauthorized) WithPayload(payload *ghcmessages.Error) *GetAllReServiceItemsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all re service items unauthorized response
func (o *GetAllReServiceItemsUnauthorized) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReServiceItemsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAllReServiceItemsNotFoundCode is the HTTP code returned for type GetAllReServiceItemsNotFound
const GetAllReServiceItemsNotFoundCode int = 404

/*
GetAllReServiceItemsNotFound The requested resource wasn't found

swagger:response getAllReServiceItemsNotFound
*/
type GetAllReServiceItemsNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetAllReServiceItemsNotFound creates GetAllReServiceItemsNotFound with default headers values
func NewGetAllReServiceItemsNotFound() *GetAllReServiceItemsNotFound {

	return &GetAllReServiceItemsNotFound{}
}

// WithPayload adds the payload to the get all re service items not found response
func (o *GetAllReServiceItemsNotFound) WithPayload(payload *ghcmessages.Error) *GetAllReServiceItemsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all re service items not found response
func (o *GetAllReServiceItemsNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReServiceItemsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAllReServiceItemsInternalServerErrorCode is the HTTP code returned for type GetAllReServiceItemsInternalServerError
const GetAllReServiceItemsInternalServerErrorCode int = 500

/*
GetAllReServiceItemsInternalServerError A server error occurred

swagger:response getAllReServiceItemsInternalServerError
*/
type GetAllReServiceItemsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetAllReServiceItemsInternalServerError creates GetAllReServiceItemsInternalServerError with default headers values
func NewGetAllReServiceItemsInternalServerError() *GetAllReServiceItemsInternalServerError {

	return &GetAllReServiceItemsInternalServerError{}
}

// WithPayload adds the payload to the get all re service items internal server error response
func (o *GetAllReServiceItemsInternalServerError) WithPayload(payload *ghcmessages.Error) *GetAllReServiceItemsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all re service items internal server error response
func (o *GetAllReServiceItemsInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReServiceItemsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
