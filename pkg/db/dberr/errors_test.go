package dberr

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (suite *DBErrSuite) TestIsDBError() {
	errCode := pgerrcode.InternalError
	dbErr := pq.Error{
		Code: pq.ErrorCode(errCode),
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
	dbErr := pq.Error{
		Code:       pq.ErrorCode(errCode),
		Constraint: constraintName,
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
