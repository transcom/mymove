// Code generated by go-swagger; DO NOT EDIT.

package mto_service_item

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
)

// UpdateMTOServiceItemStatusReader is a Reader for the UpdateMTOServiceItemStatus structure.
type UpdateMTOServiceItemStatusReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateMTOServiceItemStatusReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateMTOServiceItemStatusOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateMTOServiceItemStatusBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdateMTOServiceItemStatusUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdateMTOServiceItemStatusForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateMTOServiceItemStatusNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewUpdateMTOServiceItemStatusConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 412:
		result := NewUpdateMTOServiceItemStatusPreconditionFailed()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewUpdateMTOServiceItemStatusUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateMTOServiceItemStatusInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PATCH /mto-service-items/{mtoServiceItemID}/status] updateMTOServiceItemStatus", response, response.Code())
	}
}

// NewUpdateMTOServiceItemStatusOK creates a UpdateMTOServiceItemStatusOK with default headers values
func NewUpdateMTOServiceItemStatusOK() *UpdateMTOServiceItemStatusOK {
	return &UpdateMTOServiceItemStatusOK{}
}

/*
UpdateMTOServiceItemStatusOK describes a response with status code 200, with default header values.

Successfully updated service item status for a move task order.
*/
type UpdateMTOServiceItemStatusOK struct {
	Payload supportmessages.MTOServiceItem
}

// IsSuccess returns true when this update m t o service item status o k response has a 2xx status code
func (o *UpdateMTOServiceItemStatusOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update m t o service item status o k response has a 3xx status code
func (o *UpdateMTOServiceItemStatusOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status o k response has a 4xx status code
func (o *UpdateMTOServiceItemStatusOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update m t o service item status o k response has a 5xx status code
func (o *UpdateMTOServiceItemStatusOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status o k response a status code equal to that given
func (o *UpdateMTOServiceItemStatusOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update m t o service item status o k response
func (o *UpdateMTOServiceItemStatusOK) Code() int {
	return 200
}

func (o *UpdateMTOServiceItemStatusOK) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusOK  %+v", 200, o.Payload)
}

func (o *UpdateMTOServiceItemStatusOK) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusOK  %+v", 200, o.Payload)
}

func (o *UpdateMTOServiceItemStatusOK) GetPayload() supportmessages.MTOServiceItem {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload as interface type
	payload, err := supportmessages.UnmarshalMTOServiceItem(response.Body(), consumer)
	if err != nil {
		return err
	}
	o.Payload = payload

	return nil
}

// NewUpdateMTOServiceItemStatusBadRequest creates a UpdateMTOServiceItemStatusBadRequest with default headers values
func NewUpdateMTOServiceItemStatusBadRequest() *UpdateMTOServiceItemStatusBadRequest {
	return &UpdateMTOServiceItemStatusBadRequest{}
}

/*
UpdateMTOServiceItemStatusBadRequest describes a response with status code 400, with default header values.

The request payload is invalid.
*/
type UpdateMTOServiceItemStatusBadRequest struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status bad request response has a 2xx status code
func (o *UpdateMTOServiceItemStatusBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status bad request response has a 3xx status code
func (o *UpdateMTOServiceItemStatusBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status bad request response has a 4xx status code
func (o *UpdateMTOServiceItemStatusBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status bad request response has a 5xx status code
func (o *UpdateMTOServiceItemStatusBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status bad request response a status code equal to that given
func (o *UpdateMTOServiceItemStatusBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update m t o service item status bad request response
func (o *UpdateMTOServiceItemStatusBadRequest) Code() int {
	return 400
}

func (o *UpdateMTOServiceItemStatusBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateMTOServiceItemStatusBadRequest) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateMTOServiceItemStatusBadRequest) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusUnauthorized creates a UpdateMTOServiceItemStatusUnauthorized with default headers values
func NewUpdateMTOServiceItemStatusUnauthorized() *UpdateMTOServiceItemStatusUnauthorized {
	return &UpdateMTOServiceItemStatusUnauthorized{}
}

/*
UpdateMTOServiceItemStatusUnauthorized describes a response with status code 401, with default header values.

The request was denied.
*/
type UpdateMTOServiceItemStatusUnauthorized struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status unauthorized response has a 2xx status code
func (o *UpdateMTOServiceItemStatusUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status unauthorized response has a 3xx status code
func (o *UpdateMTOServiceItemStatusUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status unauthorized response has a 4xx status code
func (o *UpdateMTOServiceItemStatusUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status unauthorized response has a 5xx status code
func (o *UpdateMTOServiceItemStatusUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status unauthorized response a status code equal to that given
func (o *UpdateMTOServiceItemStatusUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the update m t o service item status unauthorized response
func (o *UpdateMTOServiceItemStatusUnauthorized) Code() int {
	return 401
}

func (o *UpdateMTOServiceItemStatusUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateMTOServiceItemStatusUnauthorized) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateMTOServiceItemStatusUnauthorized) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusForbidden creates a UpdateMTOServiceItemStatusForbidden with default headers values
func NewUpdateMTOServiceItemStatusForbidden() *UpdateMTOServiceItemStatusForbidden {
	return &UpdateMTOServiceItemStatusForbidden{}
}

/*
UpdateMTOServiceItemStatusForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type UpdateMTOServiceItemStatusForbidden struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status forbidden response has a 2xx status code
func (o *UpdateMTOServiceItemStatusForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status forbidden response has a 3xx status code
func (o *UpdateMTOServiceItemStatusForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status forbidden response has a 4xx status code
func (o *UpdateMTOServiceItemStatusForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status forbidden response has a 5xx status code
func (o *UpdateMTOServiceItemStatusForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status forbidden response a status code equal to that given
func (o *UpdateMTOServiceItemStatusForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the update m t o service item status forbidden response
func (o *UpdateMTOServiceItemStatusForbidden) Code() int {
	return 403
}

func (o *UpdateMTOServiceItemStatusForbidden) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusForbidden  %+v", 403, o.Payload)
}

func (o *UpdateMTOServiceItemStatusForbidden) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusForbidden  %+v", 403, o.Payload)
}

func (o *UpdateMTOServiceItemStatusForbidden) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusNotFound creates a UpdateMTOServiceItemStatusNotFound with default headers values
func NewUpdateMTOServiceItemStatusNotFound() *UpdateMTOServiceItemStatusNotFound {
	return &UpdateMTOServiceItemStatusNotFound{}
}

/*
UpdateMTOServiceItemStatusNotFound describes a response with status code 404, with default header values.

The requested resource wasn't found.
*/
type UpdateMTOServiceItemStatusNotFound struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status not found response has a 2xx status code
func (o *UpdateMTOServiceItemStatusNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status not found response has a 3xx status code
func (o *UpdateMTOServiceItemStatusNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status not found response has a 4xx status code
func (o *UpdateMTOServiceItemStatusNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status not found response has a 5xx status code
func (o *UpdateMTOServiceItemStatusNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status not found response a status code equal to that given
func (o *UpdateMTOServiceItemStatusNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the update m t o service item status not found response
func (o *UpdateMTOServiceItemStatusNotFound) Code() int {
	return 404
}

func (o *UpdateMTOServiceItemStatusNotFound) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusNotFound  %+v", 404, o.Payload)
}

func (o *UpdateMTOServiceItemStatusNotFound) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusNotFound  %+v", 404, o.Payload)
}

func (o *UpdateMTOServiceItemStatusNotFound) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusConflict creates a UpdateMTOServiceItemStatusConflict with default headers values
func NewUpdateMTOServiceItemStatusConflict() *UpdateMTOServiceItemStatusConflict {
	return &UpdateMTOServiceItemStatusConflict{}
}

/*
UpdateMTOServiceItemStatusConflict describes a response with status code 409, with default header values.

There was a conflict with the request.
*/
type UpdateMTOServiceItemStatusConflict struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status conflict response has a 2xx status code
func (o *UpdateMTOServiceItemStatusConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status conflict response has a 3xx status code
func (o *UpdateMTOServiceItemStatusConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status conflict response has a 4xx status code
func (o *UpdateMTOServiceItemStatusConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status conflict response has a 5xx status code
func (o *UpdateMTOServiceItemStatusConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status conflict response a status code equal to that given
func (o *UpdateMTOServiceItemStatusConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the update m t o service item status conflict response
func (o *UpdateMTOServiceItemStatusConflict) Code() int {
	return 409
}

func (o *UpdateMTOServiceItemStatusConflict) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusConflict  %+v", 409, o.Payload)
}

func (o *UpdateMTOServiceItemStatusConflict) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusConflict  %+v", 409, o.Payload)
}

func (o *UpdateMTOServiceItemStatusConflict) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusPreconditionFailed creates a UpdateMTOServiceItemStatusPreconditionFailed with default headers values
func NewUpdateMTOServiceItemStatusPreconditionFailed() *UpdateMTOServiceItemStatusPreconditionFailed {
	return &UpdateMTOServiceItemStatusPreconditionFailed{}
}

/*
UpdateMTOServiceItemStatusPreconditionFailed describes a response with status code 412, with default header values.

Precondition failed, likely due to a stale eTag (If-Match). Fetch the request again to get the updated eTag value.
*/
type UpdateMTOServiceItemStatusPreconditionFailed struct {
	Payload *supportmessages.ClientError
}

// IsSuccess returns true when this update m t o service item status precondition failed response has a 2xx status code
func (o *UpdateMTOServiceItemStatusPreconditionFailed) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status precondition failed response has a 3xx status code
func (o *UpdateMTOServiceItemStatusPreconditionFailed) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status precondition failed response has a 4xx status code
func (o *UpdateMTOServiceItemStatusPreconditionFailed) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status precondition failed response has a 5xx status code
func (o *UpdateMTOServiceItemStatusPreconditionFailed) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status precondition failed response a status code equal to that given
func (o *UpdateMTOServiceItemStatusPreconditionFailed) IsCode(code int) bool {
	return code == 412
}

// Code gets the status code for the update m t o service item status precondition failed response
func (o *UpdateMTOServiceItemStatusPreconditionFailed) Code() int {
	return 412
}

func (o *UpdateMTOServiceItemStatusPreconditionFailed) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusPreconditionFailed  %+v", 412, o.Payload)
}

func (o *UpdateMTOServiceItemStatusPreconditionFailed) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusPreconditionFailed  %+v", 412, o.Payload)
}

func (o *UpdateMTOServiceItemStatusPreconditionFailed) GetPayload() *supportmessages.ClientError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusPreconditionFailed) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusUnprocessableEntity creates a UpdateMTOServiceItemStatusUnprocessableEntity with default headers values
func NewUpdateMTOServiceItemStatusUnprocessableEntity() *UpdateMTOServiceItemStatusUnprocessableEntity {
	return &UpdateMTOServiceItemStatusUnprocessableEntity{}
}

/*
UpdateMTOServiceItemStatusUnprocessableEntity describes a response with status code 422, with default header values.

The payload was unprocessable.
*/
type UpdateMTOServiceItemStatusUnprocessableEntity struct {
	Payload *supportmessages.ValidationError
}

// IsSuccess returns true when this update m t o service item status unprocessable entity response has a 2xx status code
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status unprocessable entity response has a 3xx status code
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status unprocessable entity response has a 4xx status code
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this update m t o service item status unprocessable entity response has a 5xx status code
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this update m t o service item status unprocessable entity response a status code equal to that given
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the update m t o service item status unprocessable entity response
func (o *UpdateMTOServiceItemStatusUnprocessableEntity) Code() int {
	return 422
}

func (o *UpdateMTOServiceItemStatusUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdateMTOServiceItemStatusUnprocessableEntity) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *UpdateMTOServiceItemStatusUnprocessableEntity) GetPayload() *supportmessages.ValidationError {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateMTOServiceItemStatusInternalServerError creates a UpdateMTOServiceItemStatusInternalServerError with default headers values
func NewUpdateMTOServiceItemStatusInternalServerError() *UpdateMTOServiceItemStatusInternalServerError {
	return &UpdateMTOServiceItemStatusInternalServerError{}
}

/*
UpdateMTOServiceItemStatusInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type UpdateMTOServiceItemStatusInternalServerError struct {
	Payload *supportmessages.Error
}

// IsSuccess returns true when this update m t o service item status internal server error response has a 2xx status code
func (o *UpdateMTOServiceItemStatusInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update m t o service item status internal server error response has a 3xx status code
func (o *UpdateMTOServiceItemStatusInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update m t o service item status internal server error response has a 4xx status code
func (o *UpdateMTOServiceItemStatusInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update m t o service item status internal server error response has a 5xx status code
func (o *UpdateMTOServiceItemStatusInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update m t o service item status internal server error response a status code equal to that given
func (o *UpdateMTOServiceItemStatusInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update m t o service item status internal server error response
func (o *UpdateMTOServiceItemStatusInternalServerError) Code() int {
	return 500
}

func (o *UpdateMTOServiceItemStatusInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateMTOServiceItemStatusInternalServerError) String() string {
	return fmt.Sprintf("[PATCH /mto-service-items/{mtoServiceItemID}/status][%d] updateMTOServiceItemStatusInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateMTOServiceItemStatusInternalServerError) GetPayload() *supportmessages.Error {
	return o.Payload
}

func (o *UpdateMTOServiceItemStatusInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(supportmessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}