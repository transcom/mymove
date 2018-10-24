package server

import (
	"net/http"

	"net/http/httptest"

	"strings"

	"fmt"
)

var testConfig = &HostsConfig{
	MyName:     "my.move.mil",
	OfficeName: "office.move.mil",
	TspName:    "tsp.move.mil",
}

func (suite *serverSuite) TestMiddlewareConstructor() {
	adm := NewAppDetectorMiddleware(testConfig, suite.logger)
	suite.NotNil(adm)
}

func (suite *serverSuite) TestMiddleWareMyApp() {
	rr := httptest.NewRecorder()

	myMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMyApp(), "first should be myApp")
		suite.False(session.IsOfficeApp(), "first should not be officeApp")
		suite.False(session.IsTspApp(), "first should not be tspApp")
		suite.Equal(testConfig.MyName, session.Hostname)
	})
	myMoveMiddleware := NewAppDetectorMiddleware(testConfig, suite.logger)(myMoveTestHandler)

	req, _ := http.NewRequest("GET", "/some_url", nil)
	req.Host = testConfig.MyName
	session := Session{}
	myMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = strings.ToUpper(testConfig.MyName)
	session = Session{}
	myMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be myApp")
		suite.True(session.IsOfficeApp(), "should be officeApp")
		suite.False(session.IsTspApp(), "should not be tspApp")
		suite.Equal(testConfig.OfficeName, session.Hostname)
	})
	officeMiddleware := NewAppDetectorMiddleware(testConfig, suite.logger)(officeTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = fmt.Sprintf("%s:8080", testConfig.OfficeName)
	session = Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	tspTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be myApp")
		suite.False(session.IsOfficeApp(), "should not be officeApp")
		suite.True(session.IsTspApp(), "should be tspApp")
		suite.Equal(testConfig.TspName, session.Hostname)
	})
	tspMiddleware := NewAppDetectorMiddleware(testConfig, suite.logger)(tspTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = fmt.Sprintf("%s:8080", testConfig.TspName)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := NewAppDetectorMiddleware(testConfig, suite.logger)(noAppTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = "totally.bogus.hostname"
	session = Session{}
	noAppMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
