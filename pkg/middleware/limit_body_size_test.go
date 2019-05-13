package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (suite *testSuite) TestLimitBodySizeValid() {
	mw := LimitBodySize(int64(10), suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.reflect, rr, httptest.NewRequest("GET", testURL, strings.NewReader("foobar")))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                               // check that you could read full body
	suite.Equal("foobar", string(body), errBody) // check body
}

func (suite *testSuite) TestLimitBodySizeInvalid() {
	mw := LimitBodySize(int64(2), suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.reflect, rr, httptest.NewRequest("GET", testURL, strings.NewReader("foobar")))
	suite.Equal(http.StatusBadRequest, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                                                                  // check that you could read full body
	suite.Equal(http.StatusText(http.StatusBadRequest)+"\n", string(body), errBody) // check body
}
