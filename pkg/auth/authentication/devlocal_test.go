package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func getCookie(name string, cookies []*http.Cookie) (*http.Cookie, error) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}
	return nil, errors.Errorf("Unable to find cookie: %s", name)
}

func (suite *AuthSuite) TestCreateUserHandler() {
	t := suite.T()

	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", MilTestHost), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), "fake key", false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("mil_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToMilMove() {
	t := suite.T()

	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", MilTestHost), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), "fake key", false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("mil_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://mil.example.com:1234/")
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToOffice() {
	t := suite.T()

	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", MilTestHost), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), "fake key", false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("office_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://office.example.com:1234/")
}
