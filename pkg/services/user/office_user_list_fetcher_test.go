package user

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeUserListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testOfficeUserListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchMany(model)
	return m
}

func (suite *UserServiceSuite) TestFetchOfficeUserList() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.OfficeUser{ID: id})))
			return nil
		}
		builder := &testOfficeUserListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}
		fetcher := NewOfficeUserListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		officeUsers, err := fetcher.FetchOfficeUserList(filters)

		suite.NoError(err)
		suite.Equal(id, officeUsers[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no office users", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeUserListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}
		fetcher := NewOfficeUserListFetcher(builder)

		officeUsers, err := fetcher.FetchOfficeUserList([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUsers(nil), officeUsers)
	})
}
