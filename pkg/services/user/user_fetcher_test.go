package user

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

type testUserQueryBuilder struct {
	fakeFetchOne  func(model interface{}) error
	fakeUpdateOne func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testUserQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (t *testUserQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return nil, nil
}

func (suite *UserServiceSuite) TestFetchUser() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		fetcher := NewUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		user, err := fetcher.FetchUser(filters)

		suite.NoError(err)
		suite.Equal(id, user.ID)
	})

	suite.T().Run("if there is an error, we get it with zero user", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewUserFetcher(builder)

		user, err := fetcher.FetchUser([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.User{}, user)
	})
}
