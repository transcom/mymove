package notifications

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type NotificationSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *NotificationSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestNotificationSuite(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	ns := &NotificationSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, ns)
	ns.PopTestSuite.TearDown()
}
