package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

var ordersMoveMil = "orders.move.mil"

func (suite *authSuite) TestOrdersDetectorConstructor() {
	adm := HostnameDetectorMiddleware(suite.logger, ordersMoveMil)
	suite.NotNil(adm)
}

func (suite *authSuite) TestOrdersDetector() {
	rr := httptest.NewRecorder()

	ordersDetectorTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	ordersDetector := HostnameDetectorMiddleware(suite.logger, ordersMoveMil)(ordersDetectorTestHandler)

	req, _ := http.NewRequest("GET", "/some_url", nil)
	req.Host = ordersMoveMil
	session := Session{}
	ordersDetector.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusOK, rr.Code, "Should get 200 OK")

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = strings.ToUpper(ordersMoveMil)
	session = Session{}
	ordersDetector.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusOK, rr.Code, "Should get 200 OK")

	notOrdersTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	notOrdersDetector := HostnameDetectorMiddleware(suite.logger, ordersMoveMil)(notOrdersTestHandler)

	req, _ = http.NewRequest("GET", "/some_url", nil)
	req.Host = "totally.bogus.hostname"
	session = Session{}
	notOrdersDetector.ServeHTTP(rr, req.WithContext(SetSessionInRequestContext(req, &session)))
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get 400 Bad Request")
}
