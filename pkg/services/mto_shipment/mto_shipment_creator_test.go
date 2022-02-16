package mtoshipment

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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

	subtestData.appCtx = suite.AppContextForTest()

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

	suite.Run("Test requested pickup date requirement for various shipment types", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		// Default is HHG, but we set it explicitly below via the test cases
		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			Stub: true,
		})

		testCases := []struct {
			input        *time.Time
			shipmentType models.MTOShipmentType
			shouldError  bool
		}{
			{nil, models.MTOShipmentTypeHHG, true},
			{&time.Time{}, models.MTOShipmentTypeHHG, true},
			{swag.Time(time.Now()), models.MTOShipmentTypeHHG, false},
			{nil, models.MTOShipmentTypeHHGOutOfNTSDom, false},
			{&time.Time{}, models.MTOShipmentTypeHHGOutOfNTSDom, true},
			{swag.Time(time.Now()), models.MTOShipmentTypeHHGOutOfNTSDom, true},
			{nil, models.MTOShipmentTypePPM, false},
			{swag.Time(time.Now()), models.MTOShipmentTypePPM, false},
		}

		for _, testCase := range testCases {
			mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
			mtoShipmentClear.ShipmentType = testCase.shipmentType
			mtoShipmentClear.RequestedPickupDate = testCase.input
			_, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, nil)

			if testCase.shouldError {
				if suite.Errorf(err, "should have errored for a %s shipment with requested pickup date set to %s", testCase.shipmentType, testCase.input) {
					suite.IsType(apperror.InvalidInputError{}, err)
					suite.Contains(err.Error(), "RequestedPickupDate")
				}
			} else {
				suite.NoErrorf(err, "should have not errored for a %s shipment with requested pickup date set to %s", testCase.shipmentType, testCase.input)
			}
		}
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
	suite.Run("If the shipment is created successfully with a destination address type it should be returned", func() {
		destinationType := models.DestinationTypeHomeOfRecord
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move:        subtestData.move,
			MTOShipment: models.MTOShipment{DestinationType: &destinationType},
			Stub:        true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		serviceItemsList := models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, serviceItemsList)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
		suite.Equal(string(models.DestinationTypeHomeOfRecord), string(*createdShipment.DestinationType))
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

	suite.Run("If the submitted shipment has a storage facility attached", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{
			Stub: true,
		})

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				StorageFacility: &storageFacility,
				ShipmentType:    models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:          models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, nil)
		suite.NoError(err)
		suite.NotNil(createdShipment.StorageFacility)
		suite.Equal(storageFacility.Address.StreetAddress1, createdShipment.StorageFacility.Address.StreetAddress1)
	})

	suite.Run("If the submitted shipment is an NTS shipment", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		ntsRecordedWeight := unit.Pound(980)
		requestedDeliveryDate := time.Date(testdatagen.GHCTestYear, time.April, 5, 0, 0, 0, 0, time.UTC)
		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:                models.MTOShipmentStatusSubmitted,
				NTSRecordedWeight:     &ntsRecordedWeight,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, nil)
		if suite.NoError(err) {
			if suite.NotNil(createdShipment.NTSRecordedWeight) {
				suite.Equal(ntsRecordedWeight, *createdShipment.NTSRecordedWeight)
			}
			if suite.NotNil(createdShipment.RequestedDeliveryDate) {
				suite.Equal(requestedDeliveryDate, *createdShipment.RequestedDeliveryDate)
			}
		}
	})

	suite.Run("If the submitted shipment is a PPM shipment", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusDraft,
			},
			Stub: true,
		})

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear, nil)

		suite.NoError(err)
		suite.NotNil(createdShipment)
	})

	suite.Run("When NTSRecordedWeight it set for a non NTS Release shipment", func() {
		subtestData := suite.createSubtestData(testdatagen.Assertions{})
		appCtx := subtestData.appCtx
		creator := subtestData.shipmentCreator

		ntsRecordedWeight := unit.Pound(980)
		mtoShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: subtestData.move,
			MTOShipment: models.MTOShipment{
				Status:            models.MTOShipmentStatusSubmitted,
				NTSRecordedWeight: &ntsRecordedWeight,
			},
			Stub: true,
		})
		ntsrShipmentNoIDs := clearShipmentIDFields(&mtoShipment)
		ntsrShipmentNoIDs.RequestedPickupDate = swag.Time(time.Now())

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(appCtx, ntsrShipmentNoIDs, nil)

		if suite.Errorf(err, "should have errored for a %s shipment with ntsRecordedWeight set", ntsrShipmentNoIDs.ShipmentType) {
			suite.IsType(apperror.InvalidInputError{}, err)
			suite.Contains(err.Error(), "NTSRecordedWeight")
		}
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
			{"PPM", models.MTOShipmentTypePPM},
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

	if shipment.StorageFacility != nil {
		shipment.StorageFacilityID = nil
		shipment.StorageFacility.ID = uuid.Nil
		shipment.StorageFacility.AddressID = uuid.Nil
		shipment.StorageFacility.Address.ID = uuid.Nil
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
