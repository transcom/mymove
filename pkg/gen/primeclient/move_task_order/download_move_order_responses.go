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

// DownloadMoveOrderReader is a Reader for the DownloadMoveOrder structure.
type DownloadMoveOrderReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *DownloadMoveOrderReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDownloadMoveOrderOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDownloadMoveOrderBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewDownloadMoveOrderForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDownloadMoveOrderNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewDownloadMoveOrderUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDownloadMoveOrderInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /moves/{locator}/documents] downloadMoveOrder", response, response.Code())
	}
}

// NewDownloadMoveOrderOK creates a DownloadMoveOrderOK with default headers values
func NewDownloadMoveOrderOK(writer io.Writer) *DownloadMoveOrderOK {
	return &DownloadMoveOrderOK{

		Payload: writer,
	}
}

/*
DownloadMoveOrderOK describes a response with status code 200, with default header values.

Move Order PDF
*/
type DownloadMoveOrderOK struct {

	/* File name to download
	 */
	ContentDisposition string

	Payload io.Writer
}

// IsSuccess returns true when this download move order o k response has a 2xx status code
func (o *DownloadMoveOrderOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this download move order o k response has a 3xx status code
func (o *DownloadMoveOrderOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order o k response has a 4xx status code
func (o *DownloadMoveOrderOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this download move order o k response has a 5xx status code
func (o *DownloadMoveOrderOK) IsServerError() bool {
	return false
}

// IsCode returns true when this download move order o k response a status code equal to that given
func (o *DownloadMoveOrderOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the download move order o k response
func (o *DownloadMoveOrderOK) Code() int {
	return 200
}

func (o *DownloadMoveOrderOK) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderOK  %+v", 200, o.Payload)
}

func (o *DownloadMoveOrderOK) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderOK  %+v", 200, o.Payload)
}

func (o *DownloadMoveOrderOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *DownloadMoveOrderOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header Content-Disposition
	hdrContentDisposition := response.GetHeader("Content-Disposition")

	if hdrContentDisposition != "" {
		o.ContentDisposition = hdrContentDisposition
	}

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadMoveOrderBadRequest creates a DownloadMoveOrderBadRequest with default headers values
func NewDownloadMoveOrderBadRequest() *DownloadMoveOrderBadRequest {
	return &DownloadMoveOrderBadRequest{}
}

/*
DownloadMoveOrderBadRequest describes a response with status code 400, with default header values.

The request payload is invalid.
*/
type DownloadMoveOrderBadRequest struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this download move order bad request response has a 2xx status code
func (o *DownloadMoveOrderBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download move order bad request response has a 3xx status code
func (o *DownloadMoveOrderBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order bad request response has a 4xx status code
func (o *DownloadMoveOrderBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this download move order bad request response has a 5xx status code
func (o *DownloadMoveOrderBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this download move order bad request response a status code equal to that given
func (o *DownloadMoveOrderBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the download move order bad request response
func (o *DownloadMoveOrderBadRequest) Code() int {
	return 400
}

func (o *DownloadMoveOrderBadRequest) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderBadRequest  %+v", 400, o.Payload)
}

func (o *DownloadMoveOrderBadRequest) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderBadRequest  %+v", 400, o.Payload)
}

func (o *DownloadMoveOrderBadRequest) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *DownloadMoveOrderBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadMoveOrderForbidden creates a DownloadMoveOrderForbidden with default headers values
func NewDownloadMoveOrderForbidden() *DownloadMoveOrderForbidden {
	return &DownloadMoveOrderForbidden{}
}

/*
DownloadMoveOrderForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type DownloadMoveOrderForbidden struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this download move order forbidden response has a 2xx status code
func (o *DownloadMoveOrderForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download move order forbidden response has a 3xx status code
func (o *DownloadMoveOrderForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order forbidden response has a 4xx status code
func (o *DownloadMoveOrderForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this download move order forbidden response has a 5xx status code
func (o *DownloadMoveOrderForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this download move order forbidden response a status code equal to that given
func (o *DownloadMoveOrderForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the download move order forbidden response
func (o *DownloadMoveOrderForbidden) Code() int {
	return 403
}

func (o *DownloadMoveOrderForbidden) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderForbidden  %+v", 403, o.Payload)
}

func (o *DownloadMoveOrderForbidden) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderForbidden  %+v", 403, o.Payload)
}

func (o *DownloadMoveOrderForbidden) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *DownloadMoveOrderForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadMoveOrderNotFound creates a DownloadMoveOrderNotFound with default headers values
func NewDownloadMoveOrderNotFound() *DownloadMoveOrderNotFound {
	return &DownloadMoveOrderNotFound{}
}

/*
DownloadMoveOrderNotFound describes a response with status code 404, with default header values.

The requested resource wasn't found.
*/
type DownloadMoveOrderNotFound struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this download move order not found response has a 2xx status code
func (o *DownloadMoveOrderNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download move order not found response has a 3xx status code
func (o *DownloadMoveOrderNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order not found response has a 4xx status code
func (o *DownloadMoveOrderNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this download move order not found response has a 5xx status code
func (o *DownloadMoveOrderNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this download move order not found response a status code equal to that given
func (o *DownloadMoveOrderNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the download move order not found response
func (o *DownloadMoveOrderNotFound) Code() int {
	return 404
}

func (o *DownloadMoveOrderNotFound) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderNotFound  %+v", 404, o.Payload)
}

func (o *DownloadMoveOrderNotFound) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderNotFound  %+v", 404, o.Payload)
}

func (o *DownloadMoveOrderNotFound) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *DownloadMoveOrderNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadMoveOrderUnprocessableEntity creates a DownloadMoveOrderUnprocessableEntity with default headers values
func NewDownloadMoveOrderUnprocessableEntity() *DownloadMoveOrderUnprocessableEntity {
	return &DownloadMoveOrderUnprocessableEntity{}
}

/*
DownloadMoveOrderUnprocessableEntity describes a response with status code 422, with default header values.

The request was unprocessable, likely due to bad input from the requester.
*/
type DownloadMoveOrderUnprocessableEntity struct {
	Payload *primemessages.ValidationError
}

// IsSuccess returns true when this download move order unprocessable entity response has a 2xx status code
func (o *DownloadMoveOrderUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download move order unprocessable entity response has a 3xx status code
func (o *DownloadMoveOrderUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order unprocessable entity response has a 4xx status code
func (o *DownloadMoveOrderUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this download move order unprocessable entity response has a 5xx status code
func (o *DownloadMoveOrderUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this download move order unprocessable entity response a status code equal to that given
func (o *DownloadMoveOrderUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the download move order unprocessable entity response
func (o *DownloadMoveOrderUnprocessableEntity) Code() int {
	return 422
}

func (o *DownloadMoveOrderUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *DownloadMoveOrderUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *DownloadMoveOrderUnprocessableEntity) GetPayload() *primemessages.ValidationError {
	return o.Payload
}

func (o *DownloadMoveOrderUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadMoveOrderInternalServerError creates a DownloadMoveOrderInternalServerError with default headers values
func NewDownloadMoveOrderInternalServerError() *DownloadMoveOrderInternalServerError {
	return &DownloadMoveOrderInternalServerError{}
}

/*
DownloadMoveOrderInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type DownloadMoveOrderInternalServerError struct {
	Payload *primemessages.Error
}

// IsSuccess returns true when this download move order internal server error response has a 2xx status code
func (o *DownloadMoveOrderInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this download move order internal server error response has a 3xx status code
func (o *DownloadMoveOrderInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this download move order internal server error response has a 4xx status code
func (o *DownloadMoveOrderInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this download move order internal server error response has a 5xx status code
func (o *DownloadMoveOrderInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this download move order internal server error response a status code equal to that given
func (o *DownloadMoveOrderInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the download move order internal server error response
func (o *DownloadMoveOrderInternalServerError) Code() int {
	return 500
}

func (o *DownloadMoveOrderInternalServerError) Error() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderInternalServerError  %+v", 500, o.Payload)
}

func (o *DownloadMoveOrderInternalServerError) String() string {
	return fmt.Sprintf("[GET /moves/{locator}/documents][%d] downloadMoveOrderInternalServerError  %+v", 500, o.Payload)
}

func (o *DownloadMoveOrderInternalServerError) GetPayload() *primemessages.Error {
	return o.Payload
}

func (o *DownloadMoveOrderInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
