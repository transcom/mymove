package fetch

import (
	"errors"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testFetcherQueryBuilder struct {
	fakeFetch func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
}

func (t *testFetcherQueryBuilder) FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetch(appCtx, model, filters)
	return m
}

func (suite *FetchServiceSuite) TestFetchRecord() {
	suite.Run("if the user is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetch := func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.ValueOf(models.OfficeUser{ID: id}))
			return nil
		}
		builder := &testFetcherQueryBuilder{
			fakeFetch: fakeFetch,
		}

		fetcher := NewFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		officeUser := &models.OfficeUser{}
		err = fetcher.FetchRecord(suite.AppContextForTest(), officeUser, filters)

		suite.NoError(err)
		suite.Equal(id, officeUser.ID)
	})

	suite.Run("if there is an error, we get it with no office user", func() {
		fakeFetch := func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
			return errors.New("Fetch error")
		}
		builder := &testFetcherQueryBuilder{
			fakeFetch: fakeFetch,
		}

		fetcher := NewFetcher(builder)

		officeUser := &models.OfficeUser{}
		err := fetcher.FetchRecord(suite.AppContextForTest(), officeUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Resource not found: Fetch error")
		suite.Equal(models.OfficeUser{}, *officeUser)
	})

	suite.Run("reflection error", func() {
		fakeFetch := func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
			return errors.New("Fetch error")
		}
		builder := &testFetcherQueryBuilder{
			fakeFetch: fakeFetch,
		}

		fetcher := NewFetcher(builder)

		officeUser := models.OfficeUser{}
		err := fetcher.FetchRecord(suite.AppContextForTest(), officeUser, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), query.FetchOneReflectionMessage)

		err = fetcher.FetchRecord(suite.AppContextForTest(), 1, []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), query.FetchOneReflectionMessage)
	})
}
