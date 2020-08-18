package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/transcom/mymove/pkg/cli"

	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication/mocks"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AuthSuite) TestCreateTOO() {
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			Active: true,
		},
	})
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", OfficeTestHost), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           officeUser.Email,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)
	// login.gov state cookie
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   "some mis-matched hash value",
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	req.AddCookie(&cookie)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort, sessionManagers)
	h := CallbackHandler{
		authContext,
		suite.DB(),
	}
	rr := httptest.NewRecorder()
	h.SetFeatureFlag(FeatureFlag{Name: cli.FeatureFlagRoleBasedAuth, Active: true})
	officeSession.LoadAndSave(h).ServeHTTP(rr, req)
	suite.Equal(rr.Code, 307)
}

func (suite *AuthSuite) TestAssociateOfficeUser() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	testdatagen.MakeOfficeUserWithNoUser(suite.DB(), testdatagen.Assertions{OfficeUser: models.OfficeUser{
		Email:  user.LoginGovEmail,
		Active: true,
	}})
	ra := officeUserAssociator{suite.DB(), suite.logger}
	_, err := ra.AssociateOfficeUser(&user)
	suite.NoError(err)
	officeUser, err := ra.FetchOfficeUser(user.LoginGovEmail)
	suite.NoError(err)

	suite.Equal(*officeUser.UserID, user.ID)
}

func (suite *AuthSuite) TestAssociateOfficeUserInactiveIsErr() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	testdatagen.MakeOfficeUserWithNoUser(suite.DB(), testdatagen.Assertions{OfficeUser: models.OfficeUser{
		Email:  user.LoginGovEmail,
		Active: false,
	}})
	ra := officeUserAssociator{suite.DB(), suite.logger}
	_, err := ra.AssociateOfficeUser(&user)

	suite.Error(err)
	suite.IsType(err, ErrUserDeactivated)
}

func (suite *AuthSuite) TestAssociateAdminUser() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	testdatagen.MakeAdminUserWithNoUser(suite.DB(), testdatagen.Assertions{AdminUser: models.AdminUser{
		Email:  user.LoginGovEmail,
		Active: true,
	}})
	ra := adminUserAssociator{suite.DB(), suite.logger}
	_, err := ra.AssociateAdminUser(&user)
	suite.NoError(err)
	adminUser, err := ra.FetchAdminUser(user.LoginGovEmail)
	suite.NoError(err)

	suite.Equal(*adminUser.UserID, user.ID)
}

func (suite *AuthSuite) TestAssociateAdminUserInactiveIsErr() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	testdatagen.MakeAdminUserWithNoUser(suite.DB(), testdatagen.Assertions{AdminUser: models.AdminUser{
		Email:  user.LoginGovEmail,
		Active: false,
	}})
	ra := adminUserAssociator{suite.DB(), suite.logger}
	_, err := ra.AssociateAdminUser(&user)

	suite.Error(err)
	suite.IsType(err, ErrUserDeactivated)
}

func (suite *AuthSuite) TestAuthorizeUnknownUser() {
	uua := NewUnknownUserAuthorizer(suite.DB(), suite.logger)
	session := auth.Session{ApplicationName: auth.MilApp}
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	loginGovUser := goth.User{
		UserID: fakeUUID.String(),
		Email:  "sample@email.com",
	}
	err := uua.AuthorizeUnknownUser(loginGovUser, &session)
	suite.NoError(err)

	user := &models.User{}
	err = suite.DB().Where("login_gov_uuid = $1", loginGovUser.UserID).First(user)
	suite.NoError(err)
	suite.Equal(user.ID, session.UserID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserCreateUserFails() {
	oa := &mocks.OfficeUserAssociator{}
	aa := &mocks.AdminUserAssociator{}
	ra := roleAssociator{
		db:                   suite.DB(),
		logger:               suite.logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
	}
	uc := &mocks.UserCreator{}
	uc.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{}, errors.New("error"))
	uua := UnknownUserAuthorizer{
		logger:         suite.logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
	session := auth.Session{ApplicationName: auth.OfficeApp}
	user := goth.User{
		UserID: "id",
		Email:  "sample@email.com",
	}
	err := uua.AuthorizeUnknownUser(user, &session)
	suite.Error(err)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsSystemAdmin() {
	oa := &mocks.OfficeUserAssociator{}
	aa := &mocks.AdminUserAssociator{}
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	aa.On("AssociateAdminUser", mock.Anything).Return(fakeUUID, nil)
	ra := roleAssociator{
		db:                   suite.DB(),
		logger:               suite.logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
	}
	uc := &mocks.UserCreator{}
	uid, _ := uuid.NewV4()
	uc.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{ID: uid}, nil)
	uua := UnknownUserAuthorizer{
		logger:         suite.logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
	user := goth.User{
		UserID: "id",
		Email:  "sample@email.com",
	}
	session := auth.Session{ApplicationName: auth.AdminApp}

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
	aa.AssertNumberOfCalls(suite.T(), "AssociateAdminUser", 1)
	suite.NotEqual(uuid.Nil, session.UserID)
	suite.Equal(fakeUUID, session.AdminUserID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsOfficeUser() {
	oa := &mocks.OfficeUserAssociator{}
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	oa.On("AssociateOfficeUser", mock.Anything).Return(fakeUUID, nil)
	aa := &mocks.AdminUserAssociator{}
	ra := roleAssociator{
		db:                   suite.DB(),
		logger:               suite.logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
	}
	uc := &mocks.UserCreator{}
	uid, _ := uuid.NewV4()
	uc.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{ID: uid}, nil)
	uua := UnknownUserAuthorizer{
		logger:         suite.logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
	user := goth.User{
		UserID: "id",
		Email:  "sample@email.com",
	}
	session := auth.Session{ApplicationName: auth.OfficeApp}
	suite.True(session.IsOfficeApp())

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
	oa.AssertNumberOfCalls(suite.T(), "AssociateOfficeUser", 1)
	suite.Equal(fakeUUID, session.OfficeUserID)
	suite.NotEqual(uuid.Nil, session.UserID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsCustomer() {
	oa := &mocks.OfficeUserAssociator{}
	aa := &mocks.AdminUserAssociator{}
	ra := roleAssociator{
		db:                   suite.DB(),
		logger:               suite.logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
	}
	uc := &mocks.UserCreator{}
	uc.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{}, nil)
	uua := UnknownUserAuthorizer{
		logger:         suite.logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
	user := goth.User{
		UserID: "id",
		Email:  "sample@email.com",
	}
	session := auth.Session{ApplicationName: auth.MilApp}

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsTOOUser() {
	oa := &mocks.OfficeUserAssociator{}
	oa.On("AssociateOfficeUser", mock.Anything).Return(nil, ErrUnauthorized)
	tr := &mocks.TOORoleChecker{}
	tr.On("VerifyHasTOORole", mock.Anything).Return(roles.Role{RoleType: roles.RoleTypeTOO}, nil)
	tr.On("FetchUserIdentity", mock.Anything).Return(&models.UserIdentity{}, nil)
	aa := &mocks.AdminUserAssociator{}
	ra := roleAssociator{
		db:                   suite.DB(),
		logger:               suite.logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
		TOORoleChecker:       tr,
	}
	uc := &mocks.UserCreator{}
	uid, _ := uuid.NewV4()
	uc.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{ID: uid}, nil)
	uua := UnknownUserAuthorizer{
		logger:         suite.logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
	user := goth.User{
		UserID: "id",
		Email:  "sample@email.com",
	}
	session := auth.Session{ApplicationName: auth.OfficeApp}
	suite.True(session.IsOfficeApp())

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
	oa.AssertNumberOfCalls(suite.T(), "AssociateOfficeUser", 1)
	tr.AssertNumberOfCalls(suite.T(), "VerifyHasTOORole", 1)
	suite.True(session.Roles.HasRole(roles.RoleTypeTOO))
	suite.NotEqual(uuid.Nil, session.UserID)
}

func (suite *AuthSuite) TestVerifyHasTOORole() {
	rs := []roles.Role{{
		ID:       uuid.FromStringOrNil("ed2d2cd7-d427-412a-98bb-a9b391d98d32"),
		RoleType: roles.RoleTypeCustomer,
		RoleName: "Customer",
	},
		{
			ID:       uuid.FromStringOrNil("9dc423b6-33b8-493a-a59b-6a823660cb07"),
			RoleType: roles.RoleTypeTOO,
			RoleName: "Transportation Ordering Officer",
		},
	}
	suite.NoError(suite.DB().Create(&rs))
	user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: models.User{
			Active: true,
			Roles:  []roles.Role{rs[1]},
		},
	})
	ta := tooRoleChecker{suite.DB(), suite.logger}
	ui, err := ta.FetchUserIdentity(&user)
	suite.NoError(err)
	role, err := ta.VerifyHasTOORole(ui)

	suite.NoError(err)
	suite.Equal(role.RoleType, roles.RoleTypeTOO)
}

func (suite *AuthSuite) TestVerifyHasTOORoleError() {
	rs := []roles.Role{{
		ID:       uuid.FromStringOrNil("ed2d2cd7-d427-412a-98bb-a9b391d98d32"),
		RoleType: roles.RoleTypeCustomer,
	},
		{
			ID:       uuid.FromStringOrNil("9dc423b6-33b8-493a-a59b-6a823660cb07"),
			RoleType: roles.RoleTypeTOO,
		},
	}
	suite.NoError(suite.DB().Create(&rs))
	user := testdatagen.MakeDefaultUser(suite.DB())
	ta := tooRoleChecker{suite.DB(), suite.logger}
	ui, err := ta.FetchUserIdentity(&user)
	suite.NoError(err)
	_, err = ta.VerifyHasTOORole(ui)

	suite.Error(err)
}
