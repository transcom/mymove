package mtoshipment

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testMTOShipmentQueryBuilder struct {
	fakeCreateOne   func(model interface{}) (*validate.Errors, error)
	fakeFetchOne    func(model interface{}, filters []services.QueryFilter) error
	fakeTransaction func(func(tx *pop.Connection) error) error
}

func (t *testMTOShipmentQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(model)
}

func (t *testMTOShipmentQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(model, filters)
}

func (t *testMTOShipmentQueryBuilder) Transaction(fn func(tx *pop.Connection) error) error {
	return t.fakeTransaction(fn)
}

func (suite *MTOShipmentServiceSuite) TestCreateMTOShipmentRequest() {
	moveTaskOrder := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := models.MTOShipment{
		MoveTaskOrderID: moveTaskOrder.ID,
	}

	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

	// Happy path
	suite.T().Run("If the user is created successfully it should be returned", func(t *testing.T) {
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}
		fakeFetchOne := func(model interface{}, filters []services.QueryFilter) error {
			return nil
		}
		fakeTx := func(fn func(tx *pop.Connection) error) error {
			return fn(&pop.Connection{})
		}

		builder := &testMTOShipmentQueryBuilder{
			fakeCreateOne:   fakeCreateOne,
			fakeFetchOne:    fakeFetchOne,
			fakeTransaction: fakeTx,
		}

		fakeCreateNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
			return builder
		}

		creator := mtoShipmentCreator{
			builder:          builder,
			createNewBuilder: fakeCreateNewBuilder,
		}
		createdShipment, err := creator.CreateMTOShipment(&mtoShipment, eTag)

		suite.NoError(err)
		suite.NotNil(createdShipment)
	})

}