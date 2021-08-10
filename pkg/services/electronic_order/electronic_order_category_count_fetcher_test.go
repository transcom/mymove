package electronicorder

import (
	"errors"
	"testing"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testElectronicOrderCategoricalCountQueryBuilder struct {
	fakeFetchCategoricalCountsFromOneModel func(appCfg appconfig.AppConfig, model interface{}) (map[interface{}]int, error)
}

func (t *testElectronicOrderCategoricalCountQueryBuilder) FetchCategoricalCountsFromOneModel(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	m, err := t.fakeFetchCategoricalCountsFromOneModel(appCfg, model)
	return m, err
}

func (suite *ElectronicOrderServiceSuite) TestFetchElectronicOrderCategoricalCounts() {
	suite.T().Run("If we get a match on the category we should get a map with the count", func(t *testing.T) {

		fakeFetchCategoricalCountsFromOneModel := func(appCfg appconfig.AppConfig, model interface{}) (map[interface{}]int, error) {
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

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(appCfg, filters, nil)
		suite.NoError(err)
		suite.Equal(3, counts[models.IssuerArmy])
	})

	suite.T().Run("If there's an error, we get it without counts", func(t *testing.T) {
		fakeFetchCategoricalCountsFromOneModel := func(appCfg appconfig.AppConfig, model interface{}) (map[interface{}]int, error) {
			return nil, errors.New("Fetch error")
		}

		builder := &testElectronicOrderCategoricalCountQueryBuilder{
			fakeFetchCategoricalCountsFromOneModel: fakeFetchCategoricalCountsFromOneModel,
		}

		fetcher := NewElectronicOrdersCategoricalCountsFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("issuer", "=", models.IssuerArmy),
		}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		counts, err := fetcher.FetchElectronicOrderCategoricalCounts(appCfg, filters, nil)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Nil(counts)

	})

}
