// Code generated by go-swagger; DO NOT EDIT.

package payment_request

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
)

// UpdatePaymentRequestStatusReader is a Reader for the UpdatePaymentRequestStatus structure.
type UpdatePaymentRequestStatusReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdatePaymentRequestStatusReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdatePaymentRequestStatusOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdatePaymentRequestStatusBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdatePaymentRequestStatusUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdatePaymentRequestStatusForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdatePaymentRequestStatusNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewUpdatePaymentRequestStatusConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 412:
		result := NewUpdatePaymentRequestStatusPreconditionFailed()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewUpdatePaymentRequestStatusUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdatePaymentRequestStatusInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PATCH /payment-requests/{paymentRequestID}/status] updatePaymentRequestStatus", response, response.Code())
	}
}

// NewUpdatePaymentRequestStatusOK creates a UpdatePaymentRequestStatusOK with default headers values
func NewUpdatePaymentRequestStatusOK() *UpdatePaymentRequestStatusOK {
	return &UpdatePaymentRequestStatusOK{}
}

/*
UpdatePaymentRequestStatusOK describes a response with status code 200, with default header values.

Successfully updated payment request status.
*/
type UpdatePaymentRequestStatusOK struct {
	Payload *supportmessages.PaymentRequest
}

// IsSuccess returns true when this update payment request status o k response has a 2xx status code
func (o *UpdatePaymentRequestStatusOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update payment request status o k response has a 3xx status code
func (o *UpdatePaymentRequestStatusOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status o k response has a 4xx status code
func (o *UpdatePaymentRequestStatusOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update payment request status o k response has a 5xx status code
func (o *UpdatePaymentRequestStatusOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status o k response a status code equal to that given
func (o *UpdatePaymentRequestStatusOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update payment request status o k response
func (o *UpdatePaymentRequestStatusOK) Code() int {
	return 200
}

func (o *UpdatePaymentRequestStatusOK) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusOK  %+v", 200, o.Payload)
}

func (o *UpdatePaymentRequestStatusOK) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusOK  %+v", 200, o.Payload)
}

func (o *UpdatePaymentRequestStatusOK) GetPayload() *supportmessages.PaymentRequest {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.PaymentRequest)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusBadRequest creates a UpdatePaymentRequestStatusBadRequest with default headers values
func NewUpdatePaymentRequestStatusBadRequest() *UpdatePaymentRequestStatusBadRequest {
	return &UpdatePaymentRequestStatusBadRequest{}
}

/*
UpdatePaymentRequestStatusBadRequest describes a response with status code 400, with default header values.

The request payload is invalid.
*/
type UpdatePaymentRequestStatusBadRequest struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status bad request response has a 2xx status code
func (o *UpdatePaymentRequestStatusBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status bad request response has a 3xx status code
func (o *UpdatePaymentRequestStatusBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status bad request response has a 4xx status code
func (o *UpdatePaymentRequestStatusBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status bad request response has a 5xx status code
func (o *UpdatePaymentRequestStatusBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status bad request response a status code equal to that given
func (o *UpdatePaymentRequestStatusBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update payment request status bad request response
func (o *UpdatePaymentRequestStatusBadRequest) Code() int {
	return 400
}

func (o *UpdatePaymentRequestStatusBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusBadRequest  %+v", 400, o.Payload)
}

func (o *UpdatePaymentRequestStatusBadRequest) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusBadRequest  %+v", 400, o.Payload)
}

func (o *UpdatePaymentRequestStatusBadRequest) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusUnauthorized creates a UpdatePaymentRequestStatusUnauthorized with default headers values
func NewUpdatePaymentRequestStatusUnauthorized() *UpdatePaymentRequestStatusUnauthorized {
	return &UpdatePaymentRequestStatusUnauthorized{}
}

/*
UpdatePaymentRequestStatusUnauthorized describes a response with status code 401, with default header values.

The request was denied.
*/
type UpdatePaymentRequestStatusUnauthorized struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status unauthorized response has a 2xx status code
func (o *UpdatePaymentRequestStatusUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status unauthorized response has a 3xx status code
func (o *UpdatePaymentRequestStatusUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status unauthorized response has a 4xx status code
func (o *UpdatePaymentRequestStatusUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status unauthorized response has a 5xx status code
func (o *UpdatePaymentRequestStatusUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status unauthorized response a status code equal to that given
func (o *UpdatePaymentRequestStatusUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the update payment request status unauthorized response
func (o *UpdatePaymentRequestStatusUnauthorized) Code() int {
	return 401
}

func (o *UpdatePaymentRequestStatusUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdatePaymentRequestStatusUnauthorized) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdatePaymentRequestStatusUnauthorized) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusForbidden creates a UpdatePaymentRequestStatusForbidden with default headers values
func NewUpdatePaymentRequestStatusForbidden() *UpdatePaymentRequestStatusForbidden {
	return &UpdatePaymentRequestStatusForbidden{}
}

/*
UpdatePaymentRequestStatusForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type UpdatePaymentRequestStatusForbidden struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status forbidden response has a 2xx status code
func (o *UpdatePaymentRequestStatusForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status forbidden response has a 3xx status code
func (o *UpdatePaymentRequestStatusForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status forbidden response has a 4xx status code
func (o *UpdatePaymentRequestStatusForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status forbidden response has a 5xx status code
func (o *UpdatePaymentRequestStatusForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status forbidden response a status code equal to that given
func (o *UpdatePaymentRequestStatusForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the update payment request status forbidden response
func (o *UpdatePaymentRequestStatusForbidden) Code() int {
	return 403
}

func (o *UpdatePaymentRequestStatusForbidden) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusForbidden  %+v", 403, o.Payload)
}

func (o *UpdatePaymentRequestStatusForbidden) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusForbidden  %+v", 403, o.Payload)
}

func (o *UpdatePaymentRequestStatusForbidden) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusNotFound creates a UpdatePaymentRequestStatusNotFound with default headers values
func NewUpdatePaymentRequestStatusNotFound() *UpdatePaymentRequestStatusNotFound {
	return &UpdatePaymentRequestStatusNotFound{}
}

/*
UpdatePaymentRequestStatusNotFound describes a response with status code 404, with default header values.

The requested resource wasn't found.
*/
type UpdatePaymentRequestStatusNotFound struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status not found response has a 2xx status code
func (o *UpdatePaymentRequestStatusNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status not found response has a 3xx status code
func (o *UpdatePaymentRequestStatusNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status not found response has a 4xx status code
func (o *UpdatePaymentRequestStatusNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status not found response has a 5xx status code
func (o *UpdatePaymentRequestStatusNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status not found response a status code equal to that given
func (o *UpdatePaymentRequestStatusNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the update payment request status not found response
func (o *UpdatePaymentRequestStatusNotFound) Code() int {
	return 404
}

func (o *UpdatePaymentRequestStatusNotFound) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusNotFound  %+v", 404, o.Payload)
}

func (o *UpdatePaymentRequestStatusNotFound) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusNotFound  %+v", 404, o.Payload)
}

func (o *UpdatePaymentRequestStatusNotFound) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusConflict creates a UpdatePaymentRequestStatusConflict with default headers values
func NewUpdatePaymentRequestStatusConflict() *UpdatePaymentRequestStatusConflict {
	return &UpdatePaymentRequestStatusConflict{}
}

/*
UpdatePaymentRequestStatusConflict describes a response with status code 409, with default header values.

There was a conflict with the request.
*/
type UpdatePaymentRequestStatusConflict struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status conflict response has a 2xx status code
func (o *UpdatePaymentRequestStatusConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status conflict response has a 3xx status code
func (o *UpdatePaymentRequestStatusConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status conflict response has a 4xx status code
func (o *UpdatePaymentRequestStatusConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status conflict response has a 5xx status code
func (o *UpdatePaymentRequestStatusConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status conflict response a status code equal to that given
func (o *UpdatePaymentRequestStatusConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the update payment request status conflict response
func (o *UpdatePaymentRequestStatusConflict) Code() int {
	return 409
}

func (o *UpdatePaymentRequestStatusConflict) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusConflict  %+v", 409, o.Payload)
}

func (o *UpdatePaymentRequestStatusConflict) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusConflict  %+v", 409, o.Payload)
}

func (o *UpdatePaymentRequestStatusConflict) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusPreconditionFailed creates a UpdatePaymentRequestStatusPreconditionFailed with default headers values
func NewUpdatePaymentRequestStatusPreconditionFailed() *UpdatePaymentRequestStatusPreconditionFailed {
	return &UpdatePaymentRequestStatusPreconditionFailed{}
}

/*
UpdatePaymentRequestStatusPreconditionFailed describes a response with status code 412, with default header values.

Precondition failed, likely due to a stale eTag (If-Match). Fetch the request again to get the updated eTag value.
*/
type UpdatePaymentRequestStatusPreconditionFailed struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update payment request status precondition failed response has a 2xx status code
func (o *UpdatePaymentRequestStatusPreconditionFailed) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status precondition failed response has a 3xx status code
func (o *UpdatePaymentRequestStatusPreconditionFailed) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status precondition failed response has a 4xx status code
func (o *UpdatePaymentRequestStatusPreconditionFailed) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status precondition failed response has a 5xx status code
func (o *UpdatePaymentRequestStatusPreconditionFailed) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status precondition failed response a status code equal to that given
func (o *UpdatePaymentRequestStatusPreconditionFailed) IsCode(code int) bool {
	return code == 412
}

// Code gets the status code for the update payment request status precondition failed response
func (o *UpdatePaymentRequestStatusPreconditionFailed) Code() int {
	return 412
}

func (o *UpdatePaymentRequestStatusPreconditionFailed) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusPreconditionFailed  %+v", 412, o.Payload)
}

func (o *UpdatePaymentRequestStatusPreconditionFailed) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusPreconditionFailed  %+v", 412, o.Payload)
}

func (o *UpdatePaymentRequestStatusPreconditionFailed) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusPreconditionFailed) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusUnprocessableEntity creates a UpdatePaymentRequestStatusUnprocessableEntity with default headers values
func NewUpdatePaymentRequestStatusUnprocessableEntity() *UpdatePaymentRequestStatusUnprocessableEntity {
	return &UpdatePaymentRequestStatusUnprocessableEntity{}
}

/*
UpdatePaymentRequestStatusUnprocessableEntity describes a response with status code 422, with default header values.

The payload was unprocessable.
*/
type UpdatePaymentRequestStatusUnprocessableEntity struct {
	Payload *supportmessages.ValidationError
}

// IsSuccess returns true when this update payment request status unprocessable entity response has a 2xx status code
func (o *UpdatePaymentRequestStatusUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status unprocessable entity response has a 3xx status code
func (o *UpdatePaymentRequestStatusUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status unprocessable entity response has a 4xx status code
func (o *UpdatePaymentRequestStatusUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this update payment request status unprocessable entity response has a 5xx status code
func (o *UpdatePaymentRequestStatusUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this update payment request status unprocessable entity response a status code equal to that given
func (o *UpdatePaymentRequestStatusUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the update payment request status unprocessable entity response
func (o *UpdatePaymentRequestStatusUnprocessableEntity) Code() int {
	return 422
}

func (o *UpdatePaymentRequestStatusUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdatePaymentRequestStatusUnprocessableEntity) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdatePaymentRequestStatusUnprocessableEntity) GetPayload() *supportmessages.ValidationError {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentRequestStatusInternalServerError creates a UpdatePaymentRequestStatusInternalServerError with default headers values
func NewUpdatePaymentRequestStatusInternalServerError() *UpdatePaymentRequestStatusInternalServerError {
	return &UpdatePaymentRequestStatusInternalServerError{}
}

/*
UpdatePaymentRequestStatusInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type UpdatePaymentRequestStatusInternalServerError struct {
	Payload *supportmessages.Error
}

// IsSuccess returns true when this update payment request status internal server error response has a 2xx status code
func (o *UpdatePaymentRequestStatusInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update payment request status internal server error response has a 3xx status code
func (o *UpdatePaymentRequestStatusInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update payment request status internal server error response has a 4xx status code
func (o *UpdatePaymentRequestStatusInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update payment request status internal server error response has a 5xx status code
func (o *UpdatePaymentRequestStatusInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update payment request status internal server error response a status code equal to that given
func (o *UpdatePaymentRequestStatusInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update payment request status internal server error response
func (o *UpdatePaymentRequestStatusInternalServerError) Code() int {
	return 500
}

func (o *UpdatePaymentRequestStatusInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdatePaymentRequestStatusInternalServerError) String() string {
	return fmt.Sprintf("[PATCH /payment-requests/{paymentRequestID}/status][%d] updatePaymentRequestStatusInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdatePaymentRequestStatusInternalServerError) GetPayload() *supportmessages.Error {
	return o.Payload
}

func (o *UpdatePaymentRequestStatusInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}