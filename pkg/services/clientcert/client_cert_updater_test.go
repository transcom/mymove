package clientcert

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ClientCertServiceSuite) TestUpdateClientCert() {
	newUUID, _ := uuid.NewV4()

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

		updater := NewClientCertUpdater(builder)
		_, verrs, err := updater.UpdateClientCert(suite.AppContextForTest(), newUUID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
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

		updater := NewClientCertUpdater(builder)
		_, _, err := updater.UpdateClientCert(suite.AppContextForTest(), newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
