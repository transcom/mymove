package electronicorder

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testElectronicOrderListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
	fakeCount     func(model interface{}) (int, error)
}

func (t *testElectronicOrderListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(model)
	return m
}

func (t *testElectronicOrderListQueryBuilder) Count(model interface{}, filters []services.QueryFilter) (int, error) {
	count, m := t.fakeCount(model)
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

func (suite *ElectronicOrderServiceSuite) TestFetchElectronicOrderList() {
	suite.T().Run("if the transportation order is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.ElectronicOrder{ID: id})))
			return nil
		}
		builder := &testElectronicOrderListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewElectronicOrderListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		electronicOrders, err := fetcher.FetchElectronicOrderList(filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, electronicOrders[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no electronic orders", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testElectronicOrderListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewElectronicOrderListFetcher(builder)

		electronicOrders, err := fetcher.FetchElectronicOrderList([]services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.ElectronicOrders(nil), electronicOrders)
	})
}
