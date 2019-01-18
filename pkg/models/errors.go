package models

import (
	"errors"
)

// These are errors that are returned by various model functions

// ErrCreateViolatesUniqueConstraint is returned if you call create and violate a unique constraint.
var ErrCreateViolatesUniqueConstraint = errors.New("CREATE_VIOLATES_UNIQUE")

// ErrFetchNotFound means that the requested record does not exist
var ErrFetchNotFound = errors.New("FETCH_NOT_FOUND")

// ErrUserUnauthorized means that the user is not authorized to access a record
var ErrUserUnauthorized = errors.New("USER_UNAUTHORIZED")

// ErrFetchForbidden means that the record exists but that the user does not have access to it
var ErrFetchForbidden = errors.New("FETCH_FORBIDDEN")

// ErrDestroyForbidden means that a model cannot be destroyed in its current state
var ErrDestroyForbidden = errors.New("DESTROY_FORBIDDEN")

// ErrLocatorGeneration means that we got errors generating the Locator
var ErrLocatorGeneration = errors.New("LOCATOR_ERRORS")

// ErrInvalidPatchGate means that an attempt to patch a model was not given the correct set of fields
var ErrInvalidPatchGate = errors.New("INVALID_PATCH_GATE")

// ErrInvalidTransition is an error representing an invalid state transition.
var ErrInvalidTransition = errors.New("INVALID_TRANSITION")

// recordNotFoundErrorString is the error string returned when no matching rows exist in the database
// This is ugly, but the best we can do with go's Postgresql adapter
const recordNotFoundErrorString = "sql: no rows in result set"

// uniqueConstraintViolationErrorPrefix This is the error we get back from dbConnection.Create()
const uniqueConstraintViolationErrorPrefix = "pq: duplicate key value violates unique constraint"
