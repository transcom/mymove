package dberr

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBErrSuite struct {
	*testingsuite.PopTestSuite
}

func TestDBFmtSuite(t *testing.T) {
	ts := &DBErrSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
