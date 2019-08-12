package tsp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testTransportationServiceProviderPerformanceQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
}

func (t *testTransportationServiceProviderPerformanceQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *TSPServiceSuite) TestFetchTransportationServiceProviderPerformance() {
	suite.T().Run("if the TSPP is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testTransportationServiceProviderPerformanceQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewTransportationServiceProviderPerformanceFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		tspp, err := fetcher.FetchTransportationServiceProviderPerformance(filters)

		suite.NoError(err)
		suite.Equal(id, tspp.ID)
	})

	suite.T().Run("if there is an error, we get it with zero TSPP", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testTransportationServiceProviderPerformanceQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewTransportationServiceProviderPerformanceFetcher(builder)

		tspp, err := fetcher.FetchTransportationServiceProviderPerformance([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationServiceProviderPerformance{}, tspp)
	})
}
