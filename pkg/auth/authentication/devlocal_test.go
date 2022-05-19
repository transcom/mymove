//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate test data for use in the unit test
//RA: Creation of test data generation for unit test consumption does not present any unexpected states and conditions
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
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
	"github.com/transcom/mymove/pkg/handlers"
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

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateUserHandler(authContext, handlerConfig, appnames)

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

	// These roles are created during migrations but our test suite truncates all tables
	testdatagen.MakePPMOfficeRole(suite.DB())
	testdatagen.MakeTOORole(suite.DB())
	testdatagen.MakeTIORole(suite.DB())
	testdatagen.MakeServicesCounselorRole(suite.DB())
	testdatagen.MakeQaeCsrRole(suite.DB())

	appnames := ApplicationTestServername()
	callbackPort := 1234

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateUserHandler(authContext, handlerConfig, appnames)

	for _, newOfficeUser := range []struct {
		userType string
		roleType roles.RoleType
		email    string
	}{{userType: PPMOfficeUserType, roleType: roles.RoleTypePPMOfficeUsers, email: "ppm_office_user@example.com"}, {userType: TOOOfficeUserType, roleType: roles.RoleTypeTOO, email: "too_office_user@example.com"}, {userType: TIOOfficeUserType, roleType: roles.RoleTypeTIO, email: "tio_office_user@example.com"}, {userType: ServicesCounselorOfficeUserType, roleType: roles.RoleTypeServicesCounselor, email: "services_counselor_office_user@example.com"}, {userType: QaeCsrOfficeUserType, roleType: roles.RoleTypeQaeCsr, email: "qae_csr_office_user@example.com"}} {
		// Exercise all variables in the office user
		form := url.Values{}
		form.Add("userType", newOfficeUser.userType)
		form.Add("firstName", "Carol")
		form.Add("lastName", "X")
		form.Add("telephone", "222-222-2222")
		form.Add("email", newOfficeUser.email)

		req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.OfficeServername), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		req.ParseForm()

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

		err = suite.DB().Load(officeUser, "TransportationOffice", "User.Roles")
		if err != nil {
			t.Error("Could not load transportation office for this user")
		}

		suite.Equal("Carol", officeUser.FirstName)
		suite.Equal("X", officeUser.LastName)
		suite.Equal("222-222-2222", officeUser.Telephone)
		suite.Equal(newOfficeUser.roleType, officeUser.User.Roles[0].RoleType)
		suite.Equal("KKFA", officeUser.TransportationOffice.Gbloc)
	}
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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateUserHandler(authContext, handlerConfig, appnames)

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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateUserHandler(authContext, handlerConfig, appnames)

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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)
	rr := httptest.NewRecorder()
	milSession.LoadAndSave(handler).ServeHTTP(rr, req.WithContext(ctx))

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
	testdatagen.MakePPMOfficeRole(suite.DB())

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "PPM office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)

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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)

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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)

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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)
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
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)

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
	testdatagen.MakePPMOfficeRole(suite.DB())

	appnames := ApplicationTestServername()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "PPM office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.ParseForm()

	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)

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
