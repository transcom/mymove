package payloads

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestUserModel() {
	userID := uuid.Must(uuid.NewV4())

	suite.Run("success - complete input", func() {
		oktaEmail := "user@example.com"
		activeStatus := true

		inputUser := &adminmessages.UserUpdate{
			OktaEmail: &oktaEmail,
			Active:    &activeStatus,
		}

		expectedUser := &models.User{
			ID:        userID,
			OktaEmail: oktaEmail,
			Active:    activeStatus,
		}

		returnedUser, err := UserModel(inputUser, userID, false)

		suite.NoError(err)
		suite.NotNil(returnedUser)
		suite.Equal(expectedUser.ID, returnedUser.ID)
		suite.Equal(expectedUser.OktaEmail, returnedUser.OktaEmail)
		suite.Equal(expectedUser.Active, returnedUser.Active)
	})

	suite.Run("success - Active is nil (defaults to original active value)", func() {
		oktaEmail := "user@example.com"

		inputUser := &adminmessages.UserUpdate{
			OktaEmail: &oktaEmail,
			Active:    nil, // active status not provided
		}

		expectedUser := &models.User{
			ID:        userID,
			OktaEmail: oktaEmail,
			Active:    true, // userOriginalActive to be used
		}

		returnedUser, err := UserModel(inputUser, userID, true)

		suite.NoError(err)
		suite.NotNil(returnedUser)
		suite.Equal(expectedUser.ID, returnedUser.ID)
		suite.Equal(expectedUser.OktaEmail, returnedUser.OktaEmail)
		suite.Equal(expectedUser.Active, returnedUser.Active)
	})

	suite.Run("Error - Nil input", func() {
		returnedUser, err := UserModel(nil, userID, true)

		suite.Nil(returnedUser)
		suite.EqualError(err, "user payload is nil")
	})
}
