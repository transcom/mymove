package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"

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
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	h := CallbackHandler{
		authContext,
		suite.DB(),
		FakeRSAKey,
		false,
		false,
	}
	rr := httptest.NewRecorder()
	h.SetFeatureFlag(FeatureFlag{Name: cli.FeatureFlagRoleBasedAuth, Active: true})
	h.ServeHTTP(rr, req)

	suite.Equal(rr.Code, 307)
}

func (suite *AuthSuite) TestCreateAndAssociateCustomer() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		Hostname:        MilTestHost,
	}
	ra := customerAssociator{suite.DB(), suite.logger}
	err := ra.CreateAndAssociateCustomer(session.UserID)
	suite.NoError(err)
	c, err := suite.DB().Count(models.Customer{})
	suite.NoError(err)
	customer := &models.Customer{}
	err = suite.DB().Where("user_id=$1", user.ID).First(customer)
	suite.NoError(err)

	suite.Equal(1, c)
	suite.Equal(user.ID, customer.UserID)
}

func (suite *AuthSuite) TestCreateAndAssociateCustomerUserIDNil() {
	session := auth.Session{
		ApplicationName: auth.MilApp,
		Hostname:        MilTestHost,
	}
	ra := customerAssociator{suite.DB(), suite.logger}
	err := ra.CreateAndAssociateCustomer(session.UserID)
	suite.Error(err)
}

func (suite *AuthSuite) TestAssociateOfficeUser() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{OfficeUser: models.OfficeUser{
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
	testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{OfficeUser: models.OfficeUser{
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
	testdatagen.MakeAdminUser(suite.DB(), testdatagen.Assertions{AdminUser: models.AdminUser{
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
	testdatagen.MakeAdminUser(suite.DB(), testdatagen.Assertions{AdminUser: models.AdminUser{
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
	user := goth.User{
		UserID: fakeUUID.String(),
		Email:  "sample@email.com",
	}
	err := uua.AuthorizeUnknownUser(user, &session)
	suite.NoError(err)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserCreateUserFails() {
	oa := &mocks.OfficeUserAssociator{}
	aa := &mocks.AdminUserAssociator{}
	cca := &mocks.CustomerCreatorAndAssociator{}
	ra := roleAssociator{
		db:                           suite.DB(),
		logger:                       suite.logger,
		OfficeUserAssociator:         oa,
		AdminUserAssociator:          aa,
		CustomerCreatorAndAssociator: cca,
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
	cca := &mocks.CustomerCreatorAndAssociator{}
	ra := roleAssociator{
		db:                           suite.DB(),
		logger:                       suite.logger,
		OfficeUserAssociator:         oa,
		AdminUserAssociator:          aa,
		CustomerCreatorAndAssociator: cca,
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
	session := auth.Session{ApplicationName: auth.AdminApp}

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
	aa.AssertNumberOfCalls(suite.T(), "AssociateAdminUser", 1)
	suite.Equal(fakeUUID, session.AdminUserID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsOfficeUser() {
	oa := &mocks.OfficeUserAssociator{}
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	oa.On("AssociateOfficeUser", mock.Anything).Return(fakeUUID, nil)
	aa := &mocks.AdminUserAssociator{}
	cca := &mocks.CustomerCreatorAndAssociator{}
	ra := roleAssociator{
		db:                           suite.DB(),
		logger:                       suite.logger,
		OfficeUserAssociator:         oa,
		AdminUserAssociator:          aa,
		CustomerCreatorAndAssociator: cca,
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
	session := auth.Session{ApplicationName: auth.OfficeApp}
	suite.True(session.IsOfficeApp())

	err := uua.AuthorizeUnknownUser(user, &session)

	suite.NoError(err)
	oa.AssertNumberOfCalls(suite.T(), "AssociateOfficeUser", 1)
	suite.Equal(fakeUUID, session.OfficeUserID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserIsCustomer() {
	oa := &mocks.OfficeUserAssociator{}
	aa := &mocks.AdminUserAssociator{}
	cca := &mocks.CustomerCreatorAndAssociator{}
	cca.On("CreateAndAssociateCustomer", mock.Anything).Return(nil)
	ra := roleAssociator{
		db:                           suite.DB(),
		logger:                       suite.logger,
		OfficeUserAssociator:         oa,
		AdminUserAssociator:          aa,
		CustomerCreatorAndAssociator: cca,
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
	cca.AssertNumberOfCalls(suite.T(), "CreateAndAssociateCustomer", 1)
}
