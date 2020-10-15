package mtoshipment

import (
	"testing"

	"github.com/transcom/mymove/pkg/unit"

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
	mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	builder := query.NewQueryBuilder(suite.DB())
	createNewBuilder := func(db *pop.Connection) createMTOShipmentQueryBuilder {
		return builder
	}
	mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(builder)
	fetcher := fetch.NewFetcher(builder)
	creator := mtoShipmentCreator{
		suite.DB(),
		builder,
		fetcher,
		createNewBuilder,
		mtoServiceItemCreator,
	}

	// Invalid ID fields set
	suite.T().Run("invalid IDs found", func(t *testing.T) {
		// The default shipment created will have IDs filled in for subobjects, but let's make sure one is set anyway:
		moveTaskOrderID := mtoShipment.MoveTaskOrderID
		if mtoShipment.PickupAddress != nil && mtoShipment.PickupAddress.ID != uuid.Nil {
			mtoShipment.PickupAddress.ID = moveTaskOrderID
		}
		createdShipment, err := creator.CreateMTOShipment(&mtoShipment, nil)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		invalidErr := err.(services.InvalidInputError)
		suite.NotNil(invalidErr.ValidationErrors)
		suite.NotEmpty(invalidErr.ValidationErrors)
	})

	// Happy path
	suite.T().Run("If the shipment is created successfully it should be returned", func(t *testing.T) {
		mtoShipment := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(mtoShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
	})

	suite.T().Run("If the shipment is created successfully with submitted status it should be returned", func(t *testing.T) {
		mtoShipment := clearShipmentIDFields(&mtoShipment)
		mtoShipment.Status = models.MTOShipmentStatusSubmitted
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(mtoShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
	})

	suite.T().Run("If the shipment has mto service items", func(t *testing.T) {
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
				Name: "ReServiceCodeDDSHUT",
			},
		})

		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
				Name: "ReServiceCodeDOFSIT",
			},
		})

		serviceItemsList := []models.MTOServiceItem{
			{
				MoveTaskOrderID: mtoShipment.MoveTaskOrder.ID,
				MoveTaskOrder:   mtoShipment.MoveTaskOrder,
				ReService: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				MoveTaskOrderID: mtoShipment.MoveTaskOrder.ID,
				MoveTaskOrder:   mtoShipment.MoveTaskOrder,
				ReService: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}

		mtoShipment := clearShipmentIDFields(&mtoShipment)
		weight := unit.Pound(2)
		mtoShipment.PrimeEstimatedWeight = &weight // for DDSHUT service item type
		createdShipment, err := creator.CreateMTOShipment(mtoShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.NotNil(createdShipment.MTOServiceItems, "Service Items are empty")
		suite.Equal(createdShipment.MTOServiceItems[0].MTOShipmentID, &createdShipment.ID, "Service items are not the same")
	})

	suite.T().Run("If the move already has a submitted NTS shipment, it should return a validation error", func(t *testing.T) {
		ntsShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: ntsShipment.MoveTaskOrderID,
				ShipmentType:    models.MTOShipmentTypeHHGIntoNTSDom,
				Status:          models.MTOShipmentStatusDraft,
			},
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSShipment := clearShipmentIDFields(&secondNTSShipment)
		createdShipment, err := creator.CreateMTOShipment(cleanedNTSShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})
}

// Clears all the ID fields that we need to be null for a new shipment to get created:
func clearShipmentIDFields(shipment *models.MTOShipment) *models.MTOShipment {
	shipment.PickupAddressID = nil
	shipment.PickupAddress.ID = uuid.Nil
	shipment.DestinationAddressID = nil
	shipment.DestinationAddress.ID = uuid.Nil
	shipment.SecondaryPickupAddressID = nil
	shipment.SecondaryPickupAddress = nil
	shipment.SecondaryDeliveryAddressID = nil
	shipment.SecondaryDeliveryAddress = nil
	shipment.ID = uuid.Nil
	if len(shipment.MTOAgents) > 0 {
		for _, agent := range shipment.MTOAgents {
			agent.ID = uuid.Nil
			agent.MTOShipmentID = uuid.Nil
		}
	}

	return shipment
}
