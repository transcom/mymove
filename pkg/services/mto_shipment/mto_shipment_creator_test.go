package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type createShipmentSubtestData struct {
	move            models.Move
	shipmentCreator services.MTOShipmentCreator
}

func (suite *MTOShipmentServiceSuite) createSubtestData(customs []factory.Customization) (subtestData *createShipmentSubtestData) {
	subtestData = &createShipmentSubtestData{}

	subtestData.move = factory.BuildMove(suite.DB(), customs, nil)

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	fetcher := fetch.NewFetcher(builder)
	addressCreator := address.NewAddressCreator()

	subtestData.shipmentCreator = NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)

	return subtestData
}

// This func is for the PrimeAPI V2 subtest data tests for createMTOShipment
func (suite *MTOShipmentServiceSuite) createSubtestDataV2(customs []factory.Customization) (subtestData *createShipmentSubtestData) {
	subtestData = &createShipmentSubtestData{}

	subtestData.move = factory.BuildMove(suite.DB(), customs, nil)

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	fetcher := fetch.NewFetcher(builder)
	addressCreator := address.NewAddressCreator()

	subtestData.shipmentCreator = NewMTOShipmentCreatorV2(builder, fetcher, moveRouter, addressCreator)

	return subtestData
}

func (suite *MTOShipmentServiceSuite) TestCreateMTOShipment() {
	futureDate := models.TimePointer(time.Now().AddDate(0, 0, 3)) //adds 3 days to current date
	// Invalid ID fields set
	suite.Run("invalid IDs found", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), &mtoShipment)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		invalidErr := err.(apperror.InvalidInputError)
		suite.NotNil(invalidErr.ValidationErrors)
		suite.NotEmpty(invalidErr.ValidationErrors)
	})

	// Happy path
	suite.Run("If a domestic shipment is created successfully it should be returned", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
		// both pickup and destination addresses should be CONUS
		suite.False(*createdShipment.PickupAddress.IsOconus)
		suite.False(*createdShipment.DestinationAddress.IsOconus)
		suite.Equal(createdShipment.MarketCode, models.MarketCodeDomestic)
	})

	suite.Run("If an international shipment is created successfully it should be returned", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.Address{
					State: "AK",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					State: "HI",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
		// both pickup and destination addresses should be OCONUS since Alaska & Hawaii are considered OCONUS
		suite.True(*createdShipment.PickupAddress.IsOconus)
		suite.True(*createdShipment.DestinationAddress.IsOconus)
		suite.Equal(createdShipment.MarketCode, models.MarketCodeInternational)
	})

	suite.Run("If the shipment has an international address it should be returned", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		internationalAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Country{
					Country:     "GB",
					CountryName: "UNITED KINGDOM",
				},
			},
		}, nil)
		// stubbed countries need an ID
		internationalAddress.ID = uuid.Must(uuid.NewV4())

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model:    internationalAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.Error(err)
		suite.Equal("failed to create pickup address - the country GB is not supported at this time - only US is allowed", err.Error())
	})

	suite.Run("If the shipment is created successfully it should return ShipmentLocator", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		mtoShipment2 := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}
		mtoShipmentClear2 := clearShipmentIDFields(&mtoShipment2)
		mtoShipmentClear2.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)
		createdShipment2, err2 := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear2)

		suite.NoError(err)
		suite.NoError(err2)
		suite.NotNil(createdShipment)
		suite.NotEmpty(createdShipment.ShipmentLocator)

		// ShipmentLocator = move Locator + "-" + Shipment Seq Num
		// checks for proper structure of shipmentLocator
		suite.Equal(subtestData.move.Locator, (*createdShipment.ShipmentLocator)[0:6])
		suite.Equal("-", (*createdShipment.ShipmentLocator)[6:7])
		suite.Equal("01", (*createdShipment.ShipmentLocator)[7:9])

		// check if seq number is increased by 1
		suite.Equal("02", (*createdShipment2.ShipmentLocator)[7:9])
	})

	suite.Run("If the shipment is created successfully with a destination address type it should be returned", func() {
		destinationType := models.DestinationTypeHomeOfRecord
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator
		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{DestinationType: &destinationType, RequestedPickupDate: futureDate},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status)
		suite.NotEmpty(createdShipment.PickupAddressID)
		suite.NotEmpty(createdShipment.DestinationAddressID)
		suite.Equal(string(models.DestinationTypeHomeOfRecord), string(*createdShipment.DestinationType))
	})

	suite.Run("If the shipment has a nil destination address, the duty station address should be used if an HHG", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		// Make sure we have all parts we care about in the source address in order to assert against them further down
		suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.AddressID)
		if suite.NotNil(subtestData.move.Orders.NewDutyLocation.Address) {
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.ID)
			suite.Equal(subtestData.move.Orders.NewDutyLocation.AddressID, subtestData.move.Orders.NewDutyLocation.Address.ID)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.StreetAddress1)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.StreetAddress2)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.StreetAddress3)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.City)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.State)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.PostalCode)
			suite.NotEmpty(subtestData.move.Orders.NewDutyLocation.Address.CountryId)
		}

		testCases := []struct {
			shipmentType      models.MTOShipmentType
			expectDutyStation bool
		}{
			{models.MTOShipmentTypeHHG, true},
			{models.MTOShipmentTypeHHGIntoNTS, false},
			{models.MTOShipmentTypeHHGOutOfNTS, false},
			{models.MTOShipmentTypePPM, false},
		}
		for _, testCase := range testCases {
			mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
				{
					Model:    subtestData.move,
					LinkOnly: true,
				},
				{
					Model: models.MTOShipment{
						ShipmentType:        testCase.shipmentType,
						RequestedPickupDate: futureDate,
					},
				},
			}, nil)

			mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
			mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}
			mtoShipmentClear.DestinationAddress = nil
			mtoShipmentClear.StorageFacility = nil

			createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

			suite.NoError(err, testCase.shipmentType)
			suite.NotNil(createdShipment, testCase.shipmentType)
			suite.Equal(models.MTOShipmentStatusDraft, createdShipment.Status, testCase.shipmentType)

			if testCase.expectDutyStation {
				suite.NotEmpty(createdShipment.DestinationAddressID, testCase.shipmentType)
				// Original and new IDs should not match since we should be creating an entirely new address record
				suite.NotEqual(subtestData.move.Orders.NewDutyLocation.AddressID, createdShipment.DestinationAddressID, testCase.shipmentType)

				// Check address fields are set appropriately when destination duty station info is copied over
				if suite.NotNil(createdShipment.DestinationAddress, testCase.shipmentType) {
					// Original and new IDs should not match since we should be creating an entirely new address record
					suite.NotEqual(subtestData.move.Orders.NewDutyLocation.Address.ID, createdShipment.DestinationAddress.ID, testCase.shipmentType)
					suite.Equal(*createdShipment.DestinationAddressID, createdShipment.DestinationAddress.ID)
					suite.Equal("N/A", createdShipment.DestinationAddress.StreetAddress1, testCase.shipmentType)
					suite.Nil(createdShipment.DestinationAddress.StreetAddress2, testCase.shipmentType)
					suite.Nil(createdShipment.DestinationAddress.StreetAddress3, testCase.shipmentType)
					suite.Equal(subtestData.move.Orders.NewDutyLocation.Address.City, createdShipment.DestinationAddress.City, testCase.shipmentType)
					suite.Equal(subtestData.move.Orders.NewDutyLocation.Address.State, createdShipment.DestinationAddress.State, testCase.shipmentType)
					suite.Equal(subtestData.move.Orders.NewDutyLocation.Address.PostalCode, createdShipment.DestinationAddress.PostalCode, testCase.shipmentType)
				}
			} else {
				suite.Nil(createdShipment.DestinationAddressID, testCase.shipmentType)
				suite.Nil(createdShipment.DestinationAddress, testCase.shipmentType)
			}
		}
	})

	suite.Run("If the shipment is created successfully with submitted status it should be returned", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.Equal(models.MTOShipmentStatusSubmitted, createdShipment.Status)
	})

	suite.Run("If the submitted shipment has a storage facility attached", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		storageFacility := factory.BuildStorageFacility(nil, nil, nil)
		// stubbed storage facility needs an ID to be LinkOnly below
		storageFacility.ID = uuid.Must(uuid.NewV4())

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: futureDate,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)
		suite.NoError(err)
		suite.NotNil(createdShipment.StorageFacility)
		suite.Equal(storageFacility.Address.StreetAddress1, createdShipment.StorageFacility.Address.StreetAddress1)
	})

	suite.Run("If the submitted shipment is an NTS shipment", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		ntsRecordedWeight := unit.Pound(980)
		requestedDeliveryDate := time.Date(testdatagen.GHCTestYear, time.April, 5, 0, 0, 0, 0, time.UTC)
		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:          models.MTOShipmentTypeHHGOutOfNTS,
					Status:                models.MTOShipmentStatusSubmitted,
					NTSRecordedWeight:     &ntsRecordedWeight,
					RequestedDeliveryDate: &requestedDeliveryDate,
					RequestedPickupDate:   futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)
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
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
					Status:       models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
	})

	suite.Run("When NTSRecordedWeight it set for a non NTS Release shipment", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		ntsRecordedWeight := unit.Pound(980)
		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					NTSRecordedWeight:   &ntsRecordedWeight,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		ntsrShipmentNoIDs := clearShipmentIDFields(&mtoShipment)
		ntsrShipmentNoIDs.RequestedPickupDate = futureDate

		// We don't need the shipment because it only returns data that wasn't saved.
		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), ntsrShipmentNoIDs)

		if suite.Errorf(err, "should have errored for a %s shipment with ntsRecordedWeight set", ntsrShipmentNoIDs.ShipmentType) {
			suite.IsType(apperror.InvalidInputError{}, err)
			suite.Contains(err.Error(), "NTSRecordedWeight")
		}
	})

	suite.Run("If the shipment has mto service items", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeDDSHUT,
			models.ReServiceCodeDOFSIT,
		}

		for _, serviceCode := range expectedReServiceCodes {
			factory.FetchReServiceByCode(suite.DB(), serviceCode)
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

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
		mtoShipmentClear.MTOServiceItems = serviceItemsList

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.NoError(err)
		suite.NotNil(createdShipment)
		suite.NotNil(createdShipment.MTOServiceItems, "Service Items are empty")
		suite.Equal(createdShipment.MTOServiceItems[0].MTOShipmentID, &createdShipment.ID, "Service items are not the same")
	})

	suite.Run("422 Validation Error - only one mto agent of each type", func() {
		subtestData := suite.createSubtestData(nil)
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

		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MTOAgents:           agents,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		shipment.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), &shipment)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("403 Forbidden Error - shipment can only be created for service member associated with the current session", func() {
		subtestData := suite.createSubtestData(nil)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: subtestData.move.Orders.ServiceMember.ID,
		})
		creator := subtestData.shipmentCreator
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.FromStringOrNil("424d932b-cf8d-4c10-8059-be8a25ba952a"),
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ServiceMember{
					ID: uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a"),
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&shipment)
		mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}
		createdShipment, err := creator.CreateMTOShipment(appCtx, mtoShipmentClear)

		suite.Nil(createdShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Will not create MTO agent if all fields are empty", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		firstName := ""
		lastName := ""
		email := ""

		var agents models.MTOAgents

		agent1 := models.MTOAgent{
			FirstName:    &firstName,
			LastName:     &lastName,
			Email:        &email,
			MTOAgentType: models.MTOAgentReceiving,
		}

		agents = append(agents, agent1)

		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		clearedShipment := clearShipmentIDFields(&shipment)

		clearedShipment.MTOAgents = agents
		clearedShipment.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedShipment)

		suite.NoError(err)
		suite.Len(createdShipment.MTOAgents, 0)
	})

	suite.Run("Move status transitions when a new shipment is created and SUBMITTED", func() {
		// If a new shipment is added to an APPROVED move and given the SUBMITTED status,
		// the move should transition to "APPROVALS REQUESTED"
		subtestData := suite.createSubtestData([]factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		})
		creator := subtestData.shipmentCreator
		move := subtestData.move
		shipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		cleanShipment := clearShipmentIDFields(&shipment)
		cleanShipment.MTOServiceItems = models.MTOServiceItems{}

		createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), cleanShipment)

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
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt

			var mtoShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				mtoShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			} else {
				mtoShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        tt.shipmentType,
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			}

			clearedShipment := clearShipmentIDFields(&mtoShipment)

			createdShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedShipment)

			suite.NoError(err, tt.desc)

			suite.Equal(models.DefaultServiceMemberSITDaysAllowance, *createdShipment.SITDaysAllowance, tt.desc)
		}
	})

	suite.Run("Test successful diversion from non-diverted parent shipment", func() {
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt
			var err error

			var parentShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				parentShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			} else {
				parentShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        tt.shipmentType,
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			}

			clearedParentShipment := clearShipmentIDFields(&parentShipment)

			createdParentShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedParentShipment)
			suite.NoError(err)

			// Create a new shipment, diverting from the parent
			var childShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				childShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							Diversion:              true,
							DivertedFromShipmentID: &createdParentShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				childShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							Diversion:              true,
							DivertedFromShipmentID: &createdParentShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedChildShipment := clearShipmentIDFields(&childShipment)
			clearedChildShipment.PrimeActualWeight = nil

			_, err = creator.CreateMTOShipment(suite.AppContextForTest(), clearedChildShipment)
			suite.NoError(err)
		}
	})

	suite.Run("Test successful diversion from parent shipment that itself is a diversion as well", func() {
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt
			var err error

			var unDivertedParentShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				unDivertedParentShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			} else {
				unDivertedParentShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        tt.shipmentType,
							RequestedPickupDate: futureDate,
						},
					},
				}, nil)
			}

			clearedUndivertedParentShipment := clearShipmentIDFields(&unDivertedParentShipment)

			createdUndivertedParentShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedUndivertedParentShipment)
			suite.NoError(err)

			// Create a new shipment, diverting from the parent
			var childFromParentDivertedShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				childFromParentDivertedShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							Diversion:              true,
							DivertedFromShipmentID: &createdUndivertedParentShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				childFromParentDivertedShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							Diversion:              true,
							DivertedFromShipmentID: &createdUndivertedParentShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedChildFromParentDivertedShipment := clearShipmentIDFields(&childFromParentDivertedShipment)
			clearedChildFromParentDivertedShipment.PrimeActualWeight = nil

			createdChildFromParentDivertedShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedChildFromParentDivertedShipment)
			suite.NoError(err)

			// Create a new shipment, diverting from the parent
			var childOfDivertedShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				childOfDivertedShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							Diversion:              true,
							DivertedFromShipmentID: &createdChildFromParentDivertedShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				childOfDivertedShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							Diversion:              true,
							DivertedFromShipmentID: &createdChildFromParentDivertedShipment.ID,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedChildOfDivertedShipment := clearShipmentIDFields(&childOfDivertedShipment)
			clearedChildOfDivertedShipment.PrimeActualWeight = nil
			_, err = creator.CreateMTOShipment(suite.AppContextForTest(), clearedChildOfDivertedShipment)
			suite.NoError(err)
		}
	})

	suite.Run("If DivertedFromShipmentID doesn't exist", func() {
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt
			uuid, _ := uuid.NewV4()

			var parentShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				parentShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							Diversion:              true,
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				parentShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							Diversion:              true,
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedParentShipment := clearShipmentIDFields(&parentShipment)

			_, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedParentShipment)
			suite.Error(err)
		}
	})

	suite.Run("If DivertedFromShipmentID is provided without the Diversion boolean", func() {
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt
			uuid, _ := uuid.NewV4()
			var parentShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				parentShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							Diversion:              false,
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				parentShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							Diversion:              false,
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedParentShipment := clearShipmentIDFields(&parentShipment)

			_, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedParentShipment)
			suite.Error(err)
		}
	})

	suite.Run("If DivertedFromShipmentID is provided to the V1 endpoint it should fail", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		testCases := []struct {
			desc         string
			shipmentType models.MTOShipmentType
		}{
			{"HHG", models.MTOShipmentTypeHHG},
			{"HHG_INTO_NTS", models.MTOShipmentTypeHHGIntoNTS},
			{"HHG_OUTOF_NTS", models.MTOShipmentTypeHHGOutOfNTS},
			{"MOBILE_HOME", models.MTOShipmentTypeMobileHome},
			{"BOAT_HAUL_AWAY", models.MTOShipmentTypeBoatHaulAway},
			{"BOAT_TOW_AWAY", models.MTOShipmentTypeBoatTowAway},
			{"PPM", models.MTOShipmentTypePPM},
			{"UNACCOMPANIED_BAGGAGE", models.MTOShipmentTypeUnaccompaniedBaggage},
		}

		for _, tt := range testCases {
			tt := tt
			uuid, _ := uuid.NewV4()
			var parentShipment models.MTOShipment
			if tt.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				parentShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			} else {
				parentShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:           tt.shipmentType,
							DivertedFromShipmentID: &uuid,
							RequestedPickupDate:    futureDate,
						},
					},
				}, nil)
			}

			clearedParentShipment := clearShipmentIDFields(&parentShipment)

			_, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedParentShipment)
			suite.Error(err)
		}
	})

	suite.Run("Child diversion shipment creation should inherit parent's weight", func() {
		currentTime := time.Now()

		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
					ApprovedAt:         &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: nil,
					RequestedPickupDate:    futureDate,
				},
			},
		}, nil)
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator
		childShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: &parentShipment.ID,
					RequestedPickupDate:    futureDate,
				},
			},
		}, nil)

		clearedChildShipment := clearShipmentIDFields(&childShipment)
		clearedChildShipment.PrimeActualWeight = nil
		clearedChildShipment.DivertedFromShipmentID = &parentShipment.ID

		createdChildShipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedChildShipment)
		suite.NoError(err)
		suite.Equal(createdChildShipment.PrimeActualWeight, parentShipment.PrimeActualWeight)
	})
	suite.Run("Child diversion shipment creation should fail if PrimeActualWeight is provided", func() {
		currentTime := time.Now()
		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
					ApprovedAt:         &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: nil,
					RequestedPickupDate:    futureDate,
				},
			},
		}, nil)
		subtestData := suite.createSubtestDataV2(nil)
		creator := subtestData.shipmentCreator
		childShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: &parentShipment.ID,
					RequestedPickupDate:    futureDate,
				},
			},
		}, nil)

		// prmie actual weight is auto supplied
		clearedChildShipment := clearShipmentIDFields(&childShipment)

		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), clearedChildShipment)
		suite.Error(err)
	})

	suite.Run("InvalidInputError - NTS shipment cannot specify a secondary delivery address", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		pickupAddress := factory.BuildDefaultAddress(suite.DB())
		deliveryAddress := factory.BuildDefaultAddress(suite.DB())
		secondaryDeliveryAddress := factory.BuildDefaultAddress(suite.DB())

		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    deliveryAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    secondaryDeliveryAddress,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		clearShipmentIDFields(&shipment)

		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), &shipment)

		suite.Error(err)
		suite.Equal("Secondary delivery address cannot be created for shipment Type "+string(models.MTOShipmentTypeHHGIntoNTS), err.Error())
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("InvalidInputError - NTSR shipment cannot specify a secondary pickup address", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		pickupAddress := factory.BuildDefaultAddress(suite.DB())
		deliveryAddress := factory.BuildDefaultAddress(suite.DB())
		secondaryPickupAddress := factory.BuildDefaultAddress(suite.DB())

		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    deliveryAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    secondaryPickupAddress,
				Type:     &factory.Addresses.SecondaryPickupAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
					RequestedPickupDate: futureDate,
				},
			},
		}, nil)
		clearShipmentIDFields(&shipment)

		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), &shipment)

		suite.Error(err)
		suite.Equal("Secondary pickup address cannot be created for shipment Type "+string(models.MTOShipmentTypeHHGOutOfNTS), err.Error())
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("UB shipments will return an error if both addresses are CONUS", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
					RequestedPickupDate: futureDate,
				},
			},
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		mtoShipmentClear := clearShipmentIDFields(&mtoShipment)

		_, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

		suite.Error(err)
		suite.Equal("At least one address for a UB shipment must be OCONUS", err.Error())
	})

	suite.Run("RequestedPickupDate validation check - must be in the future for shipment types other than PPM", func() {
		subtestData := suite.createSubtestData(nil)
		creator := subtestData.shipmentCreator

		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		tomorrow := now.AddDate(0, 0, 1)

		testCases := []struct {
			input        *time.Time
			shipmentType models.MTOShipmentType
			shouldError  bool
		}{
			// HHG
			{nil, models.MTOShipmentTypeHHG, true},
			{&time.Time{}, models.MTOShipmentTypeHHG, true},
			{&yesterday, models.MTOShipmentTypeHHG, true},
			{&now, models.MTOShipmentTypeHHG, true},
			{&tomorrow, models.MTOShipmentTypeHHG, false},
			// NTS
			{nil, models.MTOShipmentTypeHHGIntoNTS, true},
			{&time.Time{}, models.MTOShipmentTypeHHGIntoNTS, true},
			{&yesterday, models.MTOShipmentTypeHHGIntoNTS, true},
			{&now, models.MTOShipmentTypeHHGIntoNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGIntoNTS, false},
			// NTSR
			{nil, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&time.Time{}, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&yesterday, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&now, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGOutOfNTS, false},
			// BOAT HAUL AWAY
			{nil, models.MTOShipmentTypeBoatHaulAway, false},
			{&time.Time{}, models.MTOShipmentTypeBoatHaulAway, false},
			{&yesterday, models.MTOShipmentTypeBoatHaulAway, true},
			{&now, models.MTOShipmentTypeBoatHaulAway, true},
			{&tomorrow, models.MTOShipmentTypeBoatHaulAway, false},
			// BOAT TOW AWAY
			{nil, models.MTOShipmentTypeBoatTowAway, false},
			{&time.Time{}, models.MTOShipmentTypeBoatTowAway, false},
			{&yesterday, models.MTOShipmentTypeBoatTowAway, true},
			{&now, models.MTOShipmentTypeBoatTowAway, true},
			{&tomorrow, models.MTOShipmentTypeBoatTowAway, false},
			// MOBILE HOME
			{nil, models.MTOShipmentTypeMobileHome, false},
			{&time.Time{}, models.MTOShipmentTypeMobileHome, false},
			{&yesterday, models.MTOShipmentTypeMobileHome, true},
			{&now, models.MTOShipmentTypeMobileHome, true},
			{&tomorrow, models.MTOShipmentTypeMobileHome, false},
			// UB
			{nil, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&time.Time{}, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&yesterday, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&now, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&tomorrow, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			// PPM - should always pass validation
			{nil, models.MTOShipmentTypePPM, false},
			{&time.Time{}, models.MTOShipmentTypePPM, false},
			{&yesterday, models.MTOShipmentTypePPM, false},
			{&now, models.MTOShipmentTypePPM, false},
			{&tomorrow, models.MTOShipmentTypePPM, false},
		}

		for _, testCase := range testCases {
			// Default is HHG, but we set it explicitly below via the test cases
			var mtoShipment models.MTOShipment
			if testCase.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				mtoShipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model: models.MTOShipment{
							ShipmentType: testCase.shipmentType,
						},
					},
				}, nil)
			} else {
				mtoShipment = factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model:    subtestData.move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType: testCase.shipmentType,
						},
					},
				}, nil)
			}

			mtoShipment.RequestedPickupDate = testCase.input // Zero case does not merge correctly on customization

			mtoShipmentClear := clearShipmentIDFields(&mtoShipment)
			mtoShipmentClear.MTOServiceItems = models.MTOServiceItems{}

			shipment, err := creator.CreateMTOShipment(suite.AppContextForTest(), mtoShipmentClear)

			testCaseInputString := ""
			if testCase.input == nil {
				testCaseInputString = "nil"
			} else {
				testCaseInputString = (*testCase.input).String()
			}

			if testCase.shouldError {
				suite.Nil(shipment, "Should error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.Error(err)
				if testCase.input != nil && !(*testCase.input).IsZero() {
					suite.Equal("RequestedPickupDate must be greater than or equal to tomorrow's date.", err.Error())
				} else {
					suite.Contains(err.Error(), fmt.Sprintf("RequestedPickupDate is required to create or modify %s %s shipment", GetAorAnByShipmentType(testCase.shipmentType), testCase.shipmentType))
				}
			} else {
				suite.NoError(err, "Should not error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.NotNil(shipment)
			}
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
	if shipment.HasTertiaryPickupAddress != nil {
		shipment.TertiaryPickupAddressID = nil
		shipment.TertiaryPickupAddress.ID = uuid.Nil
	}

	if shipment.HasTertiaryDeliveryAddress != nil {
		shipment.TertiaryDeliveryAddressID = nil
		shipment.TertiaryDeliveryAddress.ID = uuid.Nil
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
