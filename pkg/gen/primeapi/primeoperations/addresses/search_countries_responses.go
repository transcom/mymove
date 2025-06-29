// Code generated by go-swagger; DO NOT EDIT.

package addresses

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// SearchCountriesOKCode is the HTTP code returned for type SearchCountriesOK
const SearchCountriesOKCode int = 200

/*
SearchCountriesOK countries matching the search query

swagger:response searchCountriesOK
*/
type SearchCountriesOK struct {

	/*
	  In: Body
	*/
	Payload primemessages.Countries `json:"body,omitempty"`
}

// NewSearchCountriesOK creates SearchCountriesOK with default headers values
func NewSearchCountriesOK() *SearchCountriesOK {

	return &SearchCountriesOK{}
}

// WithPayload adds the payload to the search countries o k response
func (o *SearchCountriesOK) WithPayload(payload primemessages.Countries) *SearchCountriesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search countries o k response
func (o *SearchCountriesOK) SetPayload(payload primemessages.Countries) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchCountriesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = primemessages.Countries{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// SearchCountriesBadRequestCode is the HTTP code returned for type SearchCountriesBadRequest
const SearchCountriesBadRequestCode int = 400

/*
SearchCountriesBadRequest The request payload is invalid.

swagger:response searchCountriesBadRequest
*/
type SearchCountriesBadRequest struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewSearchCountriesBadRequest creates SearchCountriesBadRequest with default headers values
func NewSearchCountriesBadRequest() *SearchCountriesBadRequest {

	return &SearchCountriesBadRequest{}
}

// WithPayload adds the payload to the search countries bad request response
func (o *SearchCountriesBadRequest) WithPayload(payload *primemessages.ClientError) *SearchCountriesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search countries bad request response
func (o *SearchCountriesBadRequest) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchCountriesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SearchCountriesForbiddenCode is the HTTP code returned for type SearchCountriesForbidden
const SearchCountriesForbiddenCode int = 403

/*
SearchCountriesForbidden The request was denied.

swagger:response searchCountriesForbidden
*/
type SearchCountriesForbidden struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewSearchCountriesForbidden creates SearchCountriesForbidden with default headers values
func NewSearchCountriesForbidden() *SearchCountriesForbidden {

	return &SearchCountriesForbidden{}
}

// WithPayload adds the payload to the search countries forbidden response
func (o *SearchCountriesForbidden) WithPayload(payload *primemessages.ClientError) *SearchCountriesForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search countries forbidden response
func (o *SearchCountriesForbidden) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchCountriesForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SearchCountriesNotFoundCode is the HTTP code returned for type SearchCountriesNotFound
const SearchCountriesNotFoundCode int = 404

/*
SearchCountriesNotFound The requested resource wasn't found.

swagger:response searchCountriesNotFound
*/
type SearchCountriesNotFound struct {

	/*
	  In: Body
	*/
	Payload *primemessages.ClientError `json:"body,omitempty"`
}

// NewSearchCountriesNotFound creates SearchCountriesNotFound with default headers values
func NewSearchCountriesNotFound() *SearchCountriesNotFound {

	return &SearchCountriesNotFound{}
}

// WithPayload adds the payload to the search countries not found response
func (o *SearchCountriesNotFound) WithPayload(payload *primemessages.ClientError) *SearchCountriesNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search countries not found response
func (o *SearchCountriesNotFound) SetPayload(payload *primemessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchCountriesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SearchCountriesInternalServerErrorCode is the HTTP code returned for type SearchCountriesInternalServerError
const SearchCountriesInternalServerErrorCode int = 500

/*
SearchCountriesInternalServerError A server error occurred.

swagger:response searchCountriesInternalServerError
*/
type SearchCountriesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *primemessages.Error `json:"body,omitempty"`
}

// NewSearchCountriesInternalServerError creates SearchCountriesInternalServerError with default headers values
func NewSearchCountriesInternalServerError() *SearchCountriesInternalServerError {

	return &SearchCountriesInternalServerError{}
}

// WithPayload adds the payload to the search countries internal server error response
func (o *SearchCountriesInternalServerError) WithPayload(payload *primemessages.Error) *SearchCountriesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search countries internal server error response
func (o *SearchCountriesInternalServerError) SetPayload(payload *primemessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchCountriesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
