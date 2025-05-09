// Code generated by go-swagger; DO NOT EDIT.

package move

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
)

// UpdateAssignedOfficeUserOKCode is the HTTP code returned for type UpdateAssignedOfficeUserOK
const UpdateAssignedOfficeUserOKCode int = 200

/*
UpdateAssignedOfficeUserOK Successfully assigned office user to the move

swagger:response updateAssignedOfficeUserOK
*/
type UpdateAssignedOfficeUserOK struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Move `json:"body,omitempty"`
}

// NewUpdateAssignedOfficeUserOK creates UpdateAssignedOfficeUserOK with default headers values
func NewUpdateAssignedOfficeUserOK() *UpdateAssignedOfficeUserOK {

	return &UpdateAssignedOfficeUserOK{}
}

// WithPayload adds the payload to the update assigned office user o k response
func (o *UpdateAssignedOfficeUserOK) WithPayload(payload *ghcmessages.Move) *UpdateAssignedOfficeUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update assigned office user o k response
func (o *UpdateAssignedOfficeUserOK) SetPayload(payload *ghcmessages.Move) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAssignedOfficeUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateAssignedOfficeUserNotFoundCode is the HTTP code returned for type UpdateAssignedOfficeUserNotFound
const UpdateAssignedOfficeUserNotFoundCode int = 404

/*
UpdateAssignedOfficeUserNotFound The requested resource wasn't found

swagger:response updateAssignedOfficeUserNotFound
*/
type UpdateAssignedOfficeUserNotFound struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateAssignedOfficeUserNotFound creates UpdateAssignedOfficeUserNotFound with default headers values
func NewUpdateAssignedOfficeUserNotFound() *UpdateAssignedOfficeUserNotFound {

	return &UpdateAssignedOfficeUserNotFound{}
}

// WithPayload adds the payload to the update assigned office user not found response
func (o *UpdateAssignedOfficeUserNotFound) WithPayload(payload *ghcmessages.Error) *UpdateAssignedOfficeUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update assigned office user not found response
func (o *UpdateAssignedOfficeUserNotFound) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAssignedOfficeUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateAssignedOfficeUserInternalServerErrorCode is the HTTP code returned for type UpdateAssignedOfficeUserInternalServerError
const UpdateAssignedOfficeUserInternalServerErrorCode int = 500

/*
UpdateAssignedOfficeUserInternalServerError A server error occurred

swagger:response updateAssignedOfficeUserInternalServerError
*/
type UpdateAssignedOfficeUserInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *ghcmessages.Error `json:"body,omitempty"`
}

// NewUpdateAssignedOfficeUserInternalServerError creates UpdateAssignedOfficeUserInternalServerError with default headers values
func NewUpdateAssignedOfficeUserInternalServerError() *UpdateAssignedOfficeUserInternalServerError {

	return &UpdateAssignedOfficeUserInternalServerError{}
}

// WithPayload adds the payload to the update assigned office user internal server error response
func (o *UpdateAssignedOfficeUserInternalServerError) WithPayload(payload *ghcmessages.Error) *UpdateAssignedOfficeUserInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update assigned office user internal server error response
func (o *UpdateAssignedOfficeUserInternalServerError) SetPayload(payload *ghcmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAssignedOfficeUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
