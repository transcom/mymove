package adminuser

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AdminUserServiceSuite) TestCreateAdminUser() {
	organization := testdatagen.MakeDefaultOrganization(suite.DB())
	userInfo := models.AdminUser{
		LastName:       "Spaceman",
		FirstName:      "Leo",
		Email:          "spaceman@leo.org",
		OrganizationID: &organization.ID,
		Organization:   organization,
	}

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(organization.ID))
			return nil
		}
		fakeCreateOne := func(interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", organization.ID)}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewAdminUserCreator(builder)
		_, verrs, err := creator.CreateAdminUser(&userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)

	})

	// Bad organization ID
	suite.T().Run("If we are provided a organization that doesn't exist, the create should fail", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewAdminUserCreator(builder)
		_, _, err := creator.CreateAdminUser(&userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
