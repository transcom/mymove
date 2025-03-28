// Code generated by go-swagger; DO NOT EDIT.

package customer_support_remarks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// GetCustomerSupportRemarksForMoveOKCode is the HTTP code returned for type GetCustomerSupportRemarksForMoveOK
const GetCustomerSupportRemarksForMoveOKCode int = 200

/*
GetCustomerSupportRemarksForMoveOK Successfully retrieved all line items for a move task order

swagger:response getCustomerSupportRemarksForMoveOK
*/
type GetCustomerSupportRemarksForMoveOK struct {

	/*
	  In: Body
	*/
	Payload ghcmessages.CustomerSupportRemarks `json:"body,omitempty"`
}

// NewGetCustomerSupportRemarksForMoveOK creates GetCustomerSupportRemarksForMoveOK with default headers values
func NewGetCustomerSupportRemarksForMoveOK() *GetCustomerSupportRemarksForMoveOK {

	return &GetCustomerSupportRemarksForMoveOK{}
}

// WithPayload adds the payload to the get customer support remarks for move o k response
func (o *GetCustomerSupportRemarksForMoveOK) WithPayload(payload ghcmessages.CustomerSupportRemarks) *GetCustomerSupportRemarksForMoveOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer support remarks for move o k response
func (o *GetCustomerSupportRemarksForMoveOK) SetPayload(payload ghcmessages.CustomerSupportRemarks) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerSupportRemarksForMoveOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = ghcmessages.CustomerSupportRemarks{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetCustomerSupportRemarksForMoveForbiddenCode is the HTTP code returned for type GetCustomerSupportRemarksForMoveForbidden
const GetCustomerSupportRemarksForMoveForbiddenCode int = 403

/*
GetCustomerSupportRemarksForMoveForbidden The request was denied

swagger:response getCustomerSupportRemarksForMoveForbidden
*/
type GetCustomerSupportRemarksForMoveForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerSupportRemarksForMoveForbidden creates GetCustomerSupportRemarksForMoveForbidden with default headers values
func NewGetCustomerSupportRemarksForMoveForbidden() *GetCustomerSupportRemarksForMoveForbidden {

	return &GetCustomerSupportRemarksForMoveForbidden{}
}

// WithPayload adds the payload to the get customer support remarks for move forbidden response
func (o *GetCustomerSupportRemarksForMoveForbidden) WithPayload(payload *ghcmessages.Error) *GetCustomerSupportRemarksForMoveForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer support remarks for move forbidden response
func (o *GetCustomerSupportRemarksForMoveForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerSupportRemarksForMoveForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerSupportRemarksForMoveNotFoundCode is the HTTP code returned for type GetCustomerSupportRemarksForMoveNotFound
const GetCustomerSupportRemarksForMoveNotFoundCode int = 404

/*
GetCustomerSupportRemarksForMoveNotFound The requested resource wasn't found

swagger:response getCustomerSupportRemarksForMoveNotFound
*/
type GetCustomerSupportRemarksForMoveNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerSupportRemarksForMoveNotFound creates GetCustomerSupportRemarksForMoveNotFound with default headers values
func NewGetCustomerSupportRemarksForMoveNotFound() *GetCustomerSupportRemarksForMoveNotFound {

	return &GetCustomerSupportRemarksForMoveNotFound{}
}

// WithPayload adds the payload to the get customer support remarks for move not found response
func (o *GetCustomerSupportRemarksForMoveNotFound) WithPayload(payload *ghcmessages.Error) *GetCustomerSupportRemarksForMoveNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer support remarks for move not found response
func (o *GetCustomerSupportRemarksForMoveNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerSupportRemarksForMoveNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerSupportRemarksForMoveUnprocessableEntityCode is the HTTP code returned for type GetCustomerSupportRemarksForMoveUnprocessableEntity
const GetCustomerSupportRemarksForMoveUnprocessableEntityCode int = 422

/*
GetCustomerSupportRemarksForMoveUnprocessableEntity The payload was unprocessable.

swagger:response getCustomerSupportRemarksForMoveUnprocessableEntity
*/
type GetCustomerSupportRemarksForMoveUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewGetCustomerSupportRemarksForMoveUnprocessableEntity creates GetCustomerSupportRemarksForMoveUnprocessableEntity with default headers values
func NewGetCustomerSupportRemarksForMoveUnprocessableEntity() *GetCustomerSupportRemarksForMoveUnprocessableEntity {

	return &GetCustomerSupportRemarksForMoveUnprocessableEntity{}
}

// WithPayload adds the payload to the get customer support remarks for move unprocessable entity response
func (o *GetCustomerSupportRemarksForMoveUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *GetCustomerSupportRemarksForMoveUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer support remarks for move unprocessable entity response
func (o *GetCustomerSupportRemarksForMoveUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerSupportRemarksForMoveUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCustomerSupportRemarksForMoveInternalServerErrorCode is the HTTP code returned for type GetCustomerSupportRemarksForMoveInternalServerError
const GetCustomerSupportRemarksForMoveInternalServerErrorCode int = 500

/*
GetCustomerSupportRemarksForMoveInternalServerError A server error occurred

swagger:response getCustomerSupportRemarksForMoveInternalServerError
*/
type GetCustomerSupportRemarksForMoveInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewGetCustomerSupportRemarksForMoveInternalServerError creates GetCustomerSupportRemarksForMoveInternalServerError with default headers values
func NewGetCustomerSupportRemarksForMoveInternalServerError() *GetCustomerSupportRemarksForMoveInternalServerError {

	return &GetCustomerSupportRemarksForMoveInternalServerError{}
}

// WithPayload adds the payload to the get customer support remarks for move internal server error response
func (o *GetCustomerSupportRemarksForMoveInternalServerError) WithPayload(payload *ghcmessages.Error) *GetCustomerSupportRemarksForMoveInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer support remarks for move internal server error response
func (o *GetCustomerSupportRemarksForMoveInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerSupportRemarksForMoveInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
