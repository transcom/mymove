// Code generated by go-swagger; DO NOT EDIT.

package payment_request

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// CreatePaymentRequestReader is a Reader for the CreatePaymentRequest structure.
type CreatePaymentRequestReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreatePaymentRequestReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreatePaymentRequestCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreatePaymentRequestBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewCreatePaymentRequestUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewCreatePaymentRequestForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewCreatePaymentRequestNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreatePaymentRequestConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewCreatePaymentRequestUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreatePaymentRequestInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /payment-requests] createPaymentRequest", response, response.Code())
	}
}

// NewCreatePaymentRequestCreated creates a CreatePaymentRequestCreated with default headers values
func NewCreatePaymentRequestCreated() *CreatePaymentRequestCreated {
	return &CreatePaymentRequestCreated{}
}

/*
CreatePaymentRequestCreated describes a response with status code 201, with default header values.

Successfully created a paymentRequest object.
*/
type CreatePaymentRequestCreated struct {
	Payload *primemessages.PaymentRequest
}

// IsSuccess returns true when this create payment request created response has a 2xx status code
func (o *CreatePaymentRequestCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create payment request created response has a 3xx status code
func (o *CreatePaymentRequestCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request created response has a 4xx status code
func (o *CreatePaymentRequestCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this create payment request created response has a 5xx status code
func (o *CreatePaymentRequestCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request created response a status code equal to that given
func (o *CreatePaymentRequestCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the create payment request created response
func (o *CreatePaymentRequestCreated) Code() int {
	return 201
}

func (o *CreatePaymentRequestCreated) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestCreated  %+v", 201, o.Payload)
}

func (o *CreatePaymentRequestCreated) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestCreated  %+v", 201, o.Payload)
}

func (o *CreatePaymentRequestCreated) GetPayload() *primemessages.PaymentRequest {
	return o.Payload
}

func (o *CreatePaymentRequestCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.PaymentRequest)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestBadRequest creates a CreatePaymentRequestBadRequest with default headers values
func NewCreatePaymentRequestBadRequest() *CreatePaymentRequestBadRequest {
	return &CreatePaymentRequestBadRequest{}
}

/*
CreatePaymentRequestBadRequest describes a response with status code 400, with default header values.

Request payload is invalid.
*/
type CreatePaymentRequestBadRequest struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create payment request bad request response has a 2xx status code
func (o *CreatePaymentRequestBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request bad request response has a 3xx status code
func (o *CreatePaymentRequestBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request bad request response has a 4xx status code
func (o *CreatePaymentRequestBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request bad request response has a 5xx status code
func (o *CreatePaymentRequestBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request bad request response a status code equal to that given
func (o *CreatePaymentRequestBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the create payment request bad request response
func (o *CreatePaymentRequestBadRequest) Code() int {
	return 400
}

func (o *CreatePaymentRequestBadRequest) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestBadRequest  %+v", 400, o.Payload)
}

func (o *CreatePaymentRequestBadRequest) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestBadRequest  %+v", 400, o.Payload)
}

func (o *CreatePaymentRequestBadRequest) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreatePaymentRequestBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestUnauthorized creates a CreatePaymentRequestUnauthorized with default headers values
func NewCreatePaymentRequestUnauthorized() *CreatePaymentRequestUnauthorized {
	return &CreatePaymentRequestUnauthorized{}
}

/*
CreatePaymentRequestUnauthorized describes a response with status code 401, with default header values.

The request was denied.
*/
type CreatePaymentRequestUnauthorized struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create payment request unauthorized response has a 2xx status code
func (o *CreatePaymentRequestUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request unauthorized response has a 3xx status code
func (o *CreatePaymentRequestUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request unauthorized response has a 4xx status code
func (o *CreatePaymentRequestUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request unauthorized response has a 5xx status code
func (o *CreatePaymentRequestUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request unauthorized response a status code equal to that given
func (o *CreatePaymentRequestUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the create payment request unauthorized response
func (o *CreatePaymentRequestUnauthorized) Code() int {
	return 401
}

func (o *CreatePaymentRequestUnauthorized) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestUnauthorized  %+v", 401, o.Payload)
}

func (o *CreatePaymentRequestUnauthorized) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestUnauthorized  %+v", 401, o.Payload)
}

func (o *CreatePaymentRequestUnauthorized) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreatePaymentRequestUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestForbidden creates a CreatePaymentRequestForbidden with default headers values
func NewCreatePaymentRequestForbidden() *CreatePaymentRequestForbidden {
	return &CreatePaymentRequestForbidden{}
}

/*
CreatePaymentRequestForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type CreatePaymentRequestForbidden struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create payment request forbidden response has a 2xx status code
func (o *CreatePaymentRequestForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request forbidden response has a 3xx status code
func (o *CreatePaymentRequestForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request forbidden response has a 4xx status code
func (o *CreatePaymentRequestForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request forbidden response has a 5xx status code
func (o *CreatePaymentRequestForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request forbidden response a status code equal to that given
func (o *CreatePaymentRequestForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the create payment request forbidden response
func (o *CreatePaymentRequestForbidden) Code() int {
	return 403
}

func (o *CreatePaymentRequestForbidden) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestForbidden  %+v", 403, o.Payload)
}

func (o *CreatePaymentRequestForbidden) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestForbidden  %+v", 403, o.Payload)
}

func (o *CreatePaymentRequestForbidden) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreatePaymentRequestForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestNotFound creates a CreatePaymentRequestNotFound with default headers values
func NewCreatePaymentRequestNotFound() *CreatePaymentRequestNotFound {
	return &CreatePaymentRequestNotFound{}
}

/*
CreatePaymentRequestNotFound describes a response with status code 404, with default header values.

The requested resource wasn't found.
*/
type CreatePaymentRequestNotFound struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create payment request not found response has a 2xx status code
func (o *CreatePaymentRequestNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request not found response has a 3xx status code
func (o *CreatePaymentRequestNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request not found response has a 4xx status code
func (o *CreatePaymentRequestNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request not found response has a 5xx status code
func (o *CreatePaymentRequestNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request not found response a status code equal to that given
func (o *CreatePaymentRequestNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the create payment request not found response
func (o *CreatePaymentRequestNotFound) Code() int {
	return 404
}

func (o *CreatePaymentRequestNotFound) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestNotFound  %+v", 404, o.Payload)
}

func (o *CreatePaymentRequestNotFound) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestNotFound  %+v", 404, o.Payload)
}

func (o *CreatePaymentRequestNotFound) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreatePaymentRequestNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestConflict creates a CreatePaymentRequestConflict with default headers values
func NewCreatePaymentRequestConflict() *CreatePaymentRequestConflict {
	return &CreatePaymentRequestConflict{}
}

/*
CreatePaymentRequestConflict describes a response with status code 409, with default header values.

The request could not be processed because of conflict in the current state of the resource.
*/
type CreatePaymentRequestConflict struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create payment request conflict response has a 2xx status code
func (o *CreatePaymentRequestConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request conflict response has a 3xx status code
func (o *CreatePaymentRequestConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request conflict response has a 4xx status code
func (o *CreatePaymentRequestConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request conflict response has a 5xx status code
func (o *CreatePaymentRequestConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request conflict response a status code equal to that given
func (o *CreatePaymentRequestConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the create payment request conflict response
func (o *CreatePaymentRequestConflict) Code() int {
	return 409
}

func (o *CreatePaymentRequestConflict) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestConflict  %+v", 409, o.Payload)
}

func (o *CreatePaymentRequestConflict) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestConflict  %+v", 409, o.Payload)
}

func (o *CreatePaymentRequestConflict) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreatePaymentRequestConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestUnprocessableEntity creates a CreatePaymentRequestUnprocessableEntity with default headers values
func NewCreatePaymentRequestUnprocessableEntity() *CreatePaymentRequestUnprocessableEntity {
	return &CreatePaymentRequestUnprocessableEntity{}
}

/*
CreatePaymentRequestUnprocessableEntity describes a response with status code 422, with default header values.

The request was unprocessable, likely due to bad input from the requester.
*/
type CreatePaymentRequestUnprocessableEntity struct {
	Payload *primemessages.ValidationError
}

// IsSuccess returns true when this create payment request unprocessable entity response has a 2xx status code
func (o *CreatePaymentRequestUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request unprocessable entity response has a 3xx status code
func (o *CreatePaymentRequestUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request unprocessable entity response has a 4xx status code
func (o *CreatePaymentRequestUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this create payment request unprocessable entity response has a 5xx status code
func (o *CreatePaymentRequestUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this create payment request unprocessable entity response a status code equal to that given
func (o *CreatePaymentRequestUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the create payment request unprocessable entity response
func (o *CreatePaymentRequestUnprocessableEntity) Code() int {
	return 422
}

func (o *CreatePaymentRequestUnprocessableEntity) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreatePaymentRequestUnprocessableEntity) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreatePaymentRequestUnprocessableEntity) GetPayload() *primemessages.ValidationError {
	return o.Payload
}

func (o *CreatePaymentRequestUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreatePaymentRequestInternalServerError creates a CreatePaymentRequestInternalServerError with default headers values
func NewCreatePaymentRequestInternalServerError() *CreatePaymentRequestInternalServerError {
	return &CreatePaymentRequestInternalServerError{}
}

/*
CreatePaymentRequestInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type CreatePaymentRequestInternalServerError struct {
	Payload *primemessages.Error
}

// IsSuccess returns true when this create payment request internal server error response has a 2xx status code
func (o *CreatePaymentRequestInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create payment request internal server error response has a 3xx status code
func (o *CreatePaymentRequestInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create payment request internal server error response has a 4xx status code
func (o *CreatePaymentRequestInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create payment request internal server error response has a 5xx status code
func (o *CreatePaymentRequestInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create payment request internal server error response a status code equal to that given
func (o *CreatePaymentRequestInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create payment request internal server error response
func (o *CreatePaymentRequestInternalServerError) Code() int {
	return 500
}

func (o *CreatePaymentRequestInternalServerError) Error() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestInternalServerError  %+v", 500, o.Payload)
}

func (o *CreatePaymentRequestInternalServerError) String() string {
	return fmt.Sprintf("[POST /payment-requests][%d] createPaymentRequestInternalServerError  %+v", 500, o.Payload)
}

func (o *CreatePaymentRequestInternalServerError) GetPayload() *primemessages.Error {
	return o.Payload
}

func (o *CreatePaymentRequestInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}