package adminuser

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testRequestedOfficeUsersQueryBuilder struct {
	fakeFetchOne func(appConfig appcontext.AppContext, model interface{}) error
}

func (t *testRequestedOfficeUsersQueryBuilder) FetchOne(appConfig appcontext.AppContext, model interface{}, _ []services.QueryFilter) error {
	m := t.fakeFetchOne(appConfig, model)
	return m
}

func (t *testRequestedOfficeUsersQueryBuilder) UpdateOne(_ appcontext.AppContext, _ interface{}, _ *string) (*validate.Errors, error) {
	return nil, nil
}

func (t *testRequestedOfficeUsersQueryBuilder) DeleteOne(_ appcontext.AppContext, _ interface{}) error {
	return nil
}

func (t *testRequestedOfficeUsersQueryBuilder) DeleteMany(_ appcontext.AppContext, _ interface{}, _ []services.QueryFilter) error {
	return nil
}

func (suite *RequestedOfficeUsersServiceSuite) TestFetchRequestedOfficeUser() {
	suite.Run("if the requested office user is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(_ appcontext.AppContext, model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testRequestedOfficeUsersQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		fetcher := NewRequestedOfficeUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		requestedOfficeUser, err := fetcher.FetchRequestedOfficeUser(suite.AppContextForTest(), filters)

		suite.NoError(err)
		suite.Equal(id, requestedOfficeUser.ID)
	})

	suite.Run("if there is an error, we get it with zero admin user", func() {
		fakeFetchOne := func(_ appcontext.AppContext, _ interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testRequestedOfficeUsersQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewRequestedOfficeUserFetcher(builder)

		requestedOfficeUser, err := fetcher.FetchRequestedOfficeUser(suite.AppContextForTest(), []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUser{}, requestedOfficeUser)
	})
}
