package supportapi

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// AfterTest completes tests by trying to close open files
func (hs *HandlerSuite) AfterTest() {
	for _, file := range hs.TestFilesToClose() {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
		//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		// nolint:errcheck
		file.Data.Close()
	}
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
