package clientcert

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/notifications"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func setUpMockNotificationSender() notifications.NotificationSender {
	// The ClientCertCreator needs a NotificationSender for sending user activity emails to system admins.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.ClientCertModified"),
	).Return(nil)

	return &mockSender
}

func (suite *ClientCertServiceSuite) TestCreateClientCert() {
	queryBuilder := query.NewQueryBuilder()
	suite.Run("If the client cert is created successfully it should be returned", func() {
		builder := &testClientCertQueryBuilder{
			fakeCreateOne: queryBuilder.CreateOne,
		}
		mockSender := setUpMockNotificationSender()

		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{})

		clientCertInfo := models.ClientCert{
			Subject:      "fake subject",
			Sha256Digest: "fake digest",
			UserID:       user.ID,
		}

		creator := NewClientCertCreator(builder, mockSender)
		clientCert, verrs, err := creator.CreateClientCert(suite.AppContextWithSessionForTest(&auth.Session{}), &clientCertInfo)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(clientCert.ID)
		suite.Equal(clientCert.Subject, clientCertInfo.Subject)
		suite.Equal(clientCert.Sha256Digest, clientCertInfo.Sha256Digest)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
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
			fakeCreateOne: fakeCreateOne,
		}

		clientCertInfo := models.ClientCert{
			Subject:      "fake subject",
			Sha256Digest: "fake digest",
		}

		creator := NewClientCertCreator(builder, setUpMockNotificationSender())
		_, verrs, _ := creator.CreateClientCert(suite.AppContextForTest(),
			&clientCertInfo)
		suite.NotNil(verrs)
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
			fakeCreateOne: fakeCreateOne,
		}

		clientCertInfo := models.ClientCert{
			Subject:      "fake subject",
			Sha256Digest: "fake digest",
		}

		creator := NewClientCertCreator(builder, setUpMockNotificationSender())
		_, _, err := creator.CreateClientCert(suite.AppContextForTest(),
			&clientCertInfo)
		suite.EqualError(err, "uniqueness constraint conflict")
	})
}
