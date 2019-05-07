package middleware

import (
	"net/http"
	"net/http/httptest"
)

func (suite *testSuite) TestSecurityHeaders() {
	mw := SecurityHeaders(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	result := rr.Result()

	for k, v := range securityHeaders {
		suite.Contains(result.Header, http.CanonicalHeaderKey(k), errMissingHeader)
		suite.Len(result.Header[http.CanonicalHeaderKey(k)], 1, errInvalidHeader)
		suite.Equal(result.Header[http.CanonicalHeaderKey(k)][0], v, errInvalidHeader)
	}

}
