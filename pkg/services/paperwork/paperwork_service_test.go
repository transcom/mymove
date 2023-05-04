package paperwork

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaperworkServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestPaperworkServiceSuite(t *testing.T) {

	ts := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

// setUpMockGotenbergServer sets up a mock Gotenberg server and sets the corresponding env vars to make it easier to
// test functions that hit gotenberg endpoints. This asks the caller to pass in a handler function so that the caller
// can decide how the mock server should respond to their requests (e.g. setting a 404 status code or returning a
// specific file). The caller will need to close the mock server when they are done with it. The easiest way to ensure
// it happens even if the test fails is to defer the call to mockGotenbergServer.Close() right after calling this
// function.
// E.g.
//
//	mockGotenbergServer := suite.setUpMockGotenbergServer(func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusNotFound)
//	})
//
//	defer mockGotenbergServer.Close()
func (suite *PaperworkServiceSuite) setUpMockGotenbergServer(handlerFunc http.HandlerFunc) *httptest.Server {
	mockGotenbergServer := httptest.NewServer(http.HandlerFunc(handlerFunc))

	// The mock server sets its own url so we'll want to break it down and set the corresponding env vars
	url, urlParseErr := url.ParseRequestURI(mockGotenbergServer.URL)

	suite.FatalNoError(urlParseErr)

	os.Setenv(GotenbergProtocol, url.Scheme)
	os.Setenv(GotenbergHost, url.Hostname())
	os.Setenv(GotenbergPort, url.Port())

	return mockGotenbergServer
}
