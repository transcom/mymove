//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package user

import (
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	adminUser "github.com/transcom/mymove/pkg/services/admin_user"
	officeUser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
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

	suite.Run("Deactivate a user successfully", func() {
		// This case should send an email to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		user := testdatagen.MakeDefaultUser(suite.DB())
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

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		officeUser := testdatagen.MakeActiveOfficeUser(suite.DB())

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

		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		adminUser := testdatagen.MakeActiveAdminUser(suite.DB())

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

	suite.Run("Activate a user successfully", func() {
		// Under test: updateUser
		// Mocked:     notificationSender
		// Set up:     We provide an inactive user, and then activate them
		//
		// Expected outcome:
		//           	updateUser updates the user to active
		//              A notification is sent to sys admins
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		// Make an inactive user
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{})

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
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		user := testdatagen.MakeDefaultUser(suite.DB()) // Default user is active
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
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{}) // MakeUser makes an inactive user

		user.Active = inactiveStatus
		updatedUser, verr, err := updater.UpdateUser(appCtx, user.ID, &user)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 0)
	})
}
