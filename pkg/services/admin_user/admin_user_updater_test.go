package adminuser

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AdminUserServiceSuite) TestUpdateAdminUser() {
	newUUID, _ := uuid.NewV4()

	firstName := "Leo"
	payload := &adminmessages.AdminUserUpdatePayload{
		FirstName: &firstName,
	}

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return nil
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewAdminUserUpdater(builder)
		_, verrs, err := updater.UpdateAdminUser(suite.AppContextForTest(), newUUID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
	})

	// Bad organization ID
	suite.Run("If we are provided a organization that doesn't exist, the create should fail", func() {
		fakeUpdateOne := func(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewAdminUserUpdater(builder)
		_, _, err := updater.UpdateAdminUser(suite.AppContextForTest(), newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
