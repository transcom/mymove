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

// ErrWriteForbidden means that user is not permitted to write the record
var ErrWriteForbidden = errors.New("WRITE_FORBIDDEN")

// ErrWriteConflict means that the record creation or update cannot be completed due to a conflict with other records
var ErrWriteConflict = errors.New("WRITE_CONFLICT")

// ErrDestroyForbidden means that a model cannot be destroyed in its current state
var ErrDestroyForbidden = errors.New("DESTROY_FORBIDDEN")

// ErrLocatorGeneration means that we got errors generating the Locator
var ErrLocatorGeneration = errors.New("LOCATOR_ERRORS")

// ErrInvalidPatchGate means that an attempt to patch a model was not given the correct set of fields
var ErrInvalidPatchGate = errors.New("INVALID_PATCH_GATE")

// ErrInvalidTransition is an error representing an invalid state transition.
var ErrInvalidTransition = errors.New("INVALID_TRANSITION")

// RecordNotFoundErrorString is the error string returned when no matching rows exist in the database
// This is ugly, but the best we can do with go's Postgresql adapter
const RecordNotFoundErrorString = "sql: no rows in result set"

// This is for when the office user email unique idx constraint is hit
const UniqueConstraintViolationOfficeUserEmailErrorString = "pq: duplicate key value violates unique constraint \"office_users_email_idx\""

// This is for when the office user edipi unique idx constraint is hit
const UniqueConstraintViolationOfficeUserEdipiErrorString = "pq: duplicate key value violates unique constraint \"office_users_edipi_key\""

// This is for when the office user other unique id unique idx constraint is hit
const UniqueConstraintViolationOfficeUserOtherUniqueIDErrorString = "pq: duplicate key value violates unique constraint \"office_users_other_unique_id_key\""

// ErrInvalidMoveID is used if a argument is provided in cases where a move ID is provided, but may be malformed, empty, or nonexistent
var ErrInvalidMoveID = errors.New("INVALID_MOVE_ID")

// ErrInvalidOrderID is used if a argument is provided in cases where a order ID is provided, but may be malformed, empty, or nonexistent
var ErrInvalidOrderID = errors.New("INVALID_ORDER_ID")
