package clientcert

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications/mocks"
)

func (suite *ClientCertServiceSuite) TestUpdateClientCert() {
	newUUID, _ := uuid.NewV4()
	mockSender := setUpMockNotificationSender()

	allowPrime := true
	payload := &adminmessages.ClientCertUpdatePayload{
		AllowPrime: &allowPrime,
	}

	// Happy path
	suite.Run("If the client cert is updated successfully it should be returned", func() {
		fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return nil
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewClientCertUpdater(builder, mockSender)
		_, verrs, err := updater.UpdateClientCert(suite.AppContextWithSessionForTest(&auth.Session{}), newUUID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})

	// Bad cert ID
	suite.Run("If we are provided an id that doesn't exist, the update should fail", func() {
		fakeUpdateOne := func(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewClientCertUpdater(builder, mockSender)
		_, _, err := updater.UpdateClientCert(suite.AppContextWithSessionForTest(&auth.Session{}), newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
