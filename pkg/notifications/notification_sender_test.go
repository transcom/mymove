package notifications

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type NotificationSuite struct {
	*testingsuite.PopTestSuite
}

func TestNotificationSuite(t *testing.T) {

	ns := &NotificationSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ns)
	ns.PopTestSuite.TearDown()
}
