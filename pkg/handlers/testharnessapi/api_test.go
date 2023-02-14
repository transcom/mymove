package testharnessapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TestHarnessAPISuite struct {
	handlers.BaseHandlerTestSuite
}

// func (suite *TestHarnessAPISuite) SetupTest() {

// }

func TestTestHarnessAPISuite(t *testing.T) {
	hs := &TestHarnessAPISuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// tests a post without an accept header
func (suite *TestHarnessAPISuite) TestNewDefaultBuilderNoAcceptHeader() {
	req := httptest.NewRequest("POST", "/build/DefaultMove", nil)
	req = mux.SetURLVars(req, map[string]string{"action": "DefaultMove"})
	rr := httptest.NewRecorder()
	handler := NewDefaultBuilder(suite.HandlerConfig())
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-type"))
}

// tests a post without an accept header
func (suite *TestHarnessAPISuite) TestNewDefaultBuilderWithAcceptHeader() {
	req := httptest.NewRequest("POST", "/build/DefaultMove", nil)
	req = mux.SetURLVars(req, map[string]string{"action": "DefaultMove"})
	req.Header.Add("Accept", "text/html")
	rr := httptest.NewRecorder()
	handler := NewDefaultBuilder(suite.HandlerConfig())
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("text/html", rr.Header().Get("Content-type"))
}

// tests a post without an accept header
func (suite *TestHarnessAPISuite) TestNewBuilderList() {
	req := httptest.NewRequest("POST", "/list", nil)
	rr := httptest.NewRecorder()
	handler := NewBuilderList(suite.HandlerConfig())
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("text/html", rr.Header().Get("Content-type"))

	// the body contains at least one form for building
	suite.True(strings.Contains(rr.Body.String(), `form method="post" action="/testharness/build/`))
}
