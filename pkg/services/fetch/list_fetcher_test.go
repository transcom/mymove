package fetch

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testListQueryBuilder struct {
	fakeFetchMany func(appCfg appconfig.AppConfig, model interface{}) error
	fakeCount     func(appCfg appconfig.AppConfig, model interface{}) (int, error)
}

func (t *testListQueryBuilder) FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(appCfg, model)
	return m
}

func (t *testListQueryBuilder) Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error) {
	count, m := t.fakeCount(appCfg, model)
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

func (suite *FetchServiceSuite) TestFetchRecordList() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(appCfg appconfig.AppConfig, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.OfficeUser{ID: id})))
			return nil
		}
		builder := &testListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		var officeUsers models.OfficeUsers
		err = fetcher.FetchRecordList(appCfg, &officeUsers, filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, officeUsers[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no office users", func(t *testing.T) {
		fakeFetchMany := func(appCfg appconfig.AppConfig, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewListFetcher(builder)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		var officeUsers models.OfficeUsers
		err := fetcher.FetchRecordList(appCfg, &officeUsers, []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUsers(nil), officeUsers)
	})
}

func (suite *FetchServiceSuite) TestFetchRecordCount() {
	fakeCount := func(appCfg appconfig.AppConfig, model interface{}) (int, error) {
		return 5, nil
	}

	builder := &testListQueryBuilder{
		fakeCount: fakeCount,
	}
	fetcher := NewListFetcher(builder)

	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	var officeUsers models.OfficeUsers
	count, err := fetcher.FetchRecordCount(appCfg, &officeUsers, []services.QueryFilter{})
	suite.NoError(err)
	suite.Equal(5, count)
}
