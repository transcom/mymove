package usersprivileges

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UsersPrivilegesServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestUsersPrivilegesServiceSuite(t *testing.T) {
	ts := &UsersPrivilegesServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
