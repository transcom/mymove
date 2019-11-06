package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/handlers/mocks"

	"github.com/gofrs/uuid"
)

func (suite *testSuite) TestTrace() {
	handlerContext := mocks.HandlerContext{}
	handlerContext.On("SetTraceID", mock.Anything).Return(nil)
	mw := Trace(suite.logger, &handlerContext)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.trace, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.NoError(err)           // check that you could read full body
	suite.NotEmpty(string(body)) // check that handler returned the trace id
	id, err := uuid.FromString(string(body))
	suite.Nil(err, "failed to parse UUID")
	suite.logger.Debug(fmt.Sprintf("Trace ID: %s", id))
}
