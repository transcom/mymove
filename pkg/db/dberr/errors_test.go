package dberr

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (suite *DBErrSuite) TestIsDBError() {
	errCode := pgerrcode.InternalError
	dbErr := pq.Error{
		Code: pq.ErrorCode(errCode),
	}

	suite.Run("db error and code match", func() {
		suite.True(IsDBError(&dbErr, errCode))
	})

	suite.Run("not a db error", func() {
		err := errors.New("some random error")
		suite.False(IsDBError(err, errCode))
	})

	suite.Run("not the right db error code", func() {
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

	suite.Run("db error, code, and constraint match", func() {
		suite.True(IsDBErrorForConstraint(&dbErr, errCode, constraintName))
	})

	suite.Run("not a db error", func() {
		err := errors.New("some random error")
		suite.False(IsDBErrorForConstraint(err, errCode, constraintName))
	})

	suite.Run("not the right db error code", func() {
		suite.False(IsDBErrorForConstraint(&dbErr, pgerrcode.InternalError, constraintName))
	})

	suite.Run("not the right constraint name", func() {
		suite.False(IsDBErrorForConstraint(&dbErr, errCode, "bogus"))
	})
}
