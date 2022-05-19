package payloads

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type PayloadsSuite struct {
	handlers.BaseHandlerTestSuite
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	hs := &PayloadsSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
