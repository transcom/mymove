// Code generated by go-swagger; DO NOT EDIT.

package registration

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// CustomerRegistrationCreatedCode is the HTTP code returned for type CustomerRegistrationCreated
const CustomerRegistrationCreatedCode int = 201

/*
CustomerRegistrationCreated successfully registered service member

swagger:response customerRegistrationCreated
*/
type CustomerRegistrationCreated struct {
}

// NewCustomerRegistrationCreated creates CustomerRegistrationCreated with default headers values
func NewCustomerRegistrationCreated() *CustomerRegistrationCreated {

	return &CustomerRegistrationCreated{}
}

// WriteResponse to the client
func (o *CustomerRegistrationCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(201)
}

// CustomerRegistrationUnprocessableEntityCode is the HTTP code returned for type CustomerRegistrationUnprocessableEntity
const CustomerRegistrationUnprocessableEntityCode int = 422

/*
CustomerRegistrationUnprocessableEntity The payload was unprocessable.

swagger:response customerRegistrationUnprocessableEntity
*/
type CustomerRegistrationUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ValidationError `json:"body,omitempty"`
}

// NewCustomerRegistrationUnprocessableEntity creates CustomerRegistrationUnprocessableEntity with default headers values
func NewCustomerRegistrationUnprocessableEntity() *CustomerRegistrationUnprocessableEntity {

	return &CustomerRegistrationUnprocessableEntity{}
}

// WithPayload adds the payload to the customer registration unprocessable entity response
func (o *CustomerRegistrationUnprocessableEntity) WithPayload(payload *internalmessages.ValidationError) *CustomerRegistrationUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the customer registration unprocessable entity response
func (o *CustomerRegistrationUnprocessableEntity) SetPayload(payload *internalmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CustomerRegistrationUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CustomerRegistrationInternalServerErrorCode is the HTTP code returned for type CustomerRegistrationInternalServerError
const CustomerRegistrationInternalServerErrorCode int = 500

/*
CustomerRegistrationInternalServerError internal server error

swagger:response customerRegistrationInternalServerError
*/
type CustomerRegistrationInternalServerError struct {
}

// NewCustomerRegistrationInternalServerError creates CustomerRegistrationInternalServerError with default headers values
func NewCustomerRegistrationInternalServerError() *CustomerRegistrationInternalServerError {

	return &CustomerRegistrationInternalServerError{}
}

// WriteResponse to the client
func (o *CustomerRegistrationInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
