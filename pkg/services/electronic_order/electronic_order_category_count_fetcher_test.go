package electronicorder

import (
	"errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testElectronicOrderCategoricalCountQueryBuilder struct {
	fakeFetchCategoricalCountsFromOneModel func(appCtx appcontext.AppContext, model interface{}) (map[interface{}]int, error)
}

func (t *testElectronicOrderCategoricalCountQueryBuilder) FetchCategoricalCountsFromOneModel(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	m, err := t.fakeFetchCategoricalCountsFromOneModel(appCtx, model)
	return m, err
}

func (suite *ElectronicOrderServiceSuite) TestFetchElectronicOrderCategoricalCounts() {
	suite.Run("If we get a match on the category we should get a map with the count", func() {

		fakeFetchCategoricalCountsFromOneModel := func(appCtx appcontext.AppContext, model interface{}) (map[interface{}]int, error) {
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

		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(suite.AppContextForTest(), filters, nil)
		suite.NoError(err)
		suite.Equal(3, counts[models.IssuerArmy])
	})

	suite.Run("If there's an error, we get it without counts", func() {
		fakeFetchCategoricalCountsFromOneModel := func(appCtx appcontext.AppContext, model interface{}) (map[interface{}]int, error) {
			return nil, errors.New("Fetch error")
		}

		builder := &testElectronicOrderCategoricalCountQueryBuilder{
			fakeFetchCategoricalCountsFromOneModel: fakeFetchCategoricalCountsFromOneModel,
		}

		fetcher := NewElectronicOrdersCategoricalCountsFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("issuer", "=", models.IssuerArmy),
		}

		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(suite.AppContextForTest(), filters, nil)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Nil(counts)

	})

}
