package tsp

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

type testTransportationServiceProviderPerformanceListQueryBuilder struct {
	fakeFetchMany func(appCtx appcontext.AppContext, model interface{}) error
	fakeCount     func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testTransportationServiceProviderPerformanceListQueryBuilder) FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(appCtx, model)
	return m
}

func (t *testTransportationServiceProviderPerformanceListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error) {
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

func (suite *TSPServiceSuite) TestFetchTSPPList() {
	suite.Run("if the TSPP is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.TransportationServiceProviderPerformance{ID: id})))
			return nil
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList(suite.AppContextForTest(), filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, tspps[0].ID)
	})

	suite.Run("if TSPPs are fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		id2, err2 := uuid.NewV4()

		suite.NoError(err)
		suite.NoError(err2)
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.TransportationServiceProviderPerformance{ID: id})))
			value.Set(reflect.Append(value, reflect.ValueOf(models.TransportationServiceProviderPerformance{ID: id2})))
			return nil
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList(suite.AppContextForTest(), filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Len(tspps, 2)
	})

	suite.Run("if there is an error, we get it with no tspps", func() {
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList(suite.AppContextForTest(), []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationServiceProviderPerformances(nil), tspps)
	})
}

func (suite *TSPServiceSuite) TestCountTSPPs() {

	suite.Run("if TSPPs are found, they should be counted", func() {
		id, err := uuid.NewV4()

		suite.NoError(err)
		fakeCount := func(appCtx appcontext.AppContext, model interface{}) (int, error) {
			count := 2
			return count, nil
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeCount: fakeCount,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		count, err := fetcher.FetchTransportationServiceProviderPerformanceCount(suite.AppContextForTest(), filters)

		suite.NoError(err)
		suite.Equal(2, count)
	})

	suite.Run("if there is an error, we get it with no count", func() {
		fakeCount := func(appCtx appcontext.AppContext, model interface{}) (int, error) {
			return 0, errors.New("Fetch error")
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeCount: fakeCount,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)

		count, err := fetcher.FetchTransportationServiceProviderPerformanceCount(suite.AppContextForTest(), []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(0, count)
	})
}
