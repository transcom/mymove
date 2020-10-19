package dberr

import (
	"errors"

	"github.com/jackc/pgconn"
)

// IsDBError returns true if the given error is a DB error with the given code.
func IsDBError(err error, errCode string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == errCode {
		return true
	}

	return false
}

// IsDBErrorForConstraint returns true if the given error is a DB error with the given code and constraint name.
func IsDBErrorForConstraint(err error, errCode string, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == errCode && pgErr.ConstraintName == constraintName {
		return true
	}

	return false
}
