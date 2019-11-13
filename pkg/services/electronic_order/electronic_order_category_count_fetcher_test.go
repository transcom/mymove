package electronicorder

import (
	"errors"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testElectronicOrderCategoricalCountQueryBuilder struct {
	fakeFetchCategoricalCountsFromOneModel func(model interface{}) (map[interface{}]int, error)
}

func (t *testElectronicOrderCategoricalCountQueryBuilder) FetchCategoricalCountsFromOneModel(model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	m, err := t.fakeFetchCategoricalCountsFromOneModel(model)
	return m, err
}

func (suite *ElectronicOrderServiceSuite) TestFetchElectronicOrderCategoricalCounts() {
	suite.T().Run("If we get a match on the category we should get a map with the count", func(t *testing.T) {

		fakeFetchCategoricalCountsFromOneModel := func(model interface{}) (map[interface{}]int, error) {
			value := map[interface{}]int{
				models.IssuerArmy: 3,
			}
			return value, nil
		}

		builder := &testElectronicOrderCategoricalCountQueryBuilder{
			fakeFetchCategoricalCountsFromOneModel: fakeFetchCategoricalCountsFromOneModel,
		}

		fetcher := NewElectronicOrdersCategoricalCountsFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("issuer", "=", models.IssuerArmy),
		}

		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(filters, nil)
		suite.NoError(err)
		suite.Equal(3, counts[models.IssuerArmy])
	})

	suite.T().Run("If there's an error, we get it without counts", func(t *testing.T) {
		fakeFetchCategoricalCountsFromOneModel := func(model interface{}) (map[interface{}]int, error) {
			return nil, errors.New("Fetch error")
		}

		builder := &testElectronicOrderCategoricalCountQueryBuilder{
			fakeFetchCategoricalCountsFromOneModel: fakeFetchCategoricalCountsFromOneModel,
		}

		fetcher := NewElectronicOrdersCategoricalCountsFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("issuer", "=", models.IssuerArmy),
		}

		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(filters, nil)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Nil(counts)

	})

}
