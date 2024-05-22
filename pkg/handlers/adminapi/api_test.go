package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// AfterTest completes tests by trying to close open files
func (hs *HandlerSuite) AfterTest() {
	for _, file := range hs.TestFilesToClose() {
		err := file.Data.Close()
		hs.Assert().NoError(err)
	}
}

func (hs *HandlerSuite) setupAuthenticatedRequest(method string, url string) *http.Request {
	requestUser := factory.BuildUser(nil, nil, nil)
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
	return func(_ string, _ string, _ interface{}) services.QueryFilter {
		return filter
	}
}
