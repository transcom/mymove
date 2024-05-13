package clientcert

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *ClientCertServiceSuite) TestRemoveClientCert() {
	mockSender := setUpMockNotificationSender()
	queryBuilder := query.NewQueryBuilder()
	associator := usersroles.NewUsersRolesCreator()

	// Happy path
	suite.Run("Remove client cert", func() {
		clientCert := factory.BuildClientCert(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypePrime,
						},
					},
				},
			},
			{
				Model: models.ClientCert{
					AllowPrime: true,
				},
			},
		}, nil)
		suite.True(clientCert.AllowPrime)
		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))
		updater := NewClientCertRemover(queryBuilder, associator, mockSender)

		payload := &adminmessages.ClientCertificate{
			Subject:      clientCert.Subject,
			Sha256Digest: clientCert.Sha256Digest,
			AllowPrime:   clientCert.AllowPrime,
		}
		removedClientCert, verrs, err := updater.RemoveClientCert(
			suite.AppContextWithSessionForTest(&auth.Session{}),
			clientCert.ID)
		suite.NoError(err)
		suite.Nil(verrs)

		suite.Equal(payload.Subject, removedClientCert.Subject)
		suite.Equal(payload.Sha256Digest, removedClientCert.Sha256Digest)
		suite.Equal(payload.AllowPrime, removedClientCert.AllowPrime)

		var missingClientCert models.ClientCert
		findErr := suite.DB().Find(&missingClientCert, clientCert.ID)
		suite.Equal(sql.ErrNoRows, findErr)

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))

		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Bad cert ID
	suite.Run("If we are provided an id that doesn't exist, the update should fail", func() {
		missingUUID, _ := uuid.NewV4()

		fakeFetchOne := func(_ appcontext.AppContext, _ interface{}, _ []services.QueryFilter) error {
			return models.ErrFetchNotFound
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		updater := NewClientCertRemover(builder, associator, mockSender)
		_, _, err := updater.RemoveClientCert(suite.AppContextWithSessionForTest(&auth.Session{}), missingUUID)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
