package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/handlers"
)

// TestRecoveryPanic tests the recovery middleware with a handler that returns as normal.
func (suite *testSuite) TestRecoveryNormal() {
	mw := Recovery(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
}

type recoveryError struct {
	Title    string
	Detail   string
	Instance string
}

// TestRecoveryPanic tests the recovery middleware with a handler that panics.
func (suite *testSuite) TestRecoveryPanic() {
	mw := Recovery(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.panic, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusInternalServerError, rr.Code, errStatusCode) // check status code
	body, err := io.ReadAll(rr.Body)
	suite.NoError(err) // check that you could read full body
	var response recoveryError
	err = json.Unmarshal(body, &response)
	suite.Nil(err)
	suite.Equal(handlers.InternalServerErrMessage, string(response.Title), errBody)

}
