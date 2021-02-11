package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
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

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req)

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
	form.Add("userType", "PPM office")
	form.Add("firstName", "Carol")
	form.Add("lastName", "X")
	form.Add("telephone", "222-222-2222")
	form.Add("email", "office_user@example.com")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	officeSession.LoadAndSave(handler).ServeHTTP(rr, req)

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

func (suite *AuthSuite) TestCreateUserHandlerDPS() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "dps")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req)

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
		t.Error("Could not find dps user for this user.", err)
	}
}

func (suite *AuthSuite) TestCreateUserHandlerAdmin() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	adminSession := sessionManagers[1]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	adminSession.LoadAndSave(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("admin_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}

	var adminUser models.AdminUser
	queryBuilder := query.NewQueryBuilder(suite.DB())
	filters := []services.QueryFilter{
		query.NewQueryFilter("email", "=", user.LoginGovEmail),
	}

	if err := queryBuilder.FetchOne(&adminUser, filters); err != nil {
		t.Error("Couldn't find admin user record")
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

	session := auth.Session{
		ApplicationName: auth.MilApp,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)
	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req.WithContext(ctx))

	serviceMemberID := session.ServiceMemberID
	serviceMember, _ := models.FetchServiceMemberForUser(ctx, suite.DB(), &session, serviceMemberID)

	suite.NotEqual(uuid.Nil, serviceMemberID)
	suite.NotEqual(uuid.Nil, serviceMember.UserID)
	suite.Equal(false, serviceMember.RequiresAccessCode)

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
	form.Add("userType", "PPM office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	officeSession.LoadAndSave(handler).ServeHTTP(rr, req)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToAdmin() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	adminSession := sessionManagers[1]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	adminSession.LoadAndSave(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("admin_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://admin.example.com:1234/")
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

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromOfficeToAdmin() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)
	adminSession := sessionManagers[1]

	rr := httptest.NewRecorder()
	adminSession.LoadAndSave(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("admin_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	suite.Equal(rr.Result().Header.Get("Location"), "http://admin.example.com:1234/")
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromAdminToMilMove() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req)

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

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromAdminToOffice() {
	t := suite.T()

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "PPM office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), appnames)

	rr := httptest.NewRecorder()
	officeSession.LoadAndSave(handler).ServeHTTP(rr, req)

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
