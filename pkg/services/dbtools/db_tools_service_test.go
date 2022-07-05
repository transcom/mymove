package dbtools

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBToolsServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestDBToolsServiceSuite(t *testing.T) {
	ts := &DBToolsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
