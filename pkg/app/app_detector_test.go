package app

import (
	"log"
	"testing"

	"net/http"

	"net/http/httptest"

	"strings"

	"fmt"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type appSuite struct {
	suite.Suite
	logger *zap.Logger
}

func TestAppSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &appSuite{logger: logger}
	suite.Run(t, hs)
}

var myMoveMil = "my.move.mil"
var officeMoveMil = "office.move.mil"

func (suite *appSuite) TestMiddlewareConstructor() {
	adm := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil)
	suite.NotNil(adm)
}

func (suite *appSuite) TestMiddleWareMyApp() {
	rr := httptest.NewRecorder()

	myMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.True(IsMyApp(r), "first should be myApp")
		suite.False(IsOfficeApp(r), "first should not be officeApp")
		suite.Equal(myMoveMil, GetHostname(r))
	})
	myMoveMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil)(myMoveTestHandler)

	req, _ := http.NewRequest("GET", "/some_url", nil)
	req.Host = myMoveMil
	myMoveMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = strings.ToUpper(myMoveMil)
	myMoveMiddleware.ServeHTTP(rr, req)

	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.False(IsMyApp(r), "should not be myApp")
		suite.True(IsOfficeApp(r), "should be officeApp")
		suite.Equal(officeMoveMil, GetHostname(r))
	})
	officeMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil)(officeTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = fmt.Sprintf("%s:8080", officeMoveMil)
	officeMiddleware.ServeHTTP(rr, req)

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := DetectorMiddleware(suite.logger, myMoveMil, officeMoveMil)(noAppTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = "totally.bogus.hostname"
	noAppMiddleware.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
