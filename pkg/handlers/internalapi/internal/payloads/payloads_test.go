package payloads

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// HandlerSuite is an abstraction of our original suite
type PayloadsSuite struct {
	handlers.BaseHandlerTestSuite
	storer storage.FileStorer
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	hs := &PayloadsSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
		storer: storageTest.NewFakeS3Storage(true),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
