package mtoshipment

import (
	"testing"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"

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
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

	// Happy path
	suite.T().Run("If the shipment is created successfully it should be returned", func(t *testing.T) {

		mtoShipment.PickupAddressID = nil
		mtoShipment.PickupAddress.ID = uuid.Nil
		mtoShipment.DestinationAddressID = nil
		mtoShipment.DestinationAddress.ID = uuid.Nil
		mtoShipment.ID = uuid.Nil

		serviceItemsList := models.MTOServiceItems{}
		builder := query.NewQueryBuilder(suite.DB())
		createNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
			return builder
		}

		fetcher := fetch.NewFetcher(builder)
		mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(builder)

		creator := mtoShipmentCreator{
			suite.DB(),
			builder,
			fetcher,
			createNewBuilder,
			mtoServiceItemCreator,
		}

		createdShipment, err := creator.CreateMTOShipment(&mtoShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
	})

	suite.T().Run("If the shipment has mto service items", func(t *testing.T) {

		mtoShipment.PickupAddressID = nil
		mtoShipment.PickupAddress.ID = uuid.Nil
		mtoShipment.DestinationAddressID = nil
		mtoShipment.DestinationAddress.ID = uuid.Nil
		mtoShipment.ID = uuid.Nil

		builder := query.NewQueryBuilder(suite.DB())
		createNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
			return builder
		}

		fetcher := fetch.NewFetcher(builder)
		mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(builder)

		creator := mtoShipmentCreator{
			suite.DB(),
			builder,
			fetcher,
			createNewBuilder,
			mtoServiceItemCreator,
		}

		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeCS,
				Name: "ReServiceCodeCS",
			},
		})

		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDCRT,
				Name: "ReServiceCodeDCRT",
			},
		})

		serviceItemsList := []models.MTOServiceItem{
			{
				MoveTaskOrderID: mtoShipment.MoveTaskOrder.ID,
				MoveTaskOrder:   mtoShipment.MoveTaskOrder,
				Status:          models.MTOServiceItemStatusSubmitted,
				ReService: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
			{
				MoveTaskOrderID: mtoShipment.MoveTaskOrder.ID,
				MoveTaskOrder:   mtoShipment.MoveTaskOrder,
				Status:          models.MTOServiceItemStatusSubmitted,
				ReService: models.ReService{
					Code: models.ReServiceCodeDCRT,
				},
			},
		}

		createdShipment, err := creator.CreateMTOShipment(&mtoShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
	})
}
