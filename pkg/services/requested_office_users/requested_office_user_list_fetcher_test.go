package adminuser

import (
	"errors"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testRequestedOfficeUsersListQueryBuilder struct {
	fakeFetchMany func(appCtx appcontext.AppContext, model interface{}) error
	fakeCount     func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testRequestedOfficeUsersListQueryBuilder) FetchMany(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter, _ services.QueryAssociations, _ services.Pagination, _ services.QueryOrder) error {
	m := t.fakeFetchMany(appCtx, model)
	return m
}

func (t *testRequestedOfficeUsersListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter) (int, error) {
	count, m := t.fakeCount(appCtx, model)
	return count, m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func defaultAssociations() services.QueryAssociations {
	return query.NewQueryAssociations([]services.QueryAssociation{})
}

func defaultOrdering() services.QueryOrder {
	return query.NewQueryOrder(nil, nil)
}

func (suite *RequestedOfficeUsersServiceSuite) TestFetchRequestedOfficeUserList() {
	suite.Run("if the users are successfully fetched, they should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(_ appcontext.AppContext, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			requestedStatus := models.OfficeUserStatusREQUESTED
			value.Set(reflect.Append(value, reflect.ValueOf(models.OfficeUser{ID: id, Status: &requestedStatus})))
			return nil
		}
		builder := &testRequestedOfficeUsersListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewRequestedOfficeUsersListFetcher(builder)

		requestedOfficeUsers, err := fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, requestedOfficeUsers[0].ID)
	})

	suite.Run("if there is an error, we get it with no requested office users", func() {
		fakeFetchMany := func(_ appcontext.AppContext, _ interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testRequestedOfficeUsersListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewRequestedOfficeUsersListFetcher(builder)

		requestedOfficeUsers, err := fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUsers(nil), requestedOfficeUsers)
	})
}
