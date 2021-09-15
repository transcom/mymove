package mtoshipment

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type createShipmentSubtestData struct {
	appContext      appcontext.AppContext
	move            models.Move
	shipmentCreator mtoShipmentCreator
}

func (suite *MTOShipmentServiceSuite) createSubtestData(assertions testdatagen.Assertions) (subtestData *createShipmentSubtestData) {
	subtestData = &createShipmentSubtestData{}

	subtestData.appContext = suite.TestAppContext()

	subtestData.move = testdatagen.MakeMove(suite.DB(), assertions)

	builder := query.NewQueryBuilder()
	createNewBuilder := func() createMTOShipmentQueryBuilder {
		return builder
	}
	moveRouter := moverouter.NewMoveRouter()
	fetcher := fetch.NewFetcher(builder)

	subtestData.shipmentCreator = mtoShipmentCreator{
		builder,
		fetcher,
		createNewBuilder,
		moveRouter,
	}

	return subtestData
}

func (suite *MTOShipmentServiceSuite) TestCreateMTOShipment() {
	// Invalid ID fields set
	suite.T().Run("invalid IDs found", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})

		createdShipment, err := creator.CreateMTOShipment(subtestData.appContext, &mtoShipment, nil)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		invalidErr := err.(services.InvalidInputError)
		suite.NotNil(invalidErr.ValidationErrors)
		suite.NotEmpty(invalidErr.ValidationErrors)
	})

	// Unhappy path
	suite.T().Run("When required requested pickup dates are zero (required for NTS & HHG shipment types)", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		// default is HHG
		hhgShipmentFail := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})
		hhgShipmentFailClear := clearShipmentIDFields(&hhgShipmentFail)
		hhgShipmentFailClear.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(suite.TestAppContext(), hhgShipmentFailClear, nil)

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "RequestedPickupDate")
	})

	suite.T().Run("When non-required requested pickup dates are zero (not required for NTSr shipment type)", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		ntsrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
			Stub: true,
		})
		ntsrShipmentNoIDs := clearShipmentIDFields(&ntsrShipment)
		ntsrShipmentNoIDs.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(suite.TestAppContext(), ntsrShipmentNoIDs, nil)

		suite.NoError(err)
	})

	// Happy path
	suite.T().Run("If the shipment is created successfully it should be returned", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
	})

	suite.T().Run("If the shipment is created successfully with submitted status it should be returned", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
	})

	suite.T().Run("If the shipment has mto service items", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

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
				MoveTaskOrderID: subtestData.move.ID,
				MoveTaskOrder:   subtestData.move,
				ReService: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
			{
				MoveTaskOrderID: subtestData.move.ID,
				MoveTaskOrder:   subtestData.move,
				ReService: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}

		weight := unit.Pound(2) // for DDSHUT service item type
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: &weight,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.NotNil(createdShipment.MTOServiceItems, "Service Items are empty")
		suite.Equal(createdShipment.MTOServiceItems[0].MTOShipmentID, &createdShipment.ID, "Service items are not the same")
	})

	suite.T().Run("If the move already has a submitted NTS shipment, it should return a validation error", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				Status:       models.MTOShipmentStatusDraft,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSShipment := clearShipmentIDFields(&secondNTSShipment)
		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), cleanedNTSShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("If the move already has a submitted NTSr shipment, it should return a validation error", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:       models.MTOShipmentStatusDraft,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSrShipment := clearShipmentIDFields(&secondNTSrShipment)
		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), cleanedNTSrShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("422 Validation Error - only one mto agent of each type", func(t *testing.T) {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		creator := subtestData.shipmentCreator

		firstName := "First"
		lastName := "Last"
		email := "test@gmail.com"

		var agents models.MTOAgents

		agent1 := models.MTOAgent{
			FirstName:    &firstName,
			LastName:     &lastName,
			Email:        &email,
			MTOAgentType: models.MTOAgentReceiving,
		}

		agent2 := models.MTOAgent{
			FirstName:    &firstName,
			LastName:     &lastName,
			Email:        &email,
			MTOAgentType: models.MTOAgentReceiving,
		}

		agents = append(agents, agent1, agent2)

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				MTOAgents: agents,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), &shipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("Move status transitions when a new shipment is created and SUBMITTED", func(t *testing.T) {
		// If a new shipment is added to an APPROVED move and given the SUBMITTED status,
		// the move should transition to "APPROVALS REQUESTED"
		subtestData := suite.createSubtestData(testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		})
		creator := subtestData.shipmentCreator
		move := subtestData.move

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		})
		cleanShipment := clearShipmentIDFields(&shipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.TestAppContext(), cleanShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
		suite.Equal(move.ID.String(), createdShipment.MoveTaskOrderID.String())

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
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
	if shipment.SecondaryPickupAddress != nil {
		shipment.SecondaryPickupAddressID = nil
		shipment.SecondaryPickupAddress.ID = uuid.Nil
	}

	if shipment.SecondaryDeliveryAddress != nil {
		shipment.SecondaryDeliveryAddressID = nil
		shipment.SecondaryDeliveryAddress.ID = uuid.Nil
	}

	shipment.ID = uuid.Nil
	if len(shipment.MTOAgents) > 0 {
		for _, agent := range shipment.MTOAgents {
			agent.ID = uuid.Nil
			agent.MTOShipmentID = uuid.Nil
		}
	}

	return shipment
}
