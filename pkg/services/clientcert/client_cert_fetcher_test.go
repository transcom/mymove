package clientcert

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testClientCertQueryBuilder struct {
	fakeFetchOne  func(appConfig appcontext.AppContext, model interface{}) error
	fakeCreateOne func(appConfig appcontext.AppContext, models interface{}) (*validate.Errors, error)
	fakeUpdateOne func(appConfig appcontext.AppContext, models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testClientCertQueryBuilder) FetchOne(appConfig appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(appConfig, model)
	return m
}

func (t *testClientCertQueryBuilder) CreateOne(appConfig appcontext.AppContext, model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(appConfig, model)
}

func (t *testClientCertQueryBuilder) UpdateOne(appConfig appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
	return nil, nil
}

func (t *testClientCertQueryBuilder) DeleteOne(appConfig appcontext.AppContext, model interface{}) error {
	return nil
}

func (suite *ClientCertServiceSuite) TestFetchClientCert() {
	suite.Run("if the user is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(appConfig appcontext.AppContext, model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testClientCertQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		fetcher := NewClientCertFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		adminUser, err := fetcher.FetchClientCert(suite.AppContextForTest(), filters)

		suite.NoError(err)
		suite.Equal(id, adminUser.ID)
	})

	suite.Run("if there is an error, we get it with zero admin user", func() {
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testClientCertQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewClientCertFetcher(builder)

		adminUser, err := fetcher.FetchClientCert(suite.AppContextForTest(), []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.ClientCert{}, adminUser)
	})
}
