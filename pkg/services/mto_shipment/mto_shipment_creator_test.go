package mtoshipment

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

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

	// Unhappy path
	suite.T().Run("When required requested pickup dates are zero (required for NTS & HHG shipment types)", func(t *testing.T) {
		hhgShipmentFail := testdatagen.MakeDefaultMTOShipment(suite.DB()) // default is HHG
		hhgShipmentFailClear := clearShipmentIDFields(&hhgShipmentFail)
		hhgShipmentFailClear.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(hhgShipmentFailClear, nil)

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "RequestedPickupDate")
	})

	suite.T().Run("When non-required requested pickup dates are zero (not required for NTSr shipment type)", func(t *testing.T) {
		ntsrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
		})
		ntsrShipmentNoIDs := clearShipmentIDFields(&ntsrShipment)
		ntsrShipmentNoIDs.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(ntsrShipmentNoIDs, nil)

		suite.NoError(err)
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

	suite.T().Run("If the move already has a submitted NTSr shipment, it should return a validation error", func(t *testing.T) {
		ntsrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: ntsrShipment.MoveTaskOrderID,
				ShipmentType:    models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:          models.MTOShipmentStatusDraft,
			},
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSrShipment := clearShipmentIDFields(&secondNTSrShipment)
		createdShipment, err := creator.CreateMTOShipment(cleanedNTSrShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})
}

// Clears all the ID fields that we need to be null for a new shipment to get created:
func clearShipmentIDFields(shipment *models.MTOShipment) *models.MTOShipment {
	if shipment.PickupAddress != nil {
		shipment.PickupAddressID = nil
		shipment.PickupAddress.ID = uuid.Nil
	}
	if shipment.DestinationAddress != nil {
		shipment.DestinationAddressID = nil
		shipment.DestinationAddress.ID = uuid.Nil
	}
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
