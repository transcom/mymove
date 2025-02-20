// Code generated by go-swagger; DO NOT EDIT.

package payment_request

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// CreateUploadCreatedCode is the HTTP code returned for type CreateUploadCreated
const CreateUploadCreatedCode int = 201

/*
CreateUploadCreated Successfully created upload of digital file.

swagger:response createUploadCreated
*/
type CreateUploadCreated struct {

	/*
	  In: Body
	*/
	Payload *primemessages.UploadWithOmissions `json:"body,omitempty"`
}

// NewCreateUploadCreated creates CreateUploadCreated with default headers values
func NewCreateUploadCreated() *CreateUploadCreated {

	return &CreateUploadCreated{}
}

// WithPayload adds the payload to the create upload created response
func (o *CreateUploadCreated) WithPayload(payload *primemessages.UploadWithOmissions) *CreateUploadCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload created response
func (o *CreateUploadCreated) SetPayload(payload *primemessages.UploadWithOmissions) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadBadRequestCode is the HTTP code returned for type CreateUploadBadRequest
const CreateUploadBadRequestCode int = 400

/*
CreateUploadBadRequest The request payload is invalid.

swagger:response createUploadBadRequest
*/
type CreateUploadBadRequest struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewCreateUploadBadRequest creates CreateUploadBadRequest with default headers values
func NewCreateUploadBadRequest() *CreateUploadBadRequest {

	return &CreateUploadBadRequest{}
}

// WithPayload adds the payload to the create upload bad request response
func (o *CreateUploadBadRequest) WithPayload(payload *primemessages.ClientError) *CreateUploadBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload bad request response
func (o *CreateUploadBadRequest) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadUnauthorizedCode is the HTTP code returned for type CreateUploadUnauthorized
const CreateUploadUnauthorizedCode int = 401

/*
CreateUploadUnauthorized The request was denied.

swagger:response createUploadUnauthorized
*/
type CreateUploadUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewCreateUploadUnauthorized creates CreateUploadUnauthorized with default headers values
func NewCreateUploadUnauthorized() *CreateUploadUnauthorized {

	return &CreateUploadUnauthorized{}
}

// WithPayload adds the payload to the create upload unauthorized response
func (o *CreateUploadUnauthorized) WithPayload(payload *primemessages.ClientError) *CreateUploadUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload unauthorized response
func (o *CreateUploadUnauthorized) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadForbiddenCode is the HTTP code returned for type CreateUploadForbidden
const CreateUploadForbiddenCode int = 403

/*
CreateUploadForbidden The request was denied.

swagger:response createUploadForbidden
*/
type CreateUploadForbidden struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewCreateUploadForbidden creates CreateUploadForbidden with default headers values
func NewCreateUploadForbidden() *CreateUploadForbidden {

	return &CreateUploadForbidden{}
}

// WithPayload adds the payload to the create upload forbidden response
func (o *CreateUploadForbidden) WithPayload(payload *primemessages.ClientError) *CreateUploadForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload forbidden response
func (o *CreateUploadForbidden) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadNotFoundCode is the HTTP code returned for type CreateUploadNotFound
const CreateUploadNotFoundCode int = 404

/*
CreateUploadNotFound The requested resource wasn't found.

swagger:response createUploadNotFound
*/
type CreateUploadNotFound struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewCreateUploadNotFound creates CreateUploadNotFound with default headers values
func NewCreateUploadNotFound() *CreateUploadNotFound {

	return &CreateUploadNotFound{}
}

// WithPayload adds the payload to the create upload not found response
func (o *CreateUploadNotFound) WithPayload(payload *primemessages.ClientError) *CreateUploadNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload not found response
func (o *CreateUploadNotFound) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadUnprocessableEntityCode is the HTTP code returned for type CreateUploadUnprocessableEntity
const CreateUploadUnprocessableEntityCode int = 422

/*
CreateUploadUnprocessableEntity The request was unprocessable, likely due to bad input from the requester.

swagger:response createUploadUnprocessableEntity
*/
type CreateUploadUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ValidationError `json:"body,omitempty"`
}

// NewCreateUploadUnprocessableEntity creates CreateUploadUnprocessableEntity with default headers values
func NewCreateUploadUnprocessableEntity() *CreateUploadUnprocessableEntity {

	return &CreateUploadUnprocessableEntity{}
}

// WithPayload adds the payload to the create upload unprocessable entity response
func (o *CreateUploadUnprocessableEntity) WithPayload(payload *primemessages.ValidationError) *CreateUploadUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload unprocessable entity response
func (o *CreateUploadUnprocessableEntity) SetPayload(payload *primemessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUploadInternalServerErrorCode is the HTTP code returned for type CreateUploadInternalServerError
const CreateUploadInternalServerErrorCode int = 500

/*
CreateUploadInternalServerError A server error occurred.

swagger:response createUploadInternalServerError
*/
type CreateUploadInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *primemessages.Error `json:"body,omitempty"`
}

// NewCreateUploadInternalServerError creates CreateUploadInternalServerError with default headers values
func NewCreateUploadInternalServerError() *CreateUploadInternalServerError {

	return &CreateUploadInternalServerError{}
}

// WithPayload adds the payload to the create upload internal server error response
func (o *CreateUploadInternalServerError) WithPayload(payload *primemessages.Error) *CreateUploadInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create upload internal server error response
func (o *CreateUploadInternalServerError) SetPayload(payload *primemessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUploadInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
