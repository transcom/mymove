package dberr

import (
	"errors"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (suite *DBErrSuite) TestIsDBError() {
	errCode := pgerrcode.InternalError
	dbErr := pgconn.PgError{
		Code: errCode,
	}

	suite.T().Run("db error and code match", func(t *testing.T) {
		suite.True(IsDBError(&dbErr, errCode))
	})

	suite.T().Run("not a db error", func(t *testing.T) {
		err := errors.New("some random error")
		suite.False(IsDBError(err, errCode))
	})

	suite.T().Run("not the right db error code", func(t *testing.T) {
		suite.False(IsDBError(&dbErr, pgerrcode.UniqueViolation))
	})
}

func (suite *DBErrSuite) TestIsDBErrorForConstraint() {
	errCode := pgerrcode.UniqueViolation
	constraintName := "some_unique_constraint"
	dbErr := pgconn.PgError{
		Code:           errCode,
		ConstraintName: constraintName,
	}

	suite.T().Run("db error, code, and constraint match", func(t *testing.T) {
		suite.True(IsDBErrorForConstraint(&dbErr, errCode, constraintName))
	})

	suite.T().Run("not a db error", func(t *testing.T) {
		err := errors.New("some random error")
		suite.False(IsDBErrorForConstraint(err, errCode, constraintName))
	})

	suite.T().Run("not the right db error code", func(t *testing.T) {
		suite.False(IsDBErrorForConstraint(&dbErr, pgerrcode.InternalError, constraintName))
	})

	suite.T().Run("not the right constraint name", func(t *testing.T) {
		suite.False(IsDBErrorForConstraint(&dbErr, errCode, "bogus"))
	})
}
