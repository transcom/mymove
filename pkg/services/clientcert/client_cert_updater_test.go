package clientcert

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services"
	services_mocks "github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *ClientCertServiceSuite) TestUpdateClientCert() {
	newUUID, _ := uuid.NewV4()
	mockSender := setUpMockNotificationSender()
	queryBuilder := query.NewQueryBuilder()
	hash := sha256.Sum256([]byte("fake"))
	digest := hex.EncodeToString(hash[:])

	// Happy path
	suite.Run("Update client cert removing prime", func() {
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
					AllowPrime:       false,
					AllowPPTAS:       true,
					PPTASAffiliation: (*models.ServiceMemberAffiliation)(models.StringPointer("MARINES")),
				},
			},
		}, nil)
		suite.True(clientCert.AllowPrime)
		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))
		associator := usersroles.NewUsersRolesCreator()
		updater := NewClientCertUpdater(queryBuilder, associator, mockSender)

		// Update the client cert to allow prime and PPTAS, and change PPTAS affiliation to Navy
		allowPrime := true
		allowPPTAS := true
		pptasAffiliation := (*models.ServiceMemberAffiliation)(models.StringPointer("NAVY"))
		payload := &adminmessages.ClientCertificateUpdate{
			Subject:          "new-subject",
			Sha256Digest:     digest,
			AllowPrime:       &allowPrime,
			AllowPPTAS:       &allowPPTAS,
			PptasAffiliation: (*adminmessages.Affiliation)(pptasAffiliation),
		}
		updatedClientCert, verrs, err := updater.UpdateClientCert(
			suite.AppContextWithSessionForTest(&auth.Session{}),
			clientCert.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)

		suite.Equal(payload.Subject, updatedClientCert.Subject)
		suite.Equal(payload.Sha256Digest, updatedClientCert.Sha256Digest)
		suite.Equal(*payload.AllowPrime, updatedClientCert.AllowPrime)
		suite.Equal(*payload.AllowPPTAS, updatedClientCert.AllowPPTAS)
		suite.Equal((*models.ServiceMemberAffiliation)(payload.PptasAffiliation), updatedClientCert.PPTASAffiliation)

		// Update the client cert to remove PPTAS and reset the PPTAS affiliation to nil
		allowPPTAS = false
		pptasAffiliation = nil
		payload = &adminmessages.ClientCertificateUpdate{
			Subject:          "new-subject",
			Sha256Digest:     digest,
			AllowPPTAS:       &allowPPTAS,
			PptasAffiliation: (*adminmessages.Affiliation)(pptasAffiliation),
		}
		updatedClientCert, verrs, err = updater.UpdateClientCert(
			suite.AppContextWithSessionForTest(&auth.Session{}),
			clientCert.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)

		suite.Equal(payload.Subject, updatedClientCert.Subject)
		suite.Equal(payload.Sha256Digest, updatedClientCert.Sha256Digest)
		suite.Equal(*payload.AllowPPTAS, updatedClientCert.AllowPPTAS)
		suite.Nil(updatedClientCert.PPTASAffiliation)

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))

		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Bad cert ID
	suite.Run("If we are provided an id that doesn't exist, the update should fail", func() {
		allowPrime := true
		payload := &adminmessages.ClientCertificateUpdate{
			Subject:      "new-subject",
			Sha256Digest: digest,
			AllowPrime:   &allowPrime,
		}

		fakeUpdateOne := func(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
			return models.ErrFetchNotFound
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		associator := &services_mocks.UserRoleAssociator{}
		associator.On("UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.Anything,
		).Return([]models.UsersRoles{}, nil, nil)
		updater := NewClientCertUpdater(builder, associator, mockSender)
		_, _, err := updater.UpdateClientCert(suite.AppContextWithSessionForTest(&auth.Session{}), newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
