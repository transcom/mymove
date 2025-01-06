package payloads

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// HandlerSuite is an abstraction of our original suite
type PayloadsSuite struct {
	handlers.BaseHandlerTestSuite
}

func (suite *PayloadsSuite) SetupSuite() {
	suite.PreloadData(func() {
		factory.SetupDefaultAllotments(suite.DB())
	})
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
