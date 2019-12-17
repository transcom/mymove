package mtoserviceitem

import (
	"errors"
	"testing"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testMTOServiceItemQueryBuilder struct {
	fakeCreateOne func(model interface{}) (*validate.Errors, error)
}

func (t *testMTOServiceItemQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(model)
}

func (suite *MTOServiceItemServiceSuite) TestCreateMTOServiceItem() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: moveTaskOrder.ID,
		MoveTaskOrder:   moveTaskOrder,
	}

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		builder := &testMTOServiceItemQueryBuilder{
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewMTOServiceItemCreator(builder)
		createdServiceItem, verrs, err := creator.CreateMTOServiceItem(&serviceItem)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItem)
	})

	// Bad data which could be IDs that doesn't exist (MoveTaskOrderID or REServiceID)
	suite.T().Run("If error when trying to create, the create should fail", func(t *testing.T) {
		expectedError := "Can't create service item for some reason"
		verrs := validate.NewErrors()
		verrs.Add("test", expectedError)
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return verrs, errors.New(expectedError)
		}
		builder := &testMTOServiceItemQueryBuilder{
			fakeCreateOne: fakeCreateOne,
		}

		creator := NewMTOServiceItemCreator(builder)
		createdServiceItem, verrs, err := creator.CreateMTOServiceItem(&serviceItem)
		suite.Error(err)
		suite.Error(verrs)
		suite.Nil(createdServiceItem)
		suite.Equal(verrs, verrs)
		suite.Equal(expectedError, err.Error())
	})
}
