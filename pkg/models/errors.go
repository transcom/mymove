package models

import (
	"errors"
)

// These are errors that are returned by various model functions

// FetchError is a base type that can typecast for specific APIs,
// It indicicates why an attempted db fetch failed.
type FetchError string

const (
	// FetchErrorNotFound means that the requested record does not exist
	FetchErrorNotFound FetchError = "NOT_FOUND"
	// FetchErrorForbidden means that the record exists but that the user does not have access to it
	FetchErrorForbidden FetchError = "FORBIDDEN"
)

// ErrCreateViolatesUniqueConstraint is returned if you call create and violate a unique constraint.
var ErrCreateViolatesUniqueConstraint = errors.New("CREATE_VIOLATES_UNIQUE")

// RecordNotFoundErrorString is the error string returned when no matching rows exist in the database
const RecordNotFoundErrorString = "sql: no rows in result set"

// UniqueConstraintViolationErrorPrefix This is the error we get back from dbConnection.Create()
const UniqueConstraintViolationErrorPrefix = "pq: duplicate key value violates unique constraint"
