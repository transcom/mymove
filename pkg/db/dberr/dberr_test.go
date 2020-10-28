package dberr

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBErrSuite struct {
	testingsuite.BaseTestSuite
}

func TestDBFmtSuite(t *testing.T) {
	suite.Run(t, new(DBErrSuite))
}
