package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
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

func (hs *HandlerSuite) setupAuthenticatedRequest(method string, url string) *http.Request {
	requestUser := testdatagen.MakeStubbedUser(hs.DB())
	req := httptest.NewRequest(method, url, nil)         // We never need to set a body here for these tests, instead
	return hs.AuthenticateAdminRequest(req, requestUser) // we use the generated Params types to set the request body.
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("adminlocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func newMockQueryFilterBuilder(filter *mocks.QueryFilter) services.NewQueryFilter {
	return func(column string, comparator string, value interface{}) services.QueryFilter {
		return filter
	}
}
