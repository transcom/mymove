package adminuser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testAdminUserQueryBuilder struct {
	fakeFetchOne  func(appConfig appconfig.AppConfig, model interface{}) error
	fakeCreateOne func(appConfig appconfig.AppConfig, models interface{}) (*validate.Errors, error)
	fakeUpdateOne func(appConfig appconfig.AppConfig, models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testAdminUserQueryBuilder) FetchOne(appConfig appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(appConfig, model)
	return m
}

func (t *testAdminUserQueryBuilder) CreateOne(appConfig appconfig.AppConfig, model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(appConfig, model)
}

func (t *testAdminUserQueryBuilder) UpdateOne(appConfig appconfig.AppConfig, model interface{}, eTag *string) (*validate.Errors, error) {
	return nil, nil
}

func (suite *AdminUserServiceSuite) TestFetchAdminUser() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(appConfig appconfig.AppConfig, model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		appConfig := appconfig.NewAppConfig(suite.DB(), suite.logger)

		fetcher := NewAdminUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		adminUser, err := fetcher.FetchAdminUser(appConfig, filters)

		suite.NoError(err)
		suite.Equal(id, adminUser.ID)
	})

	suite.T().Run("if there is an error, we get it with zero admin user", func(t *testing.T) {
		fakeFetchOne := func(appCfg appconfig.AppConfig, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewAdminUserFetcher(builder)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		adminUser, err := fetcher.FetchAdminUser(appCfg, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.AdminUser{}, adminUser)
	})
}
