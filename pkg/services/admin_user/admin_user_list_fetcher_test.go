package adminuser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testAdminUserListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testAdminUserListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, pagination services.Pagination) error {
	m := t.fakeFetchMany(model)
	return m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func (suite *AdminUserServiceSuite) TestFetchAdminUserList() {
	suite.T().Run("if the users are successfully fetched, they should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.AdminUser{ID: id})))
			return nil
		}
		builder := &testAdminUserListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewAdminUserListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		adminUsers, err := fetcher.FetchAdminUserList(filters, defaultPagination())

		suite.NoError(err)
		suite.Equal(id, adminUsers[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no admin users", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testAdminUserListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewAdminUserListFetcher(builder)

		adminUsers, err := fetcher.FetchAdminUserList([]services.QueryFilter{}, defaultPagination())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.AdminUsers(nil), adminUsers)
	})
}
