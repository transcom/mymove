package user

import (
	"testing"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UserServiceSuite) TestUserUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewUserUpdater(builder)
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
