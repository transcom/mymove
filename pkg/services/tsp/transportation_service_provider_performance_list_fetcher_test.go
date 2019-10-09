package tsp

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

type testTransportationServiceProviderPerformanceListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testTransportationServiceProviderPerformanceListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error {
	m := t.fakeFetchMany(model)
	return m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func defaultAssociations() services.QueryAssociations {
	return query.NewQueryAssociations([]services.QueryAssociation{})
}

func (suite *TSPServiceSuite) TestFetchTSPPList() {
	suite.T().Run("if the TSPP is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
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

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList(filters, defaultAssociations(), defaultPagination())

		suite.NoError(err)
		suite.Equal(id, tspps[0].ID)
	})

	suite.T().Run("if TSPPs are fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		id2, err := uuid.NewV4()

		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
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

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList(filters, defaultAssociations(), defaultPagination())

		suite.NoError(err)
		suite.Len(tspps, 2)
	})

	suite.T().Run("if there is an error, we get it with no tspps", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testTransportationServiceProviderPerformanceListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewTransportationServiceProviderPerformanceListFetcher(builder)

		tspps, err := fetcher.FetchTransportationServiceProviderPerformanceList([]services.QueryFilter{}, defaultAssociations(), defaultPagination())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationServiceProviderPerformances(nil), tspps)
	})
}
