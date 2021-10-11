package notifications

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type NotificationSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// AppContextForTest returns the AppContext for the test suite
func (suite *NotificationSuite) AppContextForTest(session *auth.Session) appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, session)
}

func TestNotificationSuite(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	ns := &NotificationSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       logger,
	}
	suite.Run(t, ns)
	ns.PopTestSuite.TearDown()
}
