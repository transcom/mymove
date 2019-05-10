package user

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

type testOfficeUserQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
}

func (t *testOfficeUserQueryBuilder) FetchOne(model interface{}, field string, value interface{}) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *UserServiceSuite) TestFetchOfficeUser() {
	suite.T().Run("if the user it fetched, it should be returned", func(t *testing.T){
		id , err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeUserFetcher(builder)

		officeUser, err := fetcher.FetchOfficeUser("id", id)

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

		officeUser, err := fetcher.FetchOfficeUser("id", 1)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUser{}, officeUser)
	})
}
