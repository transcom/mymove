// Code generated by go-swagger; DO NOT EDIT.

package order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskOKCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskOKCode int = 200

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK updated Move

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskOK
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Move `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskOK creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskOK() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk o k response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK) WithPayload(payload *ghcmessages.Move) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk o k response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK) SetPayload(payload *ghcmessages.Move) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbiddenCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbiddenCode int = 403

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden The request was denied

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk forbidden response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden) WithPayload(payload *ghcmessages.Error) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk forbidden response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFoundCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFoundCode int = 404

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound The requested resource wasn't found

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk not found response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound) WithPayload(payload *ghcmessages.Error) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk not found response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailedCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailedCode int = 412

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed Precondition failed

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk precondition failed response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed) WithPayload(payload *ghcmessages.Error) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk precondition failed response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskPreconditionFailed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(412)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntityCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntityCode int = 422

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity The payload was unprocessable.

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk unprocessable entity response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk unprocessable entity response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerErrorCode is the HTTP code returned for type AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError
const AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerErrorCode int = 500

/*
AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError A server error occurred

swagger:response acknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError
*/
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError creates AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError with default headers values
func NewAcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError() *AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError {

	return &AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError{}
}

// WithPayload adds the payload to the acknowledge excess unaccompanied baggage weight risk internal server error response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError) WithPayload(payload *ghcmessages.Error) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the acknowledge excess unaccompanied baggage weight risk internal server error response
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}