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
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
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

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	handlerConfig.SessionManagers().Mil.LoadAndSave(handler).ServeHTTP(rr, req)

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
	// These roles are created during migrations but our test suite truncates all tables
	testdatagen.MakeTOORole(suite.DB())
	testdatagen.MakeTIORole(suite.DB())
	testdatagen.MakeServicesCounselorRole(suite.DB())
	testdatagen.MakeQaeCsrRole(suite.DB())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, handlerConfig)

	for _, newOfficeUser := range []struct {
		userType string
		roleType roles.RoleType
		email    string
	}{{userType: TOOOfficeUserType, roleType: roles.RoleTypeTOO, email: "too_office_user@example.com"}, {userType: TIOOfficeUserType, roleType: roles.RoleTypeTIO, email: "tio_office_user@example.com"}, {userType: ServicesCounselorOfficeUserType, roleType: roles.RoleTypeServicesCounselor, email: "services_counselor_office_user@example.com"}, {userType: QaeCsrOfficeUserType, roleType: roles.RoleTypeQaeCsr, email: "qae_csr_office_user@example.com"}} {
		// Exercise all variables in the office user
		form := url.Values{}
		form.Add("userType", newOfficeUser.userType)
		form.Add("firstName", "Carol")
		form.Add("lastName", "X")
		form.Add("telephone", "222-222-2222")
		form.Add("email", newOfficeUser.email)

		req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.OfficeServername), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		suite.NoError(req.ParseForm())

		rr := httptest.NewRecorder()
		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			Hostname:        appnames.OfficeServername,
		}
		officeSessionManager := handlerConfig.SessionManagers().Office
		req = suite.SetupSessionRequest(req, session, officeSessionManager)
		officeSessionManager.LoadAndSave(handler).ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

		cookies := rr.Result().Cookies()
		_, err := getCookie("office_session_token", cookies)
		suite.FatalNoError(err, "could not find session token in response")

		user := models.User{}
		err = json.Unmarshal(rr.Body.Bytes(), &user)
		suite.FatalNoError(err, "Could not unmarshal json data into User model.")

		officeUser, err := models.FetchOfficeUserByEmail(suite.DB(), user.LoginGovEmail)
		suite.FatalNoError(err, "Could not find office user for this user.")

		err = suite.DB().Load(officeUser, "TransportationOffice", "User.Roles")
		suite.FatalNoError(err, "Could not load transportation office for this user")

		suite.Equal("Carol", officeUser.FirstName)
		suite.Equal("X", officeUser.LastName)
		suite.Equal("222-222-2222", officeUser.Telephone)
		suite.Equal(newOfficeUser.roleType, officeUser.User.Roles[0].RoleType)
		suite.Equal("KKFA", officeUser.TransportationOffice.Gbloc)
	}
}

func (suite *AuthSuite) TestCreateUserHandlerAdmin() {
	t := suite.T()

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Admin.LoadAndSave(handler).ServeHTTP(rr, req)

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
	queryBuilder := query.NewQueryBuilder()
	filters := []services.QueryFilter{
		query.NewQueryFilter("email", "=", user.LoginGovEmail),
	}

	if err := queryBuilder.FetchOne(suite.AppContextForTest(), &adminUser, filters); err != nil {
		t.Error("Couldn't find admin user record")
	}

	suite.Equal("Leo", adminUser.FirstName)
	suite.Equal("Spaceman", adminUser.LastName)
	suite.Equal(models.SystemAdminRole, adminUser.Role)
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromMilMoveToMilMove() {
	t := suite.T()

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	session := auth.Session{
		ApplicationName: auth.MilApp,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)
	rr := httptest.NewRecorder()
	sessionManagers.Mil.LoadAndSave(handler).ServeHTTP(rr, req.WithContext(ctx))

	serviceMemberID := session.ServiceMemberID
	serviceMember, _ := models.FetchServiceMemberForUser(suite.DB(), &session, serviceMemberID)

	suite.NotEqual(uuid.Nil, serviceMemberID)
	suite.NotEqual(uuid.Nil, serviceMember.UserID)

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
	testdatagen.MakeTOORole(suite.DB())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "TOO office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Office.LoadAndSave(handler).ServeHTTP(rr, req)

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

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Admin.LoadAndSave(handler).ServeHTTP(rr, req)

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

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Mil.LoadAndSave(handler).ServeHTTP(rr, req)

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

	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "admin")

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.OfficeServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Admin.LoadAndSave(handler).ServeHTTP(rr, req)

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
	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	req = suite.SetupSessionRequest(req, &auth.Session{}, sessionManagers.Admin)
	sessionManagers.Mil.LoadAndSave(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")

	cookies := rr.Result().Cookies()
	_, err := getCookie("mil_session_token", cookies)
	suite.FatalNoError(err, "could not find session token in response")

	suite.Equal(rr.Result().Header.Get("Location"), "http://mil.example.com:1234/")
}

func (suite *AuthSuite) TestCreateAndLoginUserHandlerFromAdminToOffice() {
	t := suite.T()
	testdatagen.MakeTOORole(suite.DB())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "TOO office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort)
	sessionManagers := handlerConfig.SessionManagers()
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig)

	rr := httptest.NewRecorder()
	sessionManagers.Office.LoadAndSave(handler).ServeHTTP(rr, req)

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
