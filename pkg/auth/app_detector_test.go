package auth

import (
	"net/http"

	"net/http/httptest"

	"strings"

	"fmt"
)

func (suite *authSuite) TestMiddlewareConstructor() {
	adm := DetectorMiddleware(suite.logger, MyHost, OfficeHost, TspHost)
	suite.NotNil(adm)
}

func (suite *authSuite) TestMiddlewareMyApp() {
	rr := httptest.NewRecorder()

	milMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMyApp(), "first should be milmove app")
		suite.False(session.IsOfficeApp(), "first should not be office app")
		suite.False(session.IsTspApp(), "first should not be tsp app")
		suite.Equal(MyHost, session.Hostname)
	})
	milMoveMiddleware := DetectorMiddleware(suite.logger, MyHost, OfficeHost, TspHost)(milMoveTestHandler)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/some_url", MyHost), nil)
	session := Session{}
	milMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", MyHost), nil)
	session = Session{}
	milMoveMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(MyHost)), nil)
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
		suite.Equal(OfficeHost, session.Hostname)
	})
	officeMiddleware := DetectorMiddleware(suite.logger, MyHost, OfficeHost, TspHost)(officeTestHandler)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/some_url", OfficeHost), nil)
	session := Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", OfficeHost), nil)
	session = Session{}
	officeMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(OfficeHost)), nil)
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
		suite.Equal(TspHost, session.Hostname)
	})
	tspMiddleware := DetectorMiddleware(suite.logger, MyHost, OfficeHost, TspHost)(tspTestHandler)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/some_url", TspHost), nil)
	session := Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", TspHost), nil)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(TspHost)), nil)
	session = Session{}
	tspMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
}

func (suite *authSuite) TestMiddlewareBadApp() {
	rr := httptest.NewRecorder()

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := DetectorMiddleware(suite.logger, MyHost, OfficeHost, TspHost)(noAppTestHandler)

	req, _ := http.NewRequest("GET", "http://totally.bogus.hostname/some_url", nil)
	session := Session{}
	noAppMiddleware.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
