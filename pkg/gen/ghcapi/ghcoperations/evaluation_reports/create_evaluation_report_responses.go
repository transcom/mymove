// Code generated by go-swagger; DO NOT EDIT.

package evaluation_reports

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// CreateEvaluationReportOKCode is the HTTP code returned for type CreateEvaluationReportOK
const CreateEvaluationReportOKCode int = 200

/*
CreateEvaluationReportOK Successfully created evaluation report

swagger:response createEvaluationReportOK
*/
type CreateEvaluationReportOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.EvaluationReport `json:"body,omitempty"`
}

// NewCreateEvaluationReportOK creates CreateEvaluationReportOK with default headers values
func NewCreateEvaluationReportOK() *CreateEvaluationReportOK {

	return &CreateEvaluationReportOK{}
}

// WithPayload adds the payload to the create evaluation report o k response
func (o *CreateEvaluationReportOK) WithPayload(payload *ghcmessages.EvaluationReport) *CreateEvaluationReportOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create evaluation report o k response
func (o *CreateEvaluationReportOK) SetPayload(payload *ghcmessages.EvaluationReport) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEvaluationReportOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEvaluationReportBadRequestCode is the HTTP code returned for type CreateEvaluationReportBadRequest
const CreateEvaluationReportBadRequestCode int = 400

/*
CreateEvaluationReportBadRequest The request payload is invalid

swagger:response createEvaluationReportBadRequest
*/
type CreateEvaluationReportBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewCreateEvaluationReportBadRequest creates CreateEvaluationReportBadRequest with default headers values
func NewCreateEvaluationReportBadRequest() *CreateEvaluationReportBadRequest {

	return &CreateEvaluationReportBadRequest{}
}

// WithPayload adds the payload to the create evaluation report bad request response
func (o *CreateEvaluationReportBadRequest) WithPayload(payload *ghcmessages.Error) *CreateEvaluationReportBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create evaluation report bad request response
func (o *CreateEvaluationReportBadRequest) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEvaluationReportBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEvaluationReportNotFoundCode is the HTTP code returned for type CreateEvaluationReportNotFound
const CreateEvaluationReportNotFoundCode int = 404

/*
CreateEvaluationReportNotFound The requested resource wasn't found

swagger:response createEvaluationReportNotFound
*/
type CreateEvaluationReportNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewCreateEvaluationReportNotFound creates CreateEvaluationReportNotFound with default headers values
func NewCreateEvaluationReportNotFound() *CreateEvaluationReportNotFound {

	return &CreateEvaluationReportNotFound{}
}

// WithPayload adds the payload to the create evaluation report not found response
func (o *CreateEvaluationReportNotFound) WithPayload(payload *ghcmessages.Error) *CreateEvaluationReportNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create evaluation report not found response
func (o *CreateEvaluationReportNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEvaluationReportNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEvaluationReportUnprocessableEntityCode is the HTTP code returned for type CreateEvaluationReportUnprocessableEntity
const CreateEvaluationReportUnprocessableEntityCode int = 422

/*
CreateEvaluationReportUnprocessableEntity The payload was unprocessable.

swagger:response createEvaluationReportUnprocessableEntity
*/
type CreateEvaluationReportUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.ValidationError `json:"body,omitempty"`
}

// NewCreateEvaluationReportUnprocessableEntity creates CreateEvaluationReportUnprocessableEntity with default headers values
func NewCreateEvaluationReportUnprocessableEntity() *CreateEvaluationReportUnprocessableEntity {

	return &CreateEvaluationReportUnprocessableEntity{}
}

// WithPayload adds the payload to the create evaluation report unprocessable entity response
func (o *CreateEvaluationReportUnprocessableEntity) WithPayload(payload *ghcmessages.ValidationError) *CreateEvaluationReportUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create evaluation report unprocessable entity response
func (o *CreateEvaluationReportUnprocessableEntity) SetPayload(payload *ghcmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEvaluationReportUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEvaluationReportInternalServerErrorCode is the HTTP code returned for type CreateEvaluationReportInternalServerError
const CreateEvaluationReportInternalServerErrorCode int = 500

/*
CreateEvaluationReportInternalServerError A server error occurred

swagger:response createEvaluationReportInternalServerError
*/
type CreateEvaluationReportInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewCreateEvaluationReportInternalServerError creates CreateEvaluationReportInternalServerError with default headers values
func NewCreateEvaluationReportInternalServerError() *CreateEvaluationReportInternalServerError {

	return &CreateEvaluationReportInternalServerError{}
}

// WithPayload adds the payload to the create evaluation report internal server error response
func (o *CreateEvaluationReportInternalServerError) WithPayload(payload *ghcmessages.Error) *CreateEvaluationReportInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create evaluation report internal server error response
func (o *CreateEvaluationReportInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEvaluationReportInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}