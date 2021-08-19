package services

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// PreconditionFailedError is the precondition failed error
type PreconditionFailedError struct {
	id uuid.UUID
	error
}

// NewPreconditionFailedError returns an error for a failed precondition
func NewPreconditionFailedError(id uuid.UUID, err error) PreconditionFailedError {
	return PreconditionFailedError{
		id:    id,
		error: err,
	}
}

// Error is the string representation of the precondition failed error
func (e PreconditionFailedError) Error() string {
	return fmt.Sprintf("Precondition failed on update to object with id: '%s'. The If-Match header value did not match the eTag for this record.", e.id.String())
}

//NotFoundError is returned when a given struct is not found
type NotFoundError struct {
	id      uuid.UUID
	message string
	err     error
}

// NewNotFoundError returns an error for when a struct can not be found
func NewNotFoundError(id uuid.UUID, message string) NotFoundError {
	return NotFoundError{
		id:      id,
		message: message,
	}
}

func (e NotFoundError) Error() string {
	if e.id == uuid.Nil {
		return fmt.Sprintf("not found %s", e.message)
	}
	return fmt.Sprintf("id: %s not found %s", e.id.String(), e.message)
}

// Wrap lets the caller add an error to be wrapped in the NotFoundError
func (e *NotFoundError) Wrap(err error) {
	e.err = err
}

// Unwrap returns the wrapped error, could be nil
func (e *NotFoundError) Unwrap() error {
	return e.err
}

// ErrorCode contains error codes for the route package
type ErrorCode string

// BadDataCode catches errors that are due to bad data being sent and is specifically not a server side error
const BadDataCode ErrorCode = "BAD_DATA"

// Error is used for handling errors from the Route package
type Error interface {
	error
	Code() ErrorCode
}

// baseError contains basic route error functionality
type baseError struct {
	code ErrorCode
}

// Code returns the error code enum
func (b *baseError) Code() ErrorCode {
	return b.code
}

// BadDataError is the custom error type (exported for type checking)
type BadDataError struct {
	baseError
	badDataMsg string
}

// NewBadDataError creates a new BadDataError error
func NewBadDataError(badDataMsg string) *BadDataError {
	return &BadDataError{
		baseError{BadDataCode},
		badDataMsg,
	}
}

func (b *BadDataError) Error() string {
	return fmt.Sprintf("Data received from requester is bad: %s: %s", b.baseError.code, b.badDataMsg)
}

//InvalidInputError is returned when an update fails a validation rule
type InvalidInputError struct {
	id               uuid.UUID
	ValidationErrors *validate.Errors
	message          string
	error
}

// NewInvalidInputError returns an error for invalid input
func NewInvalidInputError(id uuid.UUID, err error, validationErrors *validate.Errors, message string) InvalidInputError {
	return InvalidInputError{
		id:               id,
		error:            err,
		ValidationErrors: validationErrors,
		message:          message,
	}
}

func (e InvalidInputError) Error() string {
	if e.message != "" {
		return fmt.Sprintf(e.message)
	} else if e.id == uuid.Nil && e.ValidationErrors != nil {
		return fmt.Sprintf("Invalid input received. %s", e.ValidationErrors)
	} else if e.ValidationErrors != nil {
		return fmt.Sprintf("Invalid input for id: %s. %s", e.id.String(), e.ValidationErrors)
	} else {
		return ("Invalid Input.")
	}
}

// QueryError is returned when a query in the database failed.
// Use InvalidInputError if you have validation errors to report.
// QueryError is used if you passed validation but the query still failed.
type QueryError struct {
	objectType string
	message    string
	err        error
}

func (e QueryError) Error() string {
	if e.message != "" {
		return fmt.Sprintf(e.message)
	}
	return fmt.Sprintf("Could not complete query related to object of type: %s.", e.objectType)
}

// Unwrap returns the enclosed error
func (e *QueryError) Unwrap() error {
	return e.err
}

// NewQueryError returns an error on a query to the database
// It will create a default error message based on the objectType
// You can override the default message with the msgOverride param
func NewQueryError(objectType string, err error, msgOverride string) QueryError {
	return QueryError{
		objectType: objectType,
		err:        err,
		message:    msgOverride,
	}
}

//InvalidCreateInputError is returned when an update fails a validation rule
type InvalidCreateInputError struct {
	ValidationErrors *validate.Errors
	message          string
}

// NewInvalidCreateInputError returns an error for invalid input
func NewInvalidCreateInputError(validationErrors *validate.Errors, message string) InvalidCreateInputError {
	return InvalidCreateInputError{
		ValidationErrors: validationErrors,
		message:          message,
	}
}

func (e InvalidCreateInputError) Error() string {
	if e.message != "" {
		return fmt.Sprintf(e.message)
	}
	return fmt.Sprintf("invalid input for id: %s", e.ValidationErrors)
}

//ConflictError is returned when a given struct is not found
type ConflictError struct {
	id      uuid.UUID
	message string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("id: %s is in a conflicting state %s", e.id.String(), e.message)
}

// NewConflictError returns an error for when a struct can not be found
func NewConflictError(id uuid.UUID, message string) ConflictError {
	return ConflictError{
		id:      id,
		message: message,
	}
}

// NotImplementedError is returned when some functionality is not yet implemented
type NotImplementedError struct {
	message string
}

func (e NotImplementedError) Error() string {
	return fmt.Sprintf("NotImplementedError: %s", e.message)
}

// NewNotImplementedError creates an error for some unimplemented functionality
func NewNotImplementedError(message string) NotImplementedError {
	return NotImplementedError{
		message: message,
	}
}

// ForbiddenError is the forbidden error
type ForbiddenError struct {
	message string
}

// NewForbiddenError returns an error for forbidden
func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{
		message: message,
	}
}

// Error is the string representation of the forbidden error
func (e ForbiddenError) Error() string {
	return fmt.Sprintf("User is forbidden: %s", e.message)
}

// ImplementationError is returned when a type or function has been implemented incorrectly
// (Typically a dev error)
type ImplementationError struct {
	message string
}

func (e ImplementationError) Error() string {
	return fmt.Sprintf("ImplementationError: %s", e.message)
}

// NewImplementationError creates an error for some unimplemented functionality
func NewImplementationError(message string) ImplementationError {
	return ImplementationError{
		message: message,
	}
}

// EventError is an error generated in the events/notifications system.
// We should log but not return this sort of error to the client because
// client's request could be successful but our notification subsystem
// encountered an error
type EventError struct {
	message string
	error
}

// NewEventError returns a new EventError
func NewEventError(message string, err error) EventError {
	return EventError{
		message: message,
		error:   err,
	}
}

// Error is the string representation of the EventError
func (e EventError) Error() string {
	return e.message
}

// Unwrap returns the wrapped error, could be nil
func (e EventError) Unwrap() error {
	return e.error
}
