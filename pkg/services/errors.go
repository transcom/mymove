package services

import (
	"fmt"

	"github.com/gobuffalo/validate"
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
	return fmt.Sprintf("id: '%s' could not be updated due to the record being stale", e.id.String())
}

//NotFoundError is returned when a given struct is not found
type NotFoundError struct {
	id      uuid.UUID
	message string
}

// NewNotFoundError returns an error for when a struct can not be found
func NewNotFoundError(id uuid.UUID, message string) NotFoundError {
	return NotFoundError{
		id:      id,
		message: message,
	}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("id: %s not found %s", e.id.String(), e.message)
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
	}
	return fmt.Sprintf("invalid input for id: %s. %s", e.id.String(), e.ValidationErrors)
}
