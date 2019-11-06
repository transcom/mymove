package office

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

type testOfficeListQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testOfficeListQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(model)
	return m
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

func (suite *OfficeServiceSuite) TestFetchOfficeList() {
	suite.T().Run("if the transportation office is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.TransportationOffice{ID: id})))
			return nil
		}
		builder := &testOfficeListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOfficeListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		offices, err := fetcher.FetchOfficeList(filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, offices[0].ID)
	})

	suite.T().Run("if there is an error, we get it with no offices", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewOfficeListFetcher(builder)

		offices, err := fetcher.FetchOfficeList([]services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationOffices(nil), offices)
	})
}
