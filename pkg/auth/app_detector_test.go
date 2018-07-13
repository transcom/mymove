package auth

import (
	"net/http"

	"net/http/httptest"

	"strings"

	"fmt"
)

var myMoveMil = "my.move.mil"
var officeMoveMil = "office.move.mil"
var tspMoveMil = "tsp.move.mil"

func (suite *authSuite) TestMiddlewareConstructor() {
	adm := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil, tspMoveMil)
	suite.NotNil(adm)
}

func (suite *authSuite) TestMiddleWareMyApp() {
	rr := httptest.NewRecorder()

	myMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMyApp(), "first should be myApp")
		suite.False(session.IsOfficeApp(), "first should not be officeApp")
		suite.False(session.IsTspApp(), "first should not be tspApp")
		suite.Equal(myMoveMil, session.Hostname)
	})
	myMoveMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil, tspMoveMil)(myMoveTestHandler)

	req, _ := http.NewRequest("GET", "/some_url", nil)
	req.Host = myMoveMil
	session := Session{}
	myMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = strings.ToUpper(myMoveMil)
	session = Session{}
	myMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be myApp")
		suite.True(session.IsOfficeApp(), "should be officeApp")
		suite.False(session.IsTspApp(), "should not be tspApp")
		suite.Equal(officeMoveMil, session.Hostname)
	})
	officeMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil, tspMoveMil)(officeTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = fmt.Sprintf("%s:8080", officeMoveMil)
	session = Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	tspTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be myApp")
		suite.False(session.IsOfficeApp(), "should not be officeApp")
		suite.True(session.IsTspApp(), "should be tspApp")
		suite.Equal(tspMoveMil, session.Hostname)
	})
	tspMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil, tspMoveMil)(tspTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = fmt.Sprintf("%s:8080", tspMoveMil)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil, tspMoveMil)(noAppTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = "totally.bogus.hostname"
	session = Session{}
	noAppMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
