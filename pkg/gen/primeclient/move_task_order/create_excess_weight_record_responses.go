// Code generated by go-swagger; DO NOT EDIT.

package move_task_order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// CreateExcessWeightRecordReader is a Reader for the CreateExcessWeightRecord structure.
type CreateExcessWeightRecordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateExcessWeightRecordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateExcessWeightRecordCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewCreateExcessWeightRecordUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewCreateExcessWeightRecordForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewCreateExcessWeightRecordNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewCreateExcessWeightRecordUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateExcessWeightRecordInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record] createExcessWeightRecord", response, response.Code())
	}
}

// NewCreateExcessWeightRecordCreated creates a CreateExcessWeightRecordCreated with default headers values
func NewCreateExcessWeightRecordCreated() *CreateExcessWeightRecordCreated {
	return &CreateExcessWeightRecordCreated{}
}

/*
CreateExcessWeightRecordCreated describes a response with status code 201, with default header values.

Successfully uploaded the excess weight record file.
*/
type CreateExcessWeightRecordCreated struct {
	Payload *primemessages.ExcessWeightRecord
}

// IsSuccess returns true when this create excess weight record created response has a 2xx status code
func (o *CreateExcessWeightRecordCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create excess weight record created response has a 3xx status code
func (o *CreateExcessWeightRecordCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record created response has a 4xx status code
func (o *CreateExcessWeightRecordCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this create excess weight record created response has a 5xx status code
func (o *CreateExcessWeightRecordCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this create excess weight record created response a status code equal to that given
func (o *CreateExcessWeightRecordCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the create excess weight record created response
func (o *CreateExcessWeightRecordCreated) Code() int {
	return 201
}

func (o *CreateExcessWeightRecordCreated) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordCreated  %+v", 201, o.Payload)
}

func (o *CreateExcessWeightRecordCreated) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordCreated  %+v", 201, o.Payload)
}

func (o *CreateExcessWeightRecordCreated) GetPayload() *primemessages.ExcessWeightRecord {
	return o.Payload
}

func (o *CreateExcessWeightRecordCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ExcessWeightRecord)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateExcessWeightRecordUnauthorized creates a CreateExcessWeightRecordUnauthorized with default headers values
func NewCreateExcessWeightRecordUnauthorized() *CreateExcessWeightRecordUnauthorized {
	return &CreateExcessWeightRecordUnauthorized{}
}

/*
CreateExcessWeightRecordUnauthorized describes a response with status code 401, with default header values.

The request was denied.
*/
type CreateExcessWeightRecordUnauthorized struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create excess weight record unauthorized response has a 2xx status code
func (o *CreateExcessWeightRecordUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create excess weight record unauthorized response has a 3xx status code
func (o *CreateExcessWeightRecordUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record unauthorized response has a 4xx status code
func (o *CreateExcessWeightRecordUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this create excess weight record unauthorized response has a 5xx status code
func (o *CreateExcessWeightRecordUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this create excess weight record unauthorized response a status code equal to that given
func (o *CreateExcessWeightRecordUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the create excess weight record unauthorized response
func (o *CreateExcessWeightRecordUnauthorized) Code() int {
	return 401
}

func (o *CreateExcessWeightRecordUnauthorized) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordUnauthorized  %+v", 401, o.Payload)
}

func (o *CreateExcessWeightRecordUnauthorized) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordUnauthorized  %+v", 401, o.Payload)
}

func (o *CreateExcessWeightRecordUnauthorized) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreateExcessWeightRecordUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateExcessWeightRecordForbidden creates a CreateExcessWeightRecordForbidden with default headers values
func NewCreateExcessWeightRecordForbidden() *CreateExcessWeightRecordForbidden {
	return &CreateExcessWeightRecordForbidden{}
}

/*
CreateExcessWeightRecordForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type CreateExcessWeightRecordForbidden struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create excess weight record forbidden response has a 2xx status code
func (o *CreateExcessWeightRecordForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create excess weight record forbidden response has a 3xx status code
func (o *CreateExcessWeightRecordForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record forbidden response has a 4xx status code
func (o *CreateExcessWeightRecordForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this create excess weight record forbidden response has a 5xx status code
func (o *CreateExcessWeightRecordForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this create excess weight record forbidden response a status code equal to that given
func (o *CreateExcessWeightRecordForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the create excess weight record forbidden response
func (o *CreateExcessWeightRecordForbidden) Code() int {
	return 403
}

func (o *CreateExcessWeightRecordForbidden) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordForbidden  %+v", 403, o.Payload)
}

func (o *CreateExcessWeightRecordForbidden) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordForbidden  %+v", 403, o.Payload)
}

func (o *CreateExcessWeightRecordForbidden) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreateExcessWeightRecordForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateExcessWeightRecordNotFound creates a CreateExcessWeightRecordNotFound with default headers values
func NewCreateExcessWeightRecordNotFound() *CreateExcessWeightRecordNotFound {
	return &CreateExcessWeightRecordNotFound{}
}

/*
CreateExcessWeightRecordNotFound describes a response with status code 404, with default header values.

The requested resource wasn't found.
*/
type CreateExcessWeightRecordNotFound struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this create excess weight record not found response has a 2xx status code
func (o *CreateExcessWeightRecordNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create excess weight record not found response has a 3xx status code
func (o *CreateExcessWeightRecordNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record not found response has a 4xx status code
func (o *CreateExcessWeightRecordNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this create excess weight record not found response has a 5xx status code
func (o *CreateExcessWeightRecordNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this create excess weight record not found response a status code equal to that given
func (o *CreateExcessWeightRecordNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the create excess weight record not found response
func (o *CreateExcessWeightRecordNotFound) Code() int {
	return 404
}

func (o *CreateExcessWeightRecordNotFound) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordNotFound  %+v", 404, o.Payload)
}

func (o *CreateExcessWeightRecordNotFound) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordNotFound  %+v", 404, o.Payload)
}

func (o *CreateExcessWeightRecordNotFound) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *CreateExcessWeightRecordNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateExcessWeightRecordUnprocessableEntity creates a CreateExcessWeightRecordUnprocessableEntity with default headers values
func NewCreateExcessWeightRecordUnprocessableEntity() *CreateExcessWeightRecordUnprocessableEntity {
	return &CreateExcessWeightRecordUnprocessableEntity{}
}

/*
CreateExcessWeightRecordUnprocessableEntity describes a response with status code 422, with default header values.

The request was unprocessable, likely due to bad input from the requester.
*/
type CreateExcessWeightRecordUnprocessableEntity struct {
	Payload *primemessages.ValidationError
}

// IsSuccess returns true when this create excess weight record unprocessable entity response has a 2xx status code
func (o *CreateExcessWeightRecordUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create excess weight record unprocessable entity response has a 3xx status code
func (o *CreateExcessWeightRecordUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record unprocessable entity response has a 4xx status code
func (o *CreateExcessWeightRecordUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this create excess weight record unprocessable entity response has a 5xx status code
func (o *CreateExcessWeightRecordUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this create excess weight record unprocessable entity response a status code equal to that given
func (o *CreateExcessWeightRecordUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the create excess weight record unprocessable entity response
func (o *CreateExcessWeightRecordUnprocessableEntity) Code() int {
	return 422
}

func (o *CreateExcessWeightRecordUnprocessableEntity) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreateExcessWeightRecordUnprocessableEntity) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreateExcessWeightRecordUnprocessableEntity) GetPayload() *primemessages.ValidationError {
	return o.Payload
}

func (o *CreateExcessWeightRecordUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateExcessWeightRecordInternalServerError creates a CreateExcessWeightRecordInternalServerError with default headers values
func NewCreateExcessWeightRecordInternalServerError() *CreateExcessWeightRecordInternalServerError {
	return &CreateExcessWeightRecordInternalServerError{}
}

/*
CreateExcessWeightRecordInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type CreateExcessWeightRecordInternalServerError struct {
	Payload *primemessages.Error
}

// IsSuccess returns true when this create excess weight record internal server error response has a 2xx status code
func (o *CreateExcessWeightRecordInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create excess weight record internal server error response has a 3xx status code
func (o *CreateExcessWeightRecordInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create excess weight record internal server error response has a 4xx status code
func (o *CreateExcessWeightRecordInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create excess weight record internal server error response has a 5xx status code
func (o *CreateExcessWeightRecordInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create excess weight record internal server error response a status code equal to that given
func (o *CreateExcessWeightRecordInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create excess weight record internal server error response
func (o *CreateExcessWeightRecordInternalServerError) Code() int {
	return 500
}

func (o *CreateExcessWeightRecordInternalServerError) Error() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordInternalServerError  %+v", 500, o.Payload)
}

func (o *CreateExcessWeightRecordInternalServerError) String() string {
	return fmt.Sprintf("[POST /move-task-orders/{moveTaskOrderID}/excess-weight-record][%d] createExcessWeightRecordInternalServerError  %+v", 500, o.Payload)
}

func (o *CreateExcessWeightRecordInternalServerError) GetPayload() *primemessages.Error {
	return o.Payload
}

func (o *CreateExcessWeightRecordInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}