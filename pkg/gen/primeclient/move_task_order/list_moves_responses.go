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

// ListMovesReader is a Reader for the ListMoves structure.
type ListMovesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListMovesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListMovesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListMovesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListMovesForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListMovesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /moves] listMoves", response, response.Code())
	}
}

// NewListMovesOK creates a ListMovesOK with default headers values
func NewListMovesOK() *ListMovesOK {
	return &ListMovesOK{}
}

/*
ListMovesOK describes a response with status code 200, with default header values.

Successfully retrieved moves. A successful fetch might still return zero moves.
*/
type ListMovesOK struct {
	Payload primemessages.ListMoves
}

// IsSuccess returns true when this list moves o k response has a 2xx status code
func (o *ListMovesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list moves o k response has a 3xx status code
func (o *ListMovesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list moves o k response has a 4xx status code
func (o *ListMovesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list moves o k response has a 5xx status code
func (o *ListMovesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list moves o k response a status code equal to that given
func (o *ListMovesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list moves o k response
func (o *ListMovesOK) Code() int {
	return 200
}

func (o *ListMovesOK) Error() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesOK  %+v", 200, o.Payload)
}

func (o *ListMovesOK) String() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesOK  %+v", 200, o.Payload)
}

func (o *ListMovesOK) GetPayload() primemessages.ListMoves {
	return o.Payload
}

func (o *ListMovesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListMovesUnauthorized creates a ListMovesUnauthorized with default headers values
func NewListMovesUnauthorized() *ListMovesUnauthorized {
	return &ListMovesUnauthorized{}
}

/*
ListMovesUnauthorized describes a response with status code 401, with default header values.

The request was denied.
*/
type ListMovesUnauthorized struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this list moves unauthorized response has a 2xx status code
func (o *ListMovesUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list moves unauthorized response has a 3xx status code
func (o *ListMovesUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list moves unauthorized response has a 4xx status code
func (o *ListMovesUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this list moves unauthorized response has a 5xx status code
func (o *ListMovesUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this list moves unauthorized response a status code equal to that given
func (o *ListMovesUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the list moves unauthorized response
func (o *ListMovesUnauthorized) Code() int {
	return 401
}

func (o *ListMovesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesUnauthorized  %+v", 401, o.Payload)
}

func (o *ListMovesUnauthorized) String() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesUnauthorized  %+v", 401, o.Payload)
}

func (o *ListMovesUnauthorized) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *ListMovesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListMovesForbidden creates a ListMovesForbidden with default headers values
func NewListMovesForbidden() *ListMovesForbidden {
	return &ListMovesForbidden{}
}

/*
ListMovesForbidden describes a response with status code 403, with default header values.

The request was denied.
*/
type ListMovesForbidden struct {
	Payload *primemessages.ClientError
}

// IsSuccess returns true when this list moves forbidden response has a 2xx status code
func (o *ListMovesForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list moves forbidden response has a 3xx status code
func (o *ListMovesForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list moves forbidden response has a 4xx status code
func (o *ListMovesForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this list moves forbidden response has a 5xx status code
func (o *ListMovesForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this list moves forbidden response a status code equal to that given
func (o *ListMovesForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the list moves forbidden response
func (o *ListMovesForbidden) Code() int {
	return 403
}

func (o *ListMovesForbidden) Error() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesForbidden  %+v", 403, o.Payload)
}

func (o *ListMovesForbidden) String() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesForbidden  %+v", 403, o.Payload)
}

func (o *ListMovesForbidden) GetPayload() *primemessages.ClientError {
	return o.Payload
}

func (o *ListMovesForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.ClientError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListMovesInternalServerError creates a ListMovesInternalServerError with default headers values
func NewListMovesInternalServerError() *ListMovesInternalServerError {
	return &ListMovesInternalServerError{}
}

/*
ListMovesInternalServerError describes a response with status code 500, with default header values.

A server error occurred.
*/
type ListMovesInternalServerError struct {
	Payload *primemessages.Error
}

// IsSuccess returns true when this list moves internal server error response has a 2xx status code
func (o *ListMovesInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list moves internal server error response has a 3xx status code
func (o *ListMovesInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list moves internal server error response has a 4xx status code
func (o *ListMovesInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this list moves internal server error response has a 5xx status code
func (o *ListMovesInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this list moves internal server error response a status code equal to that given
func (o *ListMovesInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the list moves internal server error response
func (o *ListMovesInternalServerError) Code() int {
	return 500
}

func (o *ListMovesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListMovesInternalServerError) String() string {
	return fmt.Sprintf("[GET /moves][%d] listMovesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListMovesInternalServerError) GetPayload() *primemessages.Error {
	return o.Payload
}

func (o *ListMovesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(primemessages.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}