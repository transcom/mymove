// Code generated by go-swagger; DO NOT EDIT.

package shipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// RequestShipmentCancellationOKCode is the HTTP code returned for type RequestShipmentCancellationOK
const RequestShipmentCancellationOKCode int = 200

/*
RequestShipmentCancellationOK Successfully requested the shipment cancellation

swagger:response requestShipmentCancellationOK
*/
type RequestShipmentCancellationOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.MTOShipment `json:"body,omitempty"`
}

// NewRequestShipmentCancellationOK creates RequestShipmentCancellationOK with default headers values
func NewRequestShipmentCancellationOK() *RequestShipmentCancellationOK {

	return &RequestShipmentCancellationOK{}
}

// WithPayload adds the payload to the request shipment cancellation o k response
func (o *RequestShipmentCancellationOK) WithPayload(payload *ghcmessages.MTOShipment) *RequestShipmentCancellationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation o k response
func (o *RequestShipmentCancellationOK) SetPayload(payload *ghcmessages.MTOShipment) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationForbiddenCode is the HTTP code returned for type RequestShipmentCancellationForbidden
const RequestShipmentCancellationForbiddenCode int = 403

/*
RequestShipmentCancellationForbidden The request was denied

swagger:response requestShipmentCancellationForbidden
*/
type RequestShipmentCancellationForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewRequestShipmentCancellationForbidden creates RequestShipmentCancellationForbidden with default headers values
func NewRequestShipmentCancellationForbidden() *RequestShipmentCancellationForbidden {

	return &RequestShipmentCancellationForbidden{}
}

// WithPayload adds the payload to the request shipment cancellation forbidden response
func (o *RequestShipmentCancellationForbidden) WithPayload(payload *ghcmessages.Error) *RequestShipmentCancellationForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation forbidden response
func (o *RequestShipmentCancellationForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationNotFoundCode is the HTTP code returned for type RequestShipmentCancellationNotFound
const RequestShipmentCancellationNotFoundCode int = 404

/*
RequestShipmentCancellationNotFound The requested resource wasn't found

swagger:response requestShipmentCancellationNotFound
*/
type RequestShipmentCancellationNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewRequestShipmentCancellationNotFound creates RequestShipmentCancellationNotFound with default headers values
func NewRequestShipmentCancellationNotFound() *RequestShipmentCancellationNotFound {

	return &RequestShipmentCancellationNotFound{}
}

// WithPayload adds the payload to the request shipment cancellation not found response
func (o *RequestShipmentCancellationNotFound) WithPayload(payload *ghcmessages.Error) *RequestShipmentCancellationNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation not found response
func (o *RequestShipmentCancellationNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationConflictCode is the HTTP code returned for type RequestShipmentCancellationConflict
const RequestShipmentCancellationConflictCode int = 409

/*
RequestShipmentCancellationConflict Conflict error

swagger:response requestShipmentCancellationConflict
*/
type RequestShipmentCancellationConflict struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewRequestShipmentCancellationConflict creates RequestShipmentCancellationConflict with default headers values
func NewRequestShipmentCancellationConflict() *RequestShipmentCancellationConflict {

	return &RequestShipmentCancellationConflict{}
}

// WithPayload adds the payload to the request shipment cancellation conflict response
func (o *RequestShipmentCancellationConflict) WithPayload(payload *ghcmessages.Error) *RequestShipmentCancellationConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation conflict response
func (o *RequestShipmentCancellationConflict) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationPreconditionFailedCode is the HTTP code returned for type RequestShipmentCancellationPreconditionFailed
const RequestShipmentCancellationPreconditionFailedCode int = 412

/*
RequestShipmentCancellationPreconditionFailed Precondition failed

swagger:response requestShipmentCancellationPreconditionFailed
*/
type RequestShipmentCancellationPreconditionFailed struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewRequestShipmentCancellationPreconditionFailed creates RequestShipmentCancellationPreconditionFailed with default headers values
func NewRequestShipmentCancellationPreconditionFailed() *RequestShipmentCancellationPreconditionFailed {

	return &RequestShipmentCancellationPreconditionFailed{}
}

// WithPayload adds the payload to the request shipment cancellation precondition failed response
func (o *RequestShipmentCancellationPreconditionFailed) WithPayload(payload *ghcmessages.Error) *RequestShipmentCancellationPreconditionFailed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation precondition failed response
func (o *RequestShipmentCancellationPreconditionFailed) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationPreconditionFailed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(412)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationUnprocessableEntityCode is the HTTP code returned for type RequestShipmentCancellationUnprocessableEntity
const RequestShipmentCancellationUnprocessableEntityCode int = 422

/*
RequestShipmentCancellationUnprocessableEntity The payload was unprocessable.

swagger:response requestShipmentCancellationUnprocessableEntity
*/
type RequestShipmentCancellationUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewRequestShipmentCancellationUnprocessableEntity creates RequestShipmentCancellationUnprocessableEntity with default headers values
func NewRequestShipmentCancellationUnprocessableEntity() *RequestShipmentCancellationUnprocessableEntity {

	return &RequestShipmentCancellationUnprocessableEntity{}
}

// WithPayload adds the payload to the request shipment cancellation unprocessable entity response
func (o *RequestShipmentCancellationUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *RequestShipmentCancellationUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation unprocessable entity response
func (o *RequestShipmentCancellationUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RequestShipmentCancellationInternalServerErrorCode is the HTTP code returned for type RequestShipmentCancellationInternalServerError
const RequestShipmentCancellationInternalServerErrorCode int = 500

/*
RequestShipmentCancellationInternalServerError A server error occurred

swagger:response requestShipmentCancellationInternalServerError
*/
type RequestShipmentCancellationInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewRequestShipmentCancellationInternalServerError creates RequestShipmentCancellationInternalServerError with default headers values
func NewRequestShipmentCancellationInternalServerError() *RequestShipmentCancellationInternalServerError {

	return &RequestShipmentCancellationInternalServerError{}
}

// WithPayload adds the payload to the request shipment cancellation internal server error response
func (o *RequestShipmentCancellationInternalServerError) WithPayload(payload *ghcmessages.Error) *RequestShipmentCancellationInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request shipment cancellation internal server error response
func (o *RequestShipmentCancellationInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestShipmentCancellationInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
