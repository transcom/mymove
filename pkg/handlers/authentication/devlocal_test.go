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
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
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

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "milmove")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/create", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitTOORole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitTIORole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitServicesCounselorRole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitQaeRole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitCustomerServiceRepresentativeRole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitHQRole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitGSRRole,
	})
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitPrimeSimulatorRole,
	})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, handlerConfig)

	for _, newOfficeUser := range []struct {
		userType  string
		roleTypes []roles.RoleType
		email     string
	}{
		{
			userType:  TOOOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeTOO},
			email:     "too_office_user@example.com",
		},
		{
			userType:  TIOOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeTIO},
			email:     "tio_office_user@example.com",
		},
		{
			userType:  ServicesCounselorOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeServicesCounselor},
			email:     "services_counselor_office_user@example.com",
		},
		{
			userType:  QaeOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeQae},
			email:     "qae_office_user@example.com",
		},
		{
			userType:  CustomerServiceRepresentativeOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeCustomerServiceRepresentative},
			email:     "customer_service_representative_office_user@example.com",
		},
		{
			userType:  PrimeSimulatorOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypePrimeSimulator},
			email:     "prime_simulator_user@example.com",
		},
		{
			userType: MultiRoleOfficeUserType,
			roleTypes: []roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeTOO,
				roles.RoleTypeServicesCounselor, roles.RoleTypePrimeSimulator, roles.RoleTypeGSR, roles.RoleTypeHQ, roles.RoleTypeQae, roles.RoleTypeCustomerServiceRepresentative},
			email: "multi_role@example.com",
		},
	} {
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

		officeUser, err := models.FetchOfficeUserByEmail(suite.DB(), user.OktaEmail)
		suite.FatalNoError(err, "Could not find office user for this user.")

		err = suite.DB().Load(officeUser, "TransportationOffice", "User.Roles")
		suite.FatalNoError(err, "Could not load transportation office for this user")

		suite.Equal("Carol", officeUser.FirstName)
		suite.Equal("X", officeUser.LastName)
		suite.Equal("222-222-2222", officeUser.Telephone)
		actualRoleTypes := []roles.RoleType{}
		for i := range officeUser.User.Roles {
			actualRoleTypes = append(actualRoleTypes, officeUser.User.Roles[i].RoleType)
		}
		for i := range newOfficeUser.roleTypes {
			suite.Contains(actualRoleTypes, newOfficeUser.roleTypes[i])
		}
		suite.Equal(len(newOfficeUser.roleTypes), len(actualRoleTypes))
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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
		query.NewQueryFilter("email", "=", user.OktaEmail),
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitTOORole,
	})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "TOO office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.MilServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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
	factory.BuildRole(suite.DB(), nil, []factory.Trait{
		factory.GetTraitTOORole,
	})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	callbackPort := 1234

	form := url.Values{}
	form.Add("userType", "TOO office")

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", appnames.AdminServername), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	suite.NoError(req.ParseForm())

	authContext := NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()), "http", callbackPort)
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
