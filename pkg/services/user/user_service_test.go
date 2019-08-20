package user

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UserServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *UserServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUserSuite(t *testing.T) {

	hs := &UserServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}
