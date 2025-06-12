// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package user

import (
	"fmt"
	"net/http/httptest"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	adminUser "github.com/transcom/mymove/pkg/services/admin_user"
	officeUser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/query"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The UserUpdater needs a NotificationSender for sending user activity emails to system admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.UserAccountModified"),
	).Return(nil)

	return &mockSender
}

func (suite *UserServiceSuite) TestUserUpdater() {
	builder := query.NewQueryBuilder()
	officeUserUpdater := officeUser.NewOfficeUserUpdater(builder)
	adminUserUpdater := adminUser.NewAdminUserUpdater(builder)

	activeStatus := true
	inactiveStatus := false

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)

	suite.Run("Deactivate a user successfully", func() {
		// This case should send an email to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		user := factory.BuildDefaultUser(suite.DB())
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		user.Active = inactiveStatus
		// Take our existing active user and change their Active status to False
		updatedUser, verr, err := updater.UpdateUser(appCtx, user.ID, &user)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	suite.Run("Deactivate an Office User successfully", func() {
		// Under test: updateUser, updateOfficeUser
		// Mocked:     Notification sender
		// Set up:     We provide an ACTIVE user/office user, and then deactivate
		//			   the user by calling updateUser.
		//
		// Expected outcome:
		//           	updateUser updates the users table and calls updateOfficeUser
		//            	to update the office_users table. Both tables have an ACTIVE
		//				status set to False.

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
				},
			},
			{
				Model: models.User{
					Active: true,
				},
			},
		}, nil)

		// Deactivate: Update the user with an inactive status. This should also
		// update their officeUser status in parallel.
		// This case should send an email to sys admins
		officeUser.User.Active = inactiveStatus
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *officeUser.UserID, &officeUser.User)

		// Fetch updated office user to confirm status
		updatedOfficeUser := models.OfficeUser{}
		suite.DB().Eager("OfficeUser.User").Find(&updatedOfficeUser, officeUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedOfficeUser.Active)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	suite.Run("Deactivate an Admin User successfully", func() {
		// Under test: updateUser, updateAdminUser
		// Mocked:     notificationSender
		// Set up:     We provide an ACTIVE user/admin user, and then deactivate
		//			   the user by calling updateUser.
		//
		// Expected outcome:
		//           	updateUser updates the users table and calls updateAdminUser
		//            	to update the admin_users table. Both tables have an ACTIVE
		//				status set to False.

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		adminUser := factory.BuildAdminUser(suite.DB(), []factory.Customization{
			{
				Model: models.AdminUser{
					Active: true,
				},
			},
			{
				Model: models.User{
					Active: true,
				},
			},
		}, nil)

		// Deactivate user. This should also update their adminUser status in parallel.
		// This case should send an email to sys admins
		adminUser.User.Active = inactiveStatus
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *adminUser.UserID, &adminUser.User)

		// Fetch updated admin user to confirm status
		updatedAdminUser := models.AdminUser{}
		suite.DB().Eager("AdminUser.User").Find(&updatedAdminUser, adminUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedAdminUser.Active)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	suite.Run("Updates email for user and associated admin user successfully", func() {
		adminUser := factory.BuildAdminUser(suite.DB(), []factory.Customization{
			{
				Model: models.AdminUser{
					Active: true,
					Email:  "adminUser@mail.mil",
				},
			},
			{
				Model: models.User{
					Active:    true,
					OktaEmail: "adminUser@mail.mil",
				},
			},
		}, nil)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%s", adminUser.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		// session.HTTPRequest = request

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		// these mocked endpoints fetch an exact user
		mockAndActivateOktaGETEndpointNoError(adminUser.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(adminUser.User.OktaID)

		// update the email (this should also update admin user email)
		adminUser.User.OktaEmail = "anotherAdminUser@mail.mil"
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *adminUser.UserID, &adminUser.User)

		// Fetch updated admin user to confirm status
		updatedAdminUser := models.AdminUser{}
		suite.DB().Eager("User").Find(&updatedAdminUser, adminUser.ID)

		suite.Nil(verr)
		suite.Nil(err)
		suite.Equal(updatedUser.OktaEmail, updatedAdminUser.Email)
	})

	suite.Run("Updates email for user and associated office user successfully", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					Email:  "officeUser@mail.mil",
				},
			},
			{
				Model: models.User{
					Active:    true,
					OktaEmail: "officeUser@mail.mil",
				},
			},
		}, nil)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%s", officeUser.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		// session.HTTPRequest = request

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		// these mocked endpoints fetch an exact user
		mockAndActivateOktaGETEndpointNoError(officeUser.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(officeUser.User.OktaID)

		// update the email (this should also update admin user email)
		officeUser.User.OktaEmail = "anotherOfficeUser@mail.mil"
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *officeUser.UserID, &officeUser.User)

		// Fetch updated admin user to confirm status
		updatedOfficeUser := models.OfficeUser{}
		suite.DB().Eager("User").Find(&updatedOfficeUser, officeUser.ID)

		suite.Nil(verr)
		suite.Nil(err)
		suite.Equal(updatedUser.OktaEmail, updatedOfficeUser.Email)
	})

	suite.Run("Updates email for user and associated service member user successfully", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					PersonalEmail: models.StringPointer("serviceMember@mail.mil"),
				},
			},
			{
				Model: models.User{
					Active:    true,
					OktaEmail: "serviceMember@mail.mil",
				},
			},
		}, nil)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%s", serviceMember.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		// session.HTTPRequest = request

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		// mocking return responses from okta
		mockAndActivateOktaGETEndpointNoError(serviceMember.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(serviceMember.User.OktaID)

		serviceMember.User.OktaEmail = "anotherServiceMember@mail.mil"
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, serviceMember.UserID, &serviceMember.User)

		updatedServiceMember := models.ServiceMember{}
		suite.DB().Eager("User").Find(&updatedServiceMember, serviceMember.ID)

		suite.Nil(verr)
		suite.Nil(err)
		suite.Equal(&updatedUser.OktaEmail, updatedServiceMember.PersonalEmail)
	})

	suite.Run("Activate a user successfully", func() {
		// Under test: updateUser
		// Mocked:     notificationSender
		// Set up:     We provide an inactive user, and then activate them
		//
		// Expected outcome:
		//           	updateUser updates the user to active
		//              A notification is sent to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		// Make an inactive user
		user := factory.BuildUser(suite.DB(), nil, nil)

		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		// Activate the user
		// Take our existing inactive user and change their Active status to True
		user.Active = activeStatus
		updatedUser, verr, err := updater.UpdateUser(appCtx, user.ID, &user)

		suite.Nil(verr)
		suite.Nil(err)
		suite.True(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	suite.Run("Make no change to active user", func() {
		// Under test: updateUser
		// Mocked:     notificationSender
		// Set up:     We provide an active user, and then "activate" them
		//
		// Expected outcome:
		//           	updateUser returns the active user
		//              A notification is NOT sent to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		user := factory.BuildUser(suite.DB(), nil,
			[]factory.Trait{
				factory.GetTraitActiveUser,
			})

		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		user.Active = activeStatus
		updatedUser, verr, err := updater.UpdateUser(appCtx, user.ID, &user)

		suite.Nil(verr)
		suite.Nil(err)
		suite.True(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 0)
	})

	suite.Run("Make no change to inactive user", func() {
		// Under test: updateUser
		// Mocked:     notificationSender
		// Set up:     We provide an inactive user, and then deactivate them
		//
		// Expected outcome:
		//           	updateUser returns the inactive user
		//              A notification is NOT sent to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		user := factory.BuildUser(suite.DB(), nil, nil)

		user.Active = inactiveStatus
		updatedUser, verr, err := updater.UpdateUser(appCtx, user.ID, &user)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 0)
	})
}

func mockAndActivateOktaGETEndpointNoError(oktaID string) {
	httpmock.Activate()
	getUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := fmt.Sprintf(`{
			"id": "%s",
			"status": "ACTIVE",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}`, oktaID)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}

func mockAndActivateOktaPOSTEndpointNoError(oktaID string) {
	httpmock.Activate()
	updateUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := fmt.Sprintf(`{
			"id": "%s",
			"status": "ACTIVE",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}`, oktaID)

	httpmock.RegisterResponder("POST", updateUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}
