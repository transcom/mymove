package adminuser

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
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

func (suite *RequestedOfficeUsersServiceSuite) TestFetchRequestedOfficeUserPop() {
	suite.Run("returns office user on success", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		fetcher := NewRequestedOfficeUserFetcherPop()

		fetchedUser, err := fetcher.FetchRequestedOfficeUserByID(suite.AppContextForTest(), officeUser.ID)

		suite.NoError(err)
		suite.Equal(officeUser.ID, fetchedUser.ID)
	})

	suite.Run("returns zero value office user on error", func() {
		fetcher := NewRequestedOfficeUserFetcherPop()
		officeUser, err := fetcher.FetchRequestedOfficeUserByID(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(uuid.Nil, officeUser.ID)
	})
}
