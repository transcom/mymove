package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"
)

func (suite *testSuite) TestTrace() {
	mw := Trace(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.trace, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)               // check that you could read full body
	suite.NotEmpty(string(body)) // check that handler returned the trace id
	id, err := uuid.FromString(string(body))
	suite.Nil(err, "failed to parse UUID")
	suite.logger.Debug(fmt.Sprintf("Trace ID: %s", id))
}
