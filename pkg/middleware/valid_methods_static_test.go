package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"github.com/pkg/errors"
)

func (suite *testSuite) TestValidMethodStaticGet() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
}

func (suite *testSuite) TestValidMethodStaticHead() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("HEAD", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
}

func (suite *testSuite) TestValidMethodStaticPost() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("POST", testURL, nil))
	suite.Equal(http.StatusMethodNotAllowed, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                                                                        // check that you could read full body
	suite.Equal(http.StatusText(http.StatusMethodNotAllowed)+"\n", string(body), errBody) // check body
}

func (suite *testSuite) TestValidMethodStaticPut() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("PUT", testURL, nil))
	suite.Equal(http.StatusMethodNotAllowed, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                                                                        // check that you could read full body
	suite.Equal(http.StatusText(http.StatusMethodNotAllowed)+"\n", string(body), errBody) // check body
}

func (suite *testSuite) TestValidMethodStaticPatch() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("PATCH", testURL, nil))
	suite.Equal(http.StatusMethodNotAllowed, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                                                                        // check that you could read full body
	suite.Equal(http.StatusText(http.StatusMethodNotAllowed)+"\n", string(body), errBody) // check body
}

func (suite *testSuite) TestValidMethodStaticDelete() {
	mw := ValidMethodsStatic(suite.logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.ok, rr, httptest.NewRequest("DELETE", testURL, nil))
	suite.Equal(http.StatusMethodNotAllowed, rr.Code, errStatusCode) // check status code
	body, err := ioutil.ReadAll(rr.Body)
	suite.Nil(err)                                                                        // check that you could read full body
	suite.Equal(http.StatusText(http.StatusMethodNotAllowed)+"\n", string(body), errBody) // check body
}
