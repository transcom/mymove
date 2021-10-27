package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type createShipmentSubtestData struct {
	appCtx          appcontext.AppContext
	move            models.Move
	shipmentCreator mtoShipmentCreator
}

func (suite *MTOShipmentServiceSuite) createSubtestData(assertions testdatagen.Assertions) (subtestData *createShipmentSubtestData) {
	subtestData = &createShipmentSubtestData{}

	subtestData.move = testdatagen.MakeMove(suite.DB(), assertions)

	subtestData.appCtx = suite.TestAppContext()

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
	suite.Run("invalid IDs found", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})

		createdShipment, err := creator.CreateMTOShipment(appCtx, &mtoShipment, nil)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		invalidErr := err.(apperror.InvalidInputError)
		suite.NotNil(invalidErr.ValidationErrors)
		suite.NotEmpty(invalidErr.ValidationErrors)
	})

	// Unhappy path
	suite.Run("When required requested pickup dates are zero (required for NTS & HHG shipment types)", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		// default is HHG
		hhgShipmentFail := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})
		hhgShipmentFailClear := clearShipmentIDFields(&hhgShipmentFail)
		hhgShipmentFailClear.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(appCtx, hhgShipmentFailClear, nil)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "RequestedPickupDate")
	})

	suite.Run("When non-required requested pickup dates are zero (not required for NTSr shipment type)", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		ntsrShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
			Stub: true,
		})
		ntsrShipmentNoIDs := clearShipmentIDFields(&ntsrShipment)
		ntsrShipmentNoIDs.RequestedPickupDate = new(time.Time)

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(appCtx, ntsrShipmentNoIDs, nil)

		suite.NoError(err)
	})

	// Happy path
	suite.Run("If the shipment is created successfully it should be returned", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
	})

	suite.Run("If the shipment is created successfully with submitted status it should be returned", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
	})

	suite.Run("If the shipment has mto service items", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeDDSHUT,
			models.ReServiceCodeDOFSIT,
		}

		var reServiceCode models.ReService
		if err := appCtx.DB().Where("code = $1", expectedReServiceCodes[0]).First(&reServiceCode); err != nil {
			// Something is truncating these when all server tests run, but we need some values for reServices
			for _, serviceCode := range expectedReServiceCodes {
				testdatagen.MakeReService(appCtx.DB(), testdatagen.Assertions{
					ReService: models.ReService{
						Code:      serviceCode,
						Name:      "test",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				})
			}
		}

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
		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: &weight,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.NotNil(createdShipment.MTOServiceItems, "Service Items are empty")
		suite.Equal(createdShipment.MTOServiceItems[0].MTOShipmentID, &createdShipment.ID, "Service items are not the same")
	})

	suite.Run("If the move already has a submitted NTS shipment, it should return a validation error", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				Status:       models.MTOShipmentStatusDraft,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSShipment := clearShipmentIDFields(&secondNTSShipment)
		createdShipment, err := creator.CreateMTOShipment(appCtx, cleanedNTSShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("If the move already has a submitted NTSr shipment, it should return a validation error", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		})

		secondNTSrShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:       models.MTOShipmentStatusDraft,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		cleanedNTSrShipment := clearShipmentIDFields(&secondNTSrShipment)
		createdShipment, err := creator.CreateMTOShipment(appCtx, cleanedNTSrShipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("422 Validation Error - only one mto agent of each type", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
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

		shipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				MTOAgents: agents,
			},
			Stub: true,
		})

		serviceItemsList := models.MTOServiceItems{}
		createdShipment, err := creator.CreateMTOShipment(appCtx, &shipment, serviceItemsList)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Move status transitions when a new shipment is created and SUBMITTED", func() {
		// If a new shipment is added to an APPROVED move and given the SUBMITTED status,
		// the move should transition to "APPROVALS REQUESTED"
		subtestData := suite.createSubtestData(testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator
		move := subtestData.move

		shipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		})
		cleanShipment := clearShipmentIDFields(&shipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(appCtx, cleanShipment, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
		suite.Equal(move.ID.String(), createdShipment.MoveTaskOrderID.String())

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	suite.Run("Sets SIT days allowance to default", func() {
		// This test will have to change in the future, but for now, service members are expected to get 90 days by
		// default.
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"INTERNATIONAL_HHG", models.MTOShipmentTypeInternationalHHG},
			{"INTERNATIONAL_UB", models.MTOShipmentTypeInternationalUB},
			{"HHG_LONGHAUL_DOMESTIC", models.MTOShipmentTypeHHGLongHaulDom},
			{"HHG_SHORTHAUL_DOMESTIC", models.MTOShipmentTypeHHGShortHaulDom},
			{"HHG_INTO_NTS_DOMESTIC", models.MTOShipmentTypeHHGIntoNTSDom},
			{"HHG_OUTOF_NTS_DOMESTIC", models.MTOShipmentTypeHHGOutOfNTSDom},
			{"MOTORHOME", models.MTOShipmentTypeMotorhome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
		}

		for _, tt := range testCases {
			tt := tt
			suite.Run(tt.desc, func() {
				mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
					Move: subtestData.move,
					MTOShipment: models.MTOShipment{
						ShipmentType: tt.shipmentType,
					},
					Stub: true,
				})

				clearedShipment := clearShipmentIDFields(&mtoShipment)

				createdShipment, err := creator.CreateMTOShipment(appCtx, clearedShipment, nil)

				suite.NoError(err)

				suite.Equal(models.DefaultServiceMemberSITDaysAllowance, *createdShipment.SITDaysAllowance)
			})
		}
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
