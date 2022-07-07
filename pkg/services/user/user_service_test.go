package user

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UserServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestUserSuite(t *testing.T) {

	ts := &UserServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
