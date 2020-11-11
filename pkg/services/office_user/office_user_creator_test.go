package officeuser

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OfficeUserServiceSuite) TestCreateOfficeUser() {
	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())
	userInfo := models.OfficeUser{
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(transportationOffice.ID))
			return nil
		}
		fakeCreateOne := func(interface{}) (*validate.Errors, error) {
			return nil, nil
		}
		fakeQueryAssociations := func(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
			return nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", transportationOffice.ID)}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:             fakeFetchOne,
			fakeCreateOne:            fakeCreateOne,
			fakeQueryForAssociations: fakeQueryAssociations,
		}

		creator := NewOfficeUserCreator(suite.DB(), builder)
		_, verrs, err := creator.CreateOfficeUser(&userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)

	})

	// Bad transportation office ID
	suite.T().Run("If we are provided a transportation office that doesn't exist, the create should fail", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewOfficeUserCreator(suite.DB(), builder)
		_, _, err := creator.CreateOfficeUser(&userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
