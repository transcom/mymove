package officeuser

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type OfficeUserServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *OfficeUserServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger)
}

func TestUserSuite(t *testing.T) {

	hs := &OfficeUserServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
