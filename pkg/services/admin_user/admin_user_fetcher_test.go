package adminuser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testAdminUserQueryBuilder struct {
	fakeFetchOne  func(model interface{}) error
	fakeCreateOne func(models interface{}) (*validate.Errors, error)
	fakeUpdateOne func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testAdminUserQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (t *testAdminUserQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return nil, nil
}

func (t *testAdminUserQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return nil, nil
}

func (suite *AdminUserServiceSuite) TestFetchAdminUser() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		fetcher := NewAdminUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		adminUser, err := fetcher.FetchAdminUser(filters)

		suite.NoError(err)
		suite.Equal(id, adminUser.ID)
	})

	suite.T().Run("if there is an error, we get it with zero admin user", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewAdminUserFetcher(builder)

		adminUser, err := fetcher.FetchAdminUser([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.AdminUser{}, adminUser)
	})
}
