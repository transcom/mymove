package adminuser

import (
	"testing"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

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
	suite.T().Run("If the user is updated successfully it should be returned", func(t *testing.T) {
		fakeUpdateOne := func(interface{}, *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(model interface{}) error {
			return nil
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewAdminUserUpdater(builder)
		_, verrs, err := updater.UpdateAdminUser(newUUID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
	})

	// Bad organization ID
	suite.T().Run("If we are provided a organization that doesn't exist, the create should fail", func(t *testing.T) {
		fakeUpdateOne := func(model interface{}, eTag *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewAdminUserUpdater(builder)
		_, _, err := updater.UpdateAdminUser(newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
