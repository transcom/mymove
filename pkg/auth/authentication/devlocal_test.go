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

func (suite *AuthSuite) TestCreateUserHandlerMilMove() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

func (suite *AuthSuite) TestCreateUserHandlerOffice() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	// Exercise all variables in the office user
	form := url.Values{}
	form.Add("userType", "office")
	form.Add("firstName", "Carol")
	form.Add("lastName", "X")
	form.Add("telephone", "222-222-2222")
	form.Add("email", "office_user@example.com")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("office_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}

	officeUser, err := models.FetchOfficeUserByEmail(suite.DB(), user.LoginGovEmail)
	if err != nil {
		t.Error("Could not find office user for this user.", err)
	}

	suite.Equal(officeUser.FirstName, "Carol")
	suite.Equal(officeUser.LastName, "X")
	suite.Equal(officeUser.Telephone, "222-222-2222")
}

func (suite *AuthSuite) TestCreateUserHandlerTSP() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "tsp")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.TspServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("tsp_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}

	tspUser, err := models.FetchTspUserByEmail(suite.DB(), user.LoginGovEmail)
	if err != nil {
		t.Error("Could not find tsp user for this user.", err)
	}

	suite.Equal(tspUser.FirstName, "Alice")
	suite.Equal(tspUser.LastName, "Bob")
	suite.Equal(tspUser.Telephone, "333-333-3333")
}

func (suite *AuthSuite) TestCreateUserHandlerDPS() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "dps")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

	_, err = models.FetchDPSUserByEmail(suite.DB(), user.LoginGovEmail)
	if err != nil {
		t.Error("Could not find tsp user for this user.", err)
	}
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToMilMove() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToTSP() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "tsp")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("tsp_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://tsp.example.com:1234/")
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromOfficeToMilMove() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromOfficeToTSP() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "tsp")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("tsp_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://tsp.example.com:1234/")
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromTspToMilMove() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.TspServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromTspToOffice() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.TspServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames, FakeRSAKey, false, false)

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
