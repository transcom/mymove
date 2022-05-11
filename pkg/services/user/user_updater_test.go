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
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
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
	appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})

	active := true
	inactive := false

	activeUser := testdatagen.MakeDefaultUser(suite.DB())

	suite.T().Run("Deactivate a user successfully", func(t *testing.T) {
		// This case should send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}
		modelToPayload, _ := payloads.UserModel(&payload, activeUser.ID, activeUser.Active)
		// Take our existing active user and change their Active status to False
		updatedUser, verr, err := updater.UpdateUser(appCtx, activeUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 1)
	})

	suite.T().Run("Deactivate an Office User successfully", func(t *testing.T) {
		// Under test: updateUser, updateOfficeUser
		//
		// Set up:     We provide an ACTIVE user/office user, and then deactivate
		//			   the user by calling updateUser.
		//
		// Expected outcome:
		//           	updateUser updates the users table and calls updateOfficeUser
		//            	to update the office_users table. Both tables have an ACTIVE
		//				status set to False.

		activeOfficeUser := testdatagen.MakeActiveOfficeUser(suite.DB())

		// Create the payload to update a user's active status. This should also
		// update their officeUser status in parallel.
		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}

		modelToPayload, _ := payloads.UserModel(&payload, *activeOfficeUser.UserID, activeOfficeUser.Active)

		// Deactivate user
		// This case should send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *activeOfficeUser.UserID, modelToPayload)

		// Fetch updated office user to confirm status
		updatedOfficeUser := models.OfficeUser{}
		suite.DB().Eager("OfficeUser.User").Find(&updatedOfficeUser, activeOfficeUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedOfficeUser.Active)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 1)
	})

	suite.T().Run("Deactivate an Admin User successfully", func(t *testing.T) {
		// Under test: updateUser, updateAdminUser
		//
		// Set up:     We provide an ACTIVE user/admin user, and then deactivate
		//			   the user by calling updateUser.
		//
		// Expected outcome:
		//           	updateUser updates the users table and calls updateAdminUser
		//            	to update the admin_users table. Both tables have an ACTIVE
		//				status set to False.

		activeAdminUser := testdatagen.MakeActiveAdminUser(suite.DB())

		// Create the payload to update a user's active status. This should also
		// update their adminUser status in parallel.
		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}

		modelToPayload, _ := payloads.UserModel(&payload, *activeAdminUser.UserID, activeAdminUser.Active)

		// Deactivate user
		// This case should send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)
		updatedUser, verr, err := updater.UpdateUser(appCtx, *activeAdminUser.UserID, modelToPayload)

		// Fetch updated admin user to confirm status
		updatedAdminUser := models.AdminUser{}
		suite.DB().Eager("AdminUser.User").Find(&updatedAdminUser, activeAdminUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedAdminUser.Active)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 1)
	})

	suite.T().Run("Activate a user successfully", func(t *testing.T) {
		// This case should send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		payload := adminmessages.UserUpdatePayload{
			Active: &active,
		}
		modelToPayload, _ := payloads.UserModel(&payload, activeUser.ID, activeUser.Active)
		// Take our existing inactive user and change their Active status to True
		updatedUser, verr, err := updater.UpdateUser(appCtx, activeUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.True(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 1)
	})

	suite.T().Run("Make no change to active user", func(t *testing.T) {
		// This case should NOT send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		payload := adminmessages.UserUpdatePayload{
			Active: nil,
		}
		modelToPayload, _ := payloads.UserModel(&payload, activeUser.ID, activeUser.Active)
		updatedUser, verr, err := updater.UpdateUser(appCtx, activeUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.True(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 0)
	})

	suite.T().Run("Make no change to inactive user", func(t *testing.T) {
		// This case should NOT send an email to sys admins
		mockSender := setUpMockNotificationSender()
		updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater, mockSender)

		inactiveUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				Active: false,
			},
		})

		payload := adminmessages.UserUpdatePayload{
			Active: nil,
		}
		modelToPayload, _ := payloads.UserModel(&payload, inactiveUser.ID, inactiveUser.Active)
		updatedUser, verr, err := updater.UpdateUser(appCtx, inactiveUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(t, "SendNotification", 0)
	})
}
