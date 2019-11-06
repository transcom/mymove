package serviceitem

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceItemServiceSuite) TestCreateServiceItem() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	userInfo := models.ServiceItem{
		MoveTaskOrderID: moveTaskOrder.ID,
		MoveTaskOrder:   moveTaskOrder,
	}

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(moveTaskOrder.ID))
			return nil
		}
		fakeCreateOne := func(interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", moveTaskOrder.ID)}

		builder := &testServiceItemQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewServiceItemCreator(builder)
		_, verrs, err := creator.CreateServiceItem(&userInfo, filter)
		suite.NoError(err)
		suite.Nil(verrs)

	})

	// Bad move task order ID
	suite.T().Run("If we are provided a move task order that doesn't exist, the create should fail", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}
		filter := []services.QueryFilter{query.NewQueryFilter("id", "=", "b9c41d03-c730-4580-bd37-9ccf4845af6c")}
		builder := &testServiceItemQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		creator := NewServiceItemCreator(builder)
		_, _, err := creator.CreateServiceItem(&userInfo, filter)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
