package payloads

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
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

func (suite *PayloadsSuite) TestOfficeUserModelFromUpdate() {
	suite.Run("success - complete input", func() {
		email := "johntest@test.com"
		firstName := "firstNameTest"
		middleInitials := "Z"
		lastName := "lastNameTest"
		telephone := "111-111-1111"

		payload := &adminmessages.OfficeUserUpdate{
			Email:          &email,
			FirstName:      &firstName,
			MiddleInitials: &middleInitials,
			LastName:       &lastName,
			Telephone:      &telephone,
			Active:         models.BoolPointer(false),
		}

		oldMiddleInitials := "H"

		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:      "John",
					LastName:       "Doe",
					MiddleInitials: &oldMiddleInitials,
					Telephone:      "555-555-5555",
					Email:          "johndoe@example.com",
					Active:         true,
				},
			},
		}, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(oldOfficeUser.ID, returnedOfficeUser.ID)
		suite.Equal(oldOfficeUser.UserID, returnedOfficeUser.UserID)
		suite.Equal(*payload.Email, returnedOfficeUser.Email)
		suite.Equal(*payload.FirstName, returnedOfficeUser.FirstName)
		suite.Equal(*payload.MiddleInitials, *returnedOfficeUser.MiddleInitials)
		suite.Equal(*payload.LastName, returnedOfficeUser.LastName)
		suite.Equal(*payload.Telephone, returnedOfficeUser.Telephone)
		suite.Equal(false, returnedOfficeUser.Active)
	})

	suite.Run("success - only update Telephone", func() {
		telephone := "111-111-1111"
		payload := &adminmessages.OfficeUserUpdate{
			Telephone: &telephone,
		}

		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Telephone: "555-555-5555",
				},
			},
		}, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(oldOfficeUser.ID, returnedOfficeUser.ID)
		suite.Equal(oldOfficeUser.UserID, returnedOfficeUser.UserID)
		suite.Equal(oldOfficeUser.Email, returnedOfficeUser.Email)
		suite.Equal(oldOfficeUser.FirstName, returnedOfficeUser.FirstName)
		suite.Equal(oldOfficeUser.MiddleInitials, returnedOfficeUser.MiddleInitials)
		suite.Equal(oldOfficeUser.LastName, returnedOfficeUser.LastName)
		suite.Equal(*payload.Telephone, returnedOfficeUser.Telephone)
		suite.Equal(oldOfficeUser.Active, returnedOfficeUser.Active)
	})

	suite.Run("Fields do not update if payload is empty", func() {
		payload := &adminmessages.OfficeUserUpdate{}

		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(oldOfficeUser.ID, returnedOfficeUser.ID)
		suite.Equal(oldOfficeUser.UserID, returnedOfficeUser.UserID)
		suite.Equal(oldOfficeUser.Email, returnedOfficeUser.Email)
		suite.Equal(oldOfficeUser.FirstName, returnedOfficeUser.FirstName)
		suite.Equal(oldOfficeUser.MiddleInitials, returnedOfficeUser.MiddleInitials)
		suite.Equal(oldOfficeUser.LastName, returnedOfficeUser.LastName)
		suite.Equal(oldOfficeUser.Telephone, returnedOfficeUser.Telephone)
		suite.Equal(oldOfficeUser.Active, returnedOfficeUser.Active)
	})

	suite.Run("Error - Return Office User if payload is nil", func() {
		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		returnedUser := OfficeUserModelFromUpdate(nil, &oldOfficeUser)

		suite.Equal(&oldOfficeUser, returnedUser)
	})

	suite.Run("Error - Return nil if Office User is nil", func() {
		telephone := "111-111-1111"
		payload := &adminmessages.OfficeUserUpdate{
			Telephone: &telephone,
		}
		returnedUser := OfficeUserModelFromUpdate(payload, nil)

		suite.Nil(returnedUser)
	})
}
