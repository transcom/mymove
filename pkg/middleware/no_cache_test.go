package middleware

import (
	"net/http"
	"net/http/httptest"
)

func (suite *testSuite) TestNoCache() {
	mw := NoCache(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	result := rr.Result()
	suite.Contains(result.Header, "Cache-Control", errMissingHeader)
	suite.Len(result.Header["Cache-Control"], 1, errInvalidHeader)
	suite.Equal(result.Header["Cache-Control"][0], "no-cache, no-store, must-revalidate", errInvalidHeader)
}
