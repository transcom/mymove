package office

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeQueryBuilder struct {
	fakeFetchOne func(appCfg appconfig.AppConfig, model interface{}) error
}

func (t *testOfficeQueryBuilder) FetchOne(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(appCfg, model)
	return m
}

func (suite *OfficeServiceSuite) TestFetchOffice() {
	suite.T().Run("if the transportation office is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(appCfg appconfig.AppConfig, model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		builder := &testOfficeQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		office, err := fetcher.FetchOffice(appCfg, filters)

		suite.NoError(err)
		suite.Equal(id, office.ID)
	})

	suite.T().Run("if there is an error, we get it with zero office", func(t *testing.T) {
		fakeFetchOne := func(appCfg appconfig.AppConfig, model interface{}) error {
			return errors.New("Fetch error")
		}
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		builder := &testOfficeQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeFetcher(builder)

		office, err := fetcher.FetchOffice(appCfg, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationOffice{}, office)
	})
}
