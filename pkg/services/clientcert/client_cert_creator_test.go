package clientcert

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	notification_mocks "github.com/transcom/mymove/pkg/notifications/mocks"
	services_mocks "github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The ClientCertCreator needs a NotificationSender for sending user activity emails to system admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := notification_mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.ClientCertModified"),
	).Return(nil)

	return &mockSender
}

func (suite *ClientCertServiceSuite) TestCreateClientCert() {
	hash := sha256.Sum256([]byte("fake"))
	digest := hex.EncodeToString(hash[:])
	queryBuilder := query.NewQueryBuilder()

	suite.Run("Create clientcert with prime to existing user", func() {
		associator := usersroles.NewUsersRolesCreator()
		mockSender := setUpMockNotificationSender()

		user := factory.BuildUser(suite.DB(), nil, nil)
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCertInfo := models.ClientCert{
			Subject:      "existingUser",
			Sha256Digest: digest,
			UserID:       user.ID,
			AllowPrime:   true,
		}

		creator := NewClientCertCreator(queryBuilder, associator, mockSender)
		clientCert, verrs, err := creator.CreateClientCert(
			suite.AppContextWithSessionForTest(&auth.Session{}),
			user.OktaEmail, &clientCertInfo)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(clientCert.ID)
		suite.Equal(clientCert.Subject, clientCertInfo.Subject)
		suite.Equal(clientCert.Sha256Digest, clientCertInfo.Sha256Digest)
		suite.Equal(clientCert.UserID, user.ID)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), user.ID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))
		mockSender.(*notification_mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	suite.Run("Create clientcert with prime to new user", func() {
		associator := usersroles.NewUsersRolesCreator()
		mockSender := setUpMockNotificationSender()

		// make sure  the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCertInfo := models.ClientCert{
			Subject:      "newUser",
			Sha256Digest: digest,
			AllowPrime:   true,
		}

		creator := NewClientCertCreator(queryBuilder, associator, mockSender)
		clientCert, verrs, err := creator.CreateClientCert(
			suite.AppContextWithSessionForTest(&auth.Session{}),
			"newuser@example.com", &clientCertInfo)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(clientCert.ID)
		suite.Equal(clientCert.Subject, clientCertInfo.Subject)
		suite.Equal(clientCert.Sha256Digest, clientCertInfo.Sha256Digest)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))
		mockSender.(*notification_mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Transaction rollback on createOne validation failure
	suite.Run("CreateOne validation error should rollback transaction", func() {
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			// Fail on the ClientCert call to CreateOne
			switch model.(type) {
			case *models.ClientCert:
				return &validate.Errors{
					Errors: map[string][]string{
						"errorKey": {"violation message"},
					},
				}, nil
			default:
				{
					return nil, nil
				}
			}
		}
		builder := &testClientCertQueryBuilder{
			fakeFetchOne:  queryBuilder.FetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		clientCertInfo := models.ClientCert{
			Subject:      "fake subject",
			Sha256Digest: digest,
		}

		associator := &services_mocks.UserRoleAssociator{}
		associator.On("UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.Anything,
		).Return([]models.UsersRoles{}, nil, nil)

		creator := NewClientCertCreator(builder, associator, setUpMockNotificationSender())
		_, verrs, _ := creator.CreateClientCert(suite.AppContextForTest(),
			"fake@example.com", &clientCertInfo)
		suite.NotNil(verrs)
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Errors)
		suite.Equal("violation message", verrs.Errors["errorKey"][0])
	})

	// Transaction rollback on createOne error failure
	suite.Run("CreateOne error should rollback transaction", func() {
		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			// Fail on the createOne call
			switch model.(type) {
			case *models.ClientCert:
				return nil, errors.New("uniqueness constraint conflict")
			default:
				return nil, nil
			}
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne:  queryBuilder.FetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		clientCertInfo := models.ClientCert{
			Subject:      "fake subject",
			Sha256Digest: digest,
		}

		associator := &services_mocks.UserRoleAssociator{}
		associator.On("UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.Anything,
		).Return([]models.UsersRoles{}, nil, nil)
		creator := NewClientCertCreator(builder, associator, setUpMockNotificationSender())
		_, _, err := creator.CreateClientCert(suite.AppContextForTest(),
			"fake@example.com", &clientCertInfo)
		suite.EqualError(err, "uniqueness constraint conflict")
	})
}
