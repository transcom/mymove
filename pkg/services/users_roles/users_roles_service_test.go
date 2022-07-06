package usersroles

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UsersRolesServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestUsersRolesServiceSuite(t *testing.T) {
	ts := &UsersRolesServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
