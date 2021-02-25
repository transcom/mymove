package user

import (
	"testing"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	adminUser "github.com/transcom/mymove/pkg/services/admin_user"
	officeUser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UserServiceSuite) TestUserUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	officeUserUpdater := officeUser.NewOfficeUserUpdater(builder)
	adminUserUpdater := adminUser.NewAdminUserUpdater(builder)
	updater := NewUserUpdater(builder, officeUserUpdater, adminUserUpdater)

	active := true
	inactive := false

	activeUser := testdatagen.MakeDefaultUser(suite.DB())

	suite.T().Run("Deactivate a user successfully", func(t *testing.T) {
		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}
		modelToPayload, _ := payloads.UserModel(&payload, activeUser.ID)
		// Take our existing active user and change their Active status to False
		updatedUser, verr, err := updater.UpdateUser(activeUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedUser.Active)

	})

	suite.T().Run("Deactivate an Office User successfully", func(t *testing.T) {
		// Create an active office user
		// deactivate the user
		// check that their office user status is also set to false

		activeOfficeUser := testdatagen.MakeActiveOfficeUser(suite.DB())

		// Create the payload to update a user's active status. This should also
		// update their officeUser status in parallel.
		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}

		modelToPayload, _ := payloads.UserModel(&payload, *activeOfficeUser.UserID)

		// Deactivate user
		updatedUser, verr, err := updater.UpdateUser(*activeOfficeUser.UserID, modelToPayload)

		// Fetch updated office user to confirm status
		updatedOfficeUser := models.OfficeUser{}
		suite.DB().Eager("OfficeUser.User").Find(&updatedOfficeUser, activeOfficeUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedOfficeUser.Active)
		suite.False(updatedUser.Active)

	})

	suite.T().Run("Deactivate an Admin User successfully", func(t *testing.T) {
		// Create an active office user
		// deactivate the user
		// check that their office user status is also set to false

		activeAdminUser := testdatagen.MakeActiveAdminUser(suite.DB())

		// Create the payload to update a user's active status. This should also
		// update their officeUser status in parallel.
		payload := adminmessages.UserUpdatePayload{
			Active: &inactive,
		}

		modelToPayload, _ := payloads.UserModel(&payload, *activeAdminUser.UserID)

		// Deactivate user
		updatedUser, verr, err := updater.UpdateUser(*activeAdminUser.UserID, modelToPayload)

		// Fetch updated office user to confirm status
		updatedAdminUser := models.AdminUser{}
		suite.DB().Eager("AdminUser.User").Find(&updatedAdminUser, activeAdminUser.ID)

		// Check that there are no errors and both statuses successfully updated
		suite.Nil(verr)
		suite.Nil(err)
		suite.False(updatedAdminUser.Active)
		suite.False(updatedUser.Active)

	})

	suite.T().Run("Activate a user successfully", func(t *testing.T) {
		payload := adminmessages.UserUpdatePayload{
			Active: &active,
		}
		modelToPayload, _ := payloads.UserModel(&payload, activeUser.ID)
		// Take our existing inactive user and change their Active status to True
		updatedUser, verr, err := updater.UpdateUser(activeUser.ID, modelToPayload)

		suite.Nil(verr)
		suite.Nil(err)
		suite.True(updatedUser.Active)

	})

}
