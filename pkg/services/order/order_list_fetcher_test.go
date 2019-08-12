package order

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOrderListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testOrderListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchMany(model)
	return m
}

func (suite *OrderServiceSuite) TestFetchOrderList() {
	suite.T().Run("if the transportation order is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.Order{ID: id})))
			return nil
		}
		builder := &testOrderListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOrderListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		orders, err := fetcher.FetchOrderList(filters)

		suite.NoError(err)
		suite.Equal(id, orders[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no orders", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOrderListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOrderListFetcher(builder)

		orders, err := fetcher.FetchOrderList([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.Orders(nil), orders)
	})
}
