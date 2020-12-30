package usersroles

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UsersRolesServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *UsersRolesServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestUsersRolesServiceSuite(t *testing.T) {
	ts := &UsersRolesServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
