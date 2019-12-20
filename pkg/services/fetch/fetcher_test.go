package fetch

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testFetcherQueryBuilder struct {
	fakeFetch func(model interface{}, filters []services.QueryFilter) error
}

func (t *testFetcherQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetch(model, filters)
	return m
}

func (suite *FetchServiceSuite) TestFetchRecord() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetch := func(model interface{}, filters []services.QueryFilter) error {
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

		ou, err := fetcher.FetchRecord(&models.OfficeUser{}, filters)
		officeUser := ou.(*models.OfficeUser)

		suite.NoError(err)
		suite.Equal(id, officeUser.ID)
	})

	suite.T().Run("if there is an error, we get it with no office user", func(t *testing.T) {
		fakeFetch := func(model interface{}, filters []services.QueryFilter) error {
			return errors.New("Fetch error")
		}
		builder := &testFetcherQueryBuilder{
			fakeFetch: fakeFetch,
		}

		fetcher := NewFetcher(builder)

		ou, err := fetcher.FetchRecord(&models.OfficeUser{}, []services.QueryFilter{})
		officeUser := ou.(*models.OfficeUser)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUser{}, *officeUser)
	})
}
