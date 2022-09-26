package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/trace"
)

func (suite *testSuite) TestTrace() {
	mw := Trace(suite.logger)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", testURL, nil)
	traceID, err := uuid.NewV4()
	suite.NoError(err, "Error creating a new trace ID.")
	req = req.WithContext(trace.NewContext(req.Context(), traceID))

	suite.do(mw, suite.trace, rr, req)
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	body, err := io.ReadAll(rr.Body)
	suite.NoError(err)           // check that you could read full body
	suite.NotEmpty(string(body)) // check that handler returned the trace id
	id, err := uuid.FromString(string(body))
	suite.Nil(err, "failed to parse UUID")
	suite.logger.Debug(fmt.Sprintf("Trace ID: %s", id))
}
