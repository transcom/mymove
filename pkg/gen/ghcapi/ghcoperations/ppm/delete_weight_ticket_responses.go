// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// DeleteWeightTicketNoContentCode is the HTTP code returned for type DeleteWeightTicketNoContent
const DeleteWeightTicketNoContentCode int = 204

/*
DeleteWeightTicketNoContent Successfully soft deleted the weight ticket

swagger:response deleteWeightTicketNoContent
*/
type DeleteWeightTicketNoContent struct {
}

// NewDeleteWeightTicketNoContent creates DeleteWeightTicketNoContent with default headers values
func NewDeleteWeightTicketNoContent() *DeleteWeightTicketNoContent {

	return &DeleteWeightTicketNoContent{}
}

// WriteResponse to the client
func (o *DeleteWeightTicketNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteWeightTicketBadRequestCode is the HTTP code returned for type DeleteWeightTicketBadRequest
const DeleteWeightTicketBadRequestCode int = 400

/*
DeleteWeightTicketBadRequest The request payload is invalid

swagger:response deleteWeightTicketBadRequest
*/
type DeleteWeightTicketBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketBadRequest creates DeleteWeightTicketBadRequest with default headers values
func NewDeleteWeightTicketBadRequest() *DeleteWeightTicketBadRequest {

	return &DeleteWeightTicketBadRequest{}
}

// WithPayload adds the payload to the delete weight ticket bad request response
func (o *DeleteWeightTicketBadRequest) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket bad request response
func (o *DeleteWeightTicketBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketUnauthorizedCode is the HTTP code returned for type DeleteWeightTicketUnauthorized
const DeleteWeightTicketUnauthorizedCode int = 401

/*
DeleteWeightTicketUnauthorized The request was denied

swagger:response deleteWeightTicketUnauthorized
*/
type DeleteWeightTicketUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketUnauthorized creates DeleteWeightTicketUnauthorized with default headers values
func NewDeleteWeightTicketUnauthorized() *DeleteWeightTicketUnauthorized {

	return &DeleteWeightTicketUnauthorized{}
}

// WithPayload adds the payload to the delete weight ticket unauthorized response
func (o *DeleteWeightTicketUnauthorized) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket unauthorized response
func (o *DeleteWeightTicketUnauthorized) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketForbiddenCode is the HTTP code returned for type DeleteWeightTicketForbidden
const DeleteWeightTicketForbiddenCode int = 403

/*
DeleteWeightTicketForbidden The request was denied

swagger:response deleteWeightTicketForbidden
*/
type DeleteWeightTicketForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketForbidden creates DeleteWeightTicketForbidden with default headers values
func NewDeleteWeightTicketForbidden() *DeleteWeightTicketForbidden {

	return &DeleteWeightTicketForbidden{}
}

// WithPayload adds the payload to the delete weight ticket forbidden response
func (o *DeleteWeightTicketForbidden) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket forbidden response
func (o *DeleteWeightTicketForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketNotFoundCode is the HTTP code returned for type DeleteWeightTicketNotFound
const DeleteWeightTicketNotFoundCode int = 404

/*
DeleteWeightTicketNotFound The requested resource wasn't found

swagger:response deleteWeightTicketNotFound
*/
type DeleteWeightTicketNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketNotFound creates DeleteWeightTicketNotFound with default headers values
func NewDeleteWeightTicketNotFound() *DeleteWeightTicketNotFound {

	return &DeleteWeightTicketNotFound{}
}

// WithPayload adds the payload to the delete weight ticket not found response
func (o *DeleteWeightTicketNotFound) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket not found response
func (o *DeleteWeightTicketNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketConflictCode is the HTTP code returned for type DeleteWeightTicketConflict
const DeleteWeightTicketConflictCode int = 409

/*
DeleteWeightTicketConflict Conflict error

swagger:response deleteWeightTicketConflict
*/
type DeleteWeightTicketConflict struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketConflict creates DeleteWeightTicketConflict with default headers values
func NewDeleteWeightTicketConflict() *DeleteWeightTicketConflict {

	return &DeleteWeightTicketConflict{}
}

// WithPayload adds the payload to the delete weight ticket conflict response
func (o *DeleteWeightTicketConflict) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket conflict response
func (o *DeleteWeightTicketConflict) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketUnprocessableEntityCode is the HTTP code returned for type DeleteWeightTicketUnprocessableEntity
const DeleteWeightTicketUnprocessableEntityCode int = 422

/*
DeleteWeightTicketUnprocessableEntity The payload was unprocessable.

swagger:response deleteWeightTicketUnprocessableEntity
*/
type DeleteWeightTicketUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewDeleteWeightTicketUnprocessableEntity creates DeleteWeightTicketUnprocessableEntity with default headers values
func NewDeleteWeightTicketUnprocessableEntity() *DeleteWeightTicketUnprocessableEntity {

	return &DeleteWeightTicketUnprocessableEntity{}
}

// WithPayload adds the payload to the delete weight ticket unprocessable entity response
func (o *DeleteWeightTicketUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *DeleteWeightTicketUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket unprocessable entity response
func (o *DeleteWeightTicketUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteWeightTicketInternalServerErrorCode is the HTTP code returned for type DeleteWeightTicketInternalServerError
const DeleteWeightTicketInternalServerErrorCode int = 500

/*
DeleteWeightTicketInternalServerError A server error occurred

swagger:response deleteWeightTicketInternalServerError
*/
type DeleteWeightTicketInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewDeleteWeightTicketInternalServerError creates DeleteWeightTicketInternalServerError with default headers values
func NewDeleteWeightTicketInternalServerError() *DeleteWeightTicketInternalServerError {

	return &DeleteWeightTicketInternalServerError{}
}

// WithPayload adds the payload to the delete weight ticket internal server error response
func (o *DeleteWeightTicketInternalServerError) WithPayload(payload *ghcmessages.Error) *DeleteWeightTicketInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete weight ticket internal server error response
func (o *DeleteWeightTicketInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteWeightTicketInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
