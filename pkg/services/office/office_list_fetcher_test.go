package office

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

type testOfficeListQueryBuilder struct {
	fakeFetchMany func(appCtx appcontext.AppContext, model interface{}) error
	fakeCount     func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testOfficeListQueryBuilder) FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(appCtx, model)
	return m
}

func (t *testOfficeListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error) {
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

func (suite *OfficeServiceSuite) TestFetchOfficeList() {
	suite.Run("if the transportation office is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.TransportationOffice{ID: id})))
			return nil
		}
		builder := &testOfficeListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOfficeListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		offices, err := fetcher.FetchOfficeList(suite.AppContextForTest(), filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, offices[0].ID)
	})

	suite.Run("if there is an error, we get it with no offices", func() {
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOfficeListFetcher(builder)

		offices, err := fetcher.FetchOfficeList(suite.AppContextForTest(), []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationOffices(nil), offices)
	})
}
