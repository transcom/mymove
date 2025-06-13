package privileges

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PrivilegesServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestPrivilegesServiceSuite(t *testing.T) {
	ts := &PrivilegesServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
