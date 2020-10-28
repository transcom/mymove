package dberr

import (
	"errors"

	"github.com/lib/pq"
)

// IsDBError returns true if the given error is a DB error with the given code.
func IsDBError(err error, errCode string) bool {
	var pgErr *pq.Error
	errors.As(err, &pgErr)
	if errors.As(err, &pgErr) && string(pgErr.Code) == errCode {
		return true
	}

	return false
}

// IsDBErrorForConstraint returns true if the given error is a DB error with the given code and constraint name.
func IsDBErrorForConstraint(err error, errCode string, constraintName string) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && string(pgErr.Code) == errCode && pgErr.Constraint == constraintName {
		return true
	}

	return false
}
