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
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return value in the file is used for test database teardown
	//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	suite.DB().TruncateAll() // nolint:errcheck
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
