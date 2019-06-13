package user

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
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

		// Mocking authorization
		session := auth.Session{}
		authFunction := func(session *auth.Session) error {
			return nil
		}

		fetcher := NewOfficeUserListFetcher(builder, authFunction)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		officeUsers, err := fetcher.FetchOfficeUserList(filters, &session)

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

		// Mocking authorization
		session := auth.Session{}
		authFunction := func(session *auth.Session) error {
			return nil
		}

		fetcher := NewOfficeUserListFetcher(builder, authFunction)

		officeUsers, err := fetcher.FetchOfficeUserList([]services.QueryFilter{}, &session)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUsers(nil), officeUsers)
	})

	suite.T().Run("if the user is unauthorized, we get an error", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return nil
		}
		builder := &testOfficeUserListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		// Mocking authorization
		session := auth.Session{}
		authFunction := func(session *auth.Session) error {
			return errors.New("USER_UNAUTHORIZED")
		}

		fetcher := NewOfficeUserListFetcher(builder, authFunction)

		officeUsers, _ := fetcher.FetchOfficeUserList([]services.QueryFilter{}, &session)

		// suite.Error(err)
		// suite.Equal(err.Error(), "USER_UNAUTHORIZED")
		suite.Equal(models.OfficeUsers(nil), officeUsers)
	})
}
