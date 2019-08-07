package office

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
}

func (t *testOfficeQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *OfficeServiceSuite) TestFetchOffice() {
	suite.T().Run("if the transportation office is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testOfficeQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		office, err := fetcher.FetchOffice(filters)

		suite.NoError(err)
		suite.Equal(id, office.ID)
	})

	suite.T().Run("if there is an error, we get it with zero office", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeFetcher(builder)

		office, err := fetcher.FetchOffice([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.TransportationOffice{}, office)
	})
}
