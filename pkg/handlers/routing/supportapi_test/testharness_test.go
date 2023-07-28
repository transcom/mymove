package supportapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func (suite *SupportAPISuite) TestTestharnessBuild() {
	req := suite.NewMilRequest("POST", "/testharness/build/DefaultMove", nil)
	rr := httptest.NewRecorder()
	suite.SetupSiteHandler().ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("content-type"))
}

func (suite *SupportAPISuite) TestTestharnessList() {
	req := suite.NewMilRequest("GET", "/testharness/list", nil)
	rr := httptest.NewRecorder()
	suite.SetupSiteHandler().ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	listBody, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Contains(string(listBody), `<form method="post"`)
}
