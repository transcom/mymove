package officeuser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeUserQueryBuilder struct {
	fakeFetchOne             func(model interface{}) error
	fakeCreateOne            func(models interface{}) (*validate.Errors, error)
	fakeUpdateOne            func(models interface{}, eTag *string) (*validate.Errors, error)
	fakeQueryForAssociations func(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error
}

func (t *testOfficeUserQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (t *testOfficeUserQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return nil, nil
}

func (t *testOfficeUserQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return nil, nil
}

func (t *testOfficeUserQueryBuilder) QueryForAssociations(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
	return nil
}

func (suite *OfficeUserServiceSuite) TestFetchOfficeUser() {
	suite.T().Run("if the user is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		fakeCreateOne := func(interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		fetcher := NewOfficeUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		officeUser, err := fetcher.FetchOfficeUser(filters)

		suite.NoError(err)
		suite.Equal(id, officeUser.ID)
	})

	suite.T().Run("if there is an error, we get it with zero office user", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeUserFetcher(builder)

		officeUser, err := fetcher.FetchOfficeUser([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUser{}, officeUser)
	})
}

func (suite *OfficeUserServiceSuite) TestFetchOfficeUserPop() {
	suite.T().Run("returns office user on success", func(t *testing.T) {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := NewOfficeUserFetcherPop(suite.DB())

		fetchedUser, err := fetcher.FetchOfficeUserByID(officeUser.ID)

		suite.NoError(err)
		suite.Equal(officeUser.ID, fetchedUser.ID)
	})

	suite.T().Run("returns zero value office user on error", func(t *testing.T) {
		fetcher := NewOfficeUserFetcherPop(suite.DB())
		officeUser, err := fetcher.FetchOfficeUserByID(uuid.Nil)

		suite.Error(err)
		suite.Equal(err.Error(), "sql: no rows in result set")
		suite.Equal(uuid.Nil, officeUser.ID)
	})
}
