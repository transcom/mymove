// Code generated by go-swagger; DO NOT EDIT.

package moves

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/pptasmessages"
)

// PptasReportsOKCode is the HTTP code returned for type PptasReportsOK
const PptasReportsOKCode int = 200

/*
PptasReportsOK Successfully retrieved pptas reports. A successful fetch might still return zero pptas reports.

swagger:response pptasReportsOK
*/
type PptasReportsOK struct {

	/*
	  In: Body
	*/
	Payload pptasmessages.PPTASReports `json:"body,omitempty"`
}

// NewPptasReportsOK creates PptasReportsOK with default headers values
func NewPptasReportsOK() *PptasReportsOK {

	return &PptasReportsOK{}
}

// WithPayload adds the payload to the pptas reports o k response
func (o *PptasReportsOK) WithPayload(payload pptasmessages.PPTASReports) *PptasReportsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the pptas reports o k response
func (o *PptasReportsOK) SetPayload(payload pptasmessages.PPTASReports) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PptasReportsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = pptasmessages.PPTASReports{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// PptasReportsUnauthorizedCode is the HTTP code returned for type PptasReportsUnauthorized
const PptasReportsUnauthorizedCode int = 401

/*
PptasReportsUnauthorized The request was denied.

swagger:response pptasReportsUnauthorized
*/
type PptasReportsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *pptasmessages.ClientError `json:"body,omitempty"`
}

// NewPptasReportsUnauthorized creates PptasReportsUnauthorized with default headers values
func NewPptasReportsUnauthorized() *PptasReportsUnauthorized {

	return &PptasReportsUnauthorized{}
}

// WithPayload adds the payload to the pptas reports unauthorized response
func (o *PptasReportsUnauthorized) WithPayload(payload *pptasmessages.ClientError) *PptasReportsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the pptas reports unauthorized response
func (o *PptasReportsUnauthorized) SetPayload(payload *pptasmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PptasReportsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PptasReportsForbiddenCode is the HTTP code returned for type PptasReportsForbidden
const PptasReportsForbiddenCode int = 403

/*
PptasReportsForbidden The request was denied.

swagger:response pptasReportsForbidden
*/
type PptasReportsForbidden struct {

	/*
	  In: Body
	*/
	Payload *pptasmessages.ClientError `json:"body,omitempty"`
}

// NewPptasReportsForbidden creates PptasReportsForbidden with default headers values
func NewPptasReportsForbidden() *PptasReportsForbidden {

	return &PptasReportsForbidden{}
}

// WithPayload adds the payload to the pptas reports forbidden response
func (o *PptasReportsForbidden) WithPayload(payload *pptasmessages.ClientError) *PptasReportsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the pptas reports forbidden response
func (o *PptasReportsForbidden) SetPayload(payload *pptasmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PptasReportsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PptasReportsInternalServerErrorCode is the HTTP code returned for type PptasReportsInternalServerError
const PptasReportsInternalServerErrorCode int = 500

/*
PptasReportsInternalServerError An unexpected error has occurred in the server.

swagger:response pptasReportsInternalServerError
*/
type PptasReportsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *pptasmessages.ClientError `json:"body,omitempty"`
}

// NewPptasReportsInternalServerError creates PptasReportsInternalServerError with default headers values
func NewPptasReportsInternalServerError() *PptasReportsInternalServerError {

	return &PptasReportsInternalServerError{}
}

// WithPayload adds the payload to the pptas reports internal server error response
func (o *PptasReportsInternalServerError) WithPayload(payload *pptasmessages.ClientError) *PptasReportsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the pptas reports internal server error response
func (o *PptasReportsInternalServerError) SetPayload(payload *pptasmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PptasReportsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
