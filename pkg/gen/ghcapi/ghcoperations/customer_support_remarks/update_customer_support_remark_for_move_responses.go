// Code generated by go-swagger; DO NOT EDIT.

package customer_support_remarks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// UpdateCustomerSupportRemarkForMoveOKCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveOK
const UpdateCustomerSupportRemarkForMoveOKCode int = 200

/*
UpdateCustomerSupportRemarkForMoveOK Successfully updated customer support remark

swagger:response updateCustomerSupportRemarkForMoveOK
*/
type UpdateCustomerSupportRemarkForMoveOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.CustomerSupportRemark `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveOK creates UpdateCustomerSupportRemarkForMoveOK with default headers values
func NewUpdateCustomerSupportRemarkForMoveOK() *UpdateCustomerSupportRemarkForMoveOK {

	return &UpdateCustomerSupportRemarkForMoveOK{}
}

// WithPayload adds the payload to the update customer support remark for move o k response
func (o *UpdateCustomerSupportRemarkForMoveOK) WithPayload(payload *ghcmessages.CustomerSupportRemark) *UpdateCustomerSupportRemarkForMoveOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move o k response
func (o *UpdateCustomerSupportRemarkForMoveOK) SetPayload(payload *ghcmessages.CustomerSupportRemark) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateCustomerSupportRemarkForMoveBadRequestCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveBadRequest
const UpdateCustomerSupportRemarkForMoveBadRequestCode int = 400

/*
UpdateCustomerSupportRemarkForMoveBadRequest The request payload is invalid

swagger:response updateCustomerSupportRemarkForMoveBadRequest
*/
type UpdateCustomerSupportRemarkForMoveBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveBadRequest creates UpdateCustomerSupportRemarkForMoveBadRequest with default headers values
func NewUpdateCustomerSupportRemarkForMoveBadRequest() *UpdateCustomerSupportRemarkForMoveBadRequest {

	return &UpdateCustomerSupportRemarkForMoveBadRequest{}
}

// WithPayload adds the payload to the update customer support remark for move bad request response
func (o *UpdateCustomerSupportRemarkForMoveBadRequest) WithPayload(payload *ghcmessages.Error) *UpdateCustomerSupportRemarkForMoveBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move bad request response
func (o *UpdateCustomerSupportRemarkForMoveBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateCustomerSupportRemarkForMoveForbiddenCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveForbidden
const UpdateCustomerSupportRemarkForMoveForbiddenCode int = 403

/*
UpdateCustomerSupportRemarkForMoveForbidden The request was denied

swagger:response updateCustomerSupportRemarkForMoveForbidden
*/
type UpdateCustomerSupportRemarkForMoveForbidden struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveForbidden creates UpdateCustomerSupportRemarkForMoveForbidden with default headers values
func NewUpdateCustomerSupportRemarkForMoveForbidden() *UpdateCustomerSupportRemarkForMoveForbidden {

	return &UpdateCustomerSupportRemarkForMoveForbidden{}
}

// WithPayload adds the payload to the update customer support remark for move forbidden response
func (o *UpdateCustomerSupportRemarkForMoveForbidden) WithPayload(payload *ghcmessages.Error) *UpdateCustomerSupportRemarkForMoveForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move forbidden response
func (o *UpdateCustomerSupportRemarkForMoveForbidden) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateCustomerSupportRemarkForMoveNotFoundCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveNotFound
const UpdateCustomerSupportRemarkForMoveNotFoundCode int = 404

/*
UpdateCustomerSupportRemarkForMoveNotFound The requested resource wasn't found

swagger:response updateCustomerSupportRemarkForMoveNotFound
*/
type UpdateCustomerSupportRemarkForMoveNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveNotFound creates UpdateCustomerSupportRemarkForMoveNotFound with default headers values
func NewUpdateCustomerSupportRemarkForMoveNotFound() *UpdateCustomerSupportRemarkForMoveNotFound {

	return &UpdateCustomerSupportRemarkForMoveNotFound{}
}

// WithPayload adds the payload to the update customer support remark for move not found response
func (o *UpdateCustomerSupportRemarkForMoveNotFound) WithPayload(payload *ghcmessages.Error) *UpdateCustomerSupportRemarkForMoveNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move not found response
func (o *UpdateCustomerSupportRemarkForMoveNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateCustomerSupportRemarkForMoveUnprocessableEntityCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveUnprocessableEntity
const UpdateCustomerSupportRemarkForMoveUnprocessableEntityCode int = 422

/*
UpdateCustomerSupportRemarkForMoveUnprocessableEntity The payload was unprocessable.

swagger:response updateCustomerSupportRemarkForMoveUnprocessableEntity
*/
type UpdateCustomerSupportRemarkForMoveUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveUnprocessableEntity creates UpdateCustomerSupportRemarkForMoveUnprocessableEntity with default headers values
func NewUpdateCustomerSupportRemarkForMoveUnprocessableEntity() *UpdateCustomerSupportRemarkForMoveUnprocessableEntity {

	return &UpdateCustomerSupportRemarkForMoveUnprocessableEntity{}
}

// WithPayload adds the payload to the update customer support remark for move unprocessable entity response
func (o *UpdateCustomerSupportRemarkForMoveUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *UpdateCustomerSupportRemarkForMoveUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move unprocessable entity response
func (o *UpdateCustomerSupportRemarkForMoveUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateCustomerSupportRemarkForMoveInternalServerErrorCode is the HTTP code returned for type UpdateCustomerSupportRemarkForMoveInternalServerError
const UpdateCustomerSupportRemarkForMoveInternalServerErrorCode int = 500

/*
UpdateCustomerSupportRemarkForMoveInternalServerError A server error occurred

swagger:response updateCustomerSupportRemarkForMoveInternalServerError
*/
type UpdateCustomerSupportRemarkForMoveInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateCustomerSupportRemarkForMoveInternalServerError creates UpdateCustomerSupportRemarkForMoveInternalServerError with default headers values
func NewUpdateCustomerSupportRemarkForMoveInternalServerError() *UpdateCustomerSupportRemarkForMoveInternalServerError {

	return &UpdateCustomerSupportRemarkForMoveInternalServerError{}
}

// WithPayload adds the payload to the update customer support remark for move internal server error response
func (o *UpdateCustomerSupportRemarkForMoveInternalServerError) WithPayload(payload *ghcmessages.Error) *UpdateCustomerSupportRemarkForMoveInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update customer support remark for move internal server error response
func (o *UpdateCustomerSupportRemarkForMoveInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateCustomerSupportRemarkForMoveInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}