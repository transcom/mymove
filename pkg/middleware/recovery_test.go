package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// TestRecoveryPanic tests the recovery middleware with a handler that returns as normal.
func (suite *testSuite) TestRecoveryNormal() {
	mw := Recovery(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
}

// TestRecoveryPanic tests the recovery middleware with a handler that panics.
func (suite *testSuite) TestRecoveryPanic() {
	mw := Recovery(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.panic, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusInternalServerError, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                         // check that you could read full body
	suite.Equal("", string(body), errBody) // check body (body was written before panic)
}
