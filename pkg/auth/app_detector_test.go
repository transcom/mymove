package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (suite *authSuite) TestMiddlewareConstructor() {
	adm := DetectorMiddleware(suite.logger, MyTestHost, OfficeTestHost, TspTestHost)
	suite.NotNil(adm)
}

func (suite *authSuite) TestMiddlewareMyApp() {
	rr := httptest.NewRecorder()

	milMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMyApp(), "first should be milmove app")
		suite.False(session.IsOfficeApp(), "first should not be office app")
		suite.False(session.IsTspApp(), "first should not be tsp app")
		suite.Equal(MyTestHost, session.Hostname)
	})
	milMoveMiddleware := DetectorMiddleware(suite.logger, MyTestHost, OfficeTestHost, TspTestHost)(milMoveTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", MyTestHost), nil)
	session := Session{}
	milMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", MyTestHost), nil)
	session = Session{}
	milMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(MyTestHost)), nil)
	session = Session{}
	milMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
}

func (suite *authSuite) TestMiddlwareOfficeApp() {
	rr := httptest.NewRecorder()

	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be milmove app")
		suite.True(session.IsOfficeApp(), "should be office app")
		suite.False(session.IsTspApp(), "should not be tsp app")
		suite.Equal(OfficeTestHost, session.Hostname)
	})
	officeMiddleware := DetectorMiddleware(suite.logger, MyTestHost, OfficeTestHost, TspTestHost)(officeTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", OfficeTestHost), nil)
	session := Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", OfficeTestHost), nil)
	session = Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(OfficeTestHost)), nil)
	session = Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
}

func (suite *authSuite) TestMiddlwareTspApp() {
	rr := httptest.NewRecorder()

	tspTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMyApp(), "should not be milmove app")
		suite.False(session.IsOfficeApp(), "should not be office app")
		suite.True(session.IsTspApp(), "should be tsp app")
		suite.Equal(TspTestHost, session.Hostname)
	})
	tspMiddleware := DetectorMiddleware(suite.logger, MyTestHost, OfficeTestHost, TspTestHost)(tspTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", TspTestHost), nil)
	session := Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", TspTestHost), nil)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(TspTestHost)), nil)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
}

func (suite *authSuite) TestMiddlewareBadApp() {
	rr := httptest.NewRecorder()

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := DetectorMiddleware(suite.logger, MyTestHost, OfficeTestHost, TspTestHost)(noAppTestHandler)

	req := httptest.NewRequest("GET", "http://totally.bogus.hostname/some_url", nil)
	session := Session{}
	noAppMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
