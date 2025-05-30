package shipmentaddressupdate

import (
	"fmt"
	"slices"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ShipmentAddressUpdateServiceSuite) setupServiceItemData() {
	startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, time.December, 31, 12, 0, 0, 0, time.UTC)

	testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			StartDate: startDate,
			EndDate:   endDate,
		},
	})

	originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.AppContextForTest().DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea:      "004",
			ServicesSchedule: 2,
		},
		ReContract: testdatagen.FetchOrMakeReContract(suite.AppContextForTest().DB(), testdatagen.Assertions{}),
	})

	testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			Contract:              originalDomesticServiceArea.Contract,
			ContractID:            originalDomesticServiceArea.ContractID,
			DomesticServiceArea:   originalDomesticServiceArea,
			DomesticServiceAreaID: originalDomesticServiceArea.ID,
			WeightLower:           unit.Pound(500),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            500,
			MilesUpper:            9999,
			PriceMillicents:       unit.Millicents(606800),
			IsPeakPeriod:          false,
		},
	})
}

func (suite *ShipmentAddressUpdateServiceSuite) TestCreateApprovedShipmentAddressUpdate() {
	setupTestData := func() models.Move {
		originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		// Common ZIP3s used in these tests
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "895",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "902",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "945",
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		return move
	}
	addressCreator := address.NewAddressCreator()
	mockPlanner := &routemocks.Planner{}
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter)

	suite.Run("Successfully create ShipmentAddressUpdate for a domestic shipment", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Twice()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(2500, nil).Once()
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// New destination address with same postal code should not change pricing
		newAddress := models.Address{
			StreetAddress1: "987 Any Avenue",
			City:           "Fairfield",
			State:          "CA",
			PostalCode:     "94535",
		}
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)

		// Make sure the destination address on the shipment was updated
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress").Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(newAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(newAddress.PostalCode, updatedShipment.DestinationAddress.PostalCode)
		suite.Equal(newAddress.State, updatedShipment.DestinationAddress.State)
		suite.Equal(newAddress.City, updatedShipment.DestinationAddress.City)
	})

	suite.Run("Successfully create ShipmentAddressUpdate for an international shipment that requires approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Twice()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(2500, nil).Once()
		move := setupTestData()

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
				},
			},
		}, nil)

		newAddress := models.Address{
			StreetAddress1: "Colder Ave.",
			City:           "Klawock",
			State:          "AK",
			PostalCode:     "99925",
		}
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)

		// Make sure the destination address on the shipment was NOT updated
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress").Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		suite.NotEqual(newAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.NotEqual(newAddress.PostalCode, updatedShipment.DestinationAddress.PostalCode)
		suite.NotEqual(newAddress.City, updatedShipment.DestinationAddress.City)
	})

	suite.Run("Successfully create ShipmentAddressUpdate for an international shipment that requires approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"99505",
			"99506",
		).Return(49, nil)
		move := setupTestData()

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
				},
			},
		}, nil)

		// this shouldn't change the rate area
		newAddress := models.Address{
			StreetAddress1: "Elsewhere Ave.",
			City:           "Anchorage",
			State:          "AK",
			PostalCode:     "99506",
		}
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)

		// Make sure the destination address on the shipment was updated
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress").Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(newAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(newAddress.PostalCode, updatedShipment.DestinationAddress.PostalCode)
		suite.Equal(newAddress.City, updatedShipment.DestinationAddress.City)
	})

	suite.Run("Update with invalid etag should fail", func() {
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// New destination address with same postal code should not change pricing
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     shipment.DestinationAddress.PostalCode,
		}
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt.Add(-1)))
		suite.Error(err)
		suite.Nil(update)
	})

	suite.Run("Failed distance calculation should error", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(0, fmt.Errorf("error calculating distance 2")).Once()

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		// building DDASIT service item to get dest SIT checks
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:                          models.MTOServiceItemStatusApproved,
					SITDestinationOriginalAddressID: shipment.DestinationAddressID,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
		}, nil)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.Error(err)
		suite.Nil(update)
	})

	suite.Run("Should be able to use this service to update a shipment with origin SIT", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		year, month, day := time.Now().Date()
		lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
		lastMonthDeparture := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &lastMonthEntry,
					SITDepartureDate: &lastMonthDeparture,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
		}, nil)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
	})

	suite.Run("Should not be able to update invalid shipments delivery address", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		contractRemarks := "we really need to change the address"

		// NTS and PPM
		invalidTypes := []models.MTOShipment{
			factory.BuildNTSShipment(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypePPM,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
		}

		for i := 0; i < len(invalidTypes); i++ {
			update, err :=
				addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(),
					invalidTypes[i].ID, newAddress, contractRemarks,
					etag.GenerateEtag(invalidTypes[i].UpdatedAt))
			suite.Error(err)
			suite.Nil(update)
		}
	})

	suite.Run("Should be able to update valid shipments delivery address", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)
		contractRemarks := "we really need to change the address"

		// NTS-R, Mobile Home, Boat, UB, HHG
		validTypes := []models.MTOShipment{
			factory.BuildNTSRShipment(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    storageFacility,
					LinkOnly: true,
				},
			}, nil),
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypeMobileHome,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypeBoatHaulAway,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
						MarketCode:   models.MarketCodeInternational,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
			factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypeHHG,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
			}, nil),
		}

		for i := 0; i < len(validTypes); i++ {
			update, err :=
				addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(
					suite.AppContextForTest(), validTypes[i].ID, newAddress,
					contractRemarks,
					etag.GenerateEtag(validTypes[i].UpdatedAt))
			suite.NoError(err)
			suite.NotNil(update)
		}
	})

	suite.Run("Request destination address changes on the same shipment multiple times", func() {
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "FAIRFIELD",
			State:          "CA",
			PostalCode:     shipment.DestinationAddress.PostalCode,
		}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Times(4)
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(2500, nil).Twice()

		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		// Need to re-request the shipment to get the updated etag
		var updatedShipment models.MTOShipment
		err = suite.DB().Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		update, err = addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address again", etag.GenerateEtag(updatedShipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address again", update.ContractorRemarks)
	})

	suite.Run("Shorthaul to linehaul should be flagged", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "BEVERLY HILLS",
			State:          "CA",
			PostalCode:     "90210",
		}
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
					City:       "RENO",
					State:      "NV",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "89503",
					City:       "RENO",
					State:      "NV",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	suite.Run("linehaul to shorthaul should be flagged", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
					City:       "RENO",
					State:      "NV",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "RENO",
			State:          "NV",
			PostalCode:     "89503",
		}

		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	suite.Run("service area change should be flagged", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(0, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"89503",
		).Return(200, nil).Once()
		originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "902",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "945",
			},
		})
		newDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "005",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            newDomesticServiceArea.Contract,
				ContractID:          newDomesticServiceArea.ContractID,
				DomesticServiceArea: newDomesticServiceArea,
				Zip3:                "895",
			},
		})

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "RENO",
			State:          "NV",
			PostalCode:     "89503",
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "94535",
					City:       "FAIRFIELD",
					State:      "CA",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	suite.Run("mileage bracket change should be flagged", func() {
		originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "871",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "870",
			},
		})

		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "87108",
					City:       "ALBUQUERQUE",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "SAN YSIDRO",
			State:          "NM",
			PostalCode:     "87053",
		}

		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			"87108",
		).Return(500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			"87053",
		).Return(501, nil).Once()
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})

	suite.Run("destination address request succeeds when containing destination SIT", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// building DDASIT service item to get dest SIT checks
		serviceItemDDASIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:                          models.MTOServiceItemStatusApproved,
					SITDestinationOriginalAddressID: shipment.DestinationAddressID,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
		}, nil)

		// mock ZipTransitDistance function
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(0, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94523",
			"90210",
		).Return(500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"90210",
		).Return(501, nil).Once()

		// request the update
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)

		// querying the address update to make sure that SIT data was populated
		var addressUpdate models.ShipmentAddressUpdate
		err = suite.DB().Find(&addressUpdate, update.ID)
		suite.NoError(err)
		suite.Equal(*addressUpdate.NewSitDistanceBetween, 501)
		suite.Equal(*addressUpdate.OldSitDistanceBetween, 0)
		suite.Equal(*addressUpdate.SitOriginalAddressID, *serviceItemDDASIT.SITDestinationOriginalAddressID)
	})

	suite.Run("destination address request succeeds when containing international destination SIT", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Anchorage",
			State:          "AK",
			PostalCode:     "99695",
		}

		setupInternationalSITCodes := []models.ReServiceCode{
			models.ReServiceCodeIDASIT,
			models.ReServiceCodeIDDSIT,
			models.ReServiceCodeIDFSIT,
			models.ReServiceCodeIDSFSC,
		}

		// loop through test codes to verify updates are applied for expected international SITs
		for _, reServiceCode := range setupInternationalSITCodes {
			shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
			}, nil)

			// building service item to get dest SIT checks
			serviceItemDDASIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.MTOServiceItem{
						Status:                          models.MTOServiceItemStatusApproved,
						SITDestinationOriginalAddressID: shipment.DestinationAddressID,
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    shipment,
					LinkOnly: true,
				},
				{
					Model: models.ReService{
						Code: reServiceCode,
					},
				},
			}, nil)

			// mock ZipTransitDistance function
			mockPlanner.On("ZipTransitDistance",
				mock.AnythingOfType("*appcontext.appContext"),
				"94535",
				"94535",
			).Return(0, nil).Once()
			mockPlanner.On("ZipTransitDistance",
				mock.AnythingOfType("*appcontext.appContext"),
				"94523",
				"99695",
			).Return(500, nil).Once()
			mockPlanner.On("ZipTransitDistance",
				mock.AnythingOfType("*appcontext.appContext"),
				"94535",
				"99695",
			).Return(1000, nil).Once()

			// request the update
			update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
			suite.NoError(err)
			suite.NotNil(update)

			// querying the address update to make sure that SIT data was populated
			var addressUpdate models.ShipmentAddressUpdate
			err = suite.DB().Find(&addressUpdate, update.ID)
			suite.NoError(err)
			suite.Equal(*addressUpdate.NewSitDistanceBetween, 1000)
			suite.Equal(*addressUpdate.OldSitDistanceBetween, 0)
			suite.Equal(*addressUpdate.SitOriginalAddressID, *serviceItemDDASIT.SITDestinationOriginalAddressID)
		}
	})
}

func (suite *ShipmentAddressUpdateServiceSuite) TestTOOApprovedShipmentAddressUpdateRequest() {
	addressCreator := address.NewAddressCreator()
	mockPlanner := &routemocks.Planner{}
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter)

	suite.Run("TOO approves address change", func() {

		suite.setupServiceItemData()

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})
		officeRemarks := "This is a TOO remark"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

	})

	suite.Run("TOO rejects address change", func() {

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})
		officeRemarks := "This is a TOO remark"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "REJECTED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRejected, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

	})

	suite.Run("TOO approves address change and left no remarks", func() {

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)
		officeRemarks := ""

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Nil(update)
	})

	suite.Run("After TOO approval, move transitions from approvals requested to approved", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{{
			Model: models.Move{
				Status: "APPROVALS REQUESTED",
			},
		}}, nil)
		officeRemarks := "Looks good!"

		var updatedMove models.Move
		err := suite.DB().Find(&updatedMove, addressChange.Shipment.MoveTaskOrderID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)
		suite.NoError(err)

		err = suite.DB().Find(&updatedMove, addressChange.Shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)

	})

	suite.Run("TOO approves address change and service items final destination address changes", func() {
		// creating an address change that shares the same address to avoid hitting lineHaulChange check
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: models.StringPointer("Apt 2"),
					StreetAddress3: models.StringPointer("Suite 200"),
					City:           "New York",
					State:          "NY",
					PostalCode:     "10001",
				},
				Type: &factory.Addresses.OriginalAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: models.StringPointer("Apt 2"),
					StreetAddress3: models.StringPointer("Suite 200"),
					City:           "New York",
					State:          "NY",
					PostalCode:     "10001",
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, nil)
		shipment := addressChange.Shipment
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		sitDestinationOriginalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: shipment.MoveTaskOrderID,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    sitDestinationOriginalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
			},
		}, nil)
		officeRemarks := "This is a TOO remark"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

		// Make sure the destination address on the shipment was updated
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress", "MTOServiceItems").Find(&updatedShipment, update.ShipmentID)
		suite.NoError(err)

		// service item status should be changed to submitted
		suite.Equal(models.MTOServiceItemStatusSubmitted, updatedShipment.MTOServiceItems[0].Status)
		// delivery and final destination addresses should be the same
		suite.Equal(updatedShipment.DestinationAddressID, updatedShipment.MTOServiceItems[0].SITDestinationFinalAddressID)
	})

	suite.Run("TOO approves address change that triggers market code change of shipment", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: models.StringPointer("Apt 2"),
					StreetAddress3: models.StringPointer("Suite 200"),
					City:           "New York",
					State:          "NY",
					PostalCode:     "10001",
				},
				Type: &factory.Addresses.OriginalAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "456 Northern Lights Blvd",
					StreetAddress2: models.StringPointer("Apt 5B"),
					StreetAddress3: models.StringPointer("Suite 300"),
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99503",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, nil)
		shipment := addressChange.Shipment
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		sitDestinationOriginalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: shipment.MoveTaskOrderID,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    sitDestinationOriginalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
			},
		}, nil)
		officeRemarks := "Changing to OCONUS address"

		// check to make sure the market code is "d" prior to updating with OCONUS address
		suite.Equal(shipment.MarketCode, models.MarketCodeDomestic)
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)
		suite.NoError(err)
		suite.NotNil(update)

		// Make sure the market code changed on the shipment
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress", "MTOServiceItems").Find(&updatedShipment, update.ShipmentID)
		suite.NoError(err)

		suite.Equal(updatedShipment.MarketCode, models.MarketCodeInternational)
	})

	suite.Run("TOO approves address change on invalid ShipmentType Mobile Home", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeMobileHome,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: models.StringPointer("Apt 2"),
					StreetAddress3: models.StringPointer("Suite 200"),
					City:           "New York",
					State:          "NY",
					PostalCode:     "10001",
				},
				Type: &factory.Addresses.OriginalAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "456 Northern Lights Blvd",
					StreetAddress2: models.StringPointer("Apt 5B"),
					StreetAddress3: models.StringPointer("Suite 300"),
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99503",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, nil)
		shipment := addressChange.Shipment
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		sitDestinationOriginalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: shipment.MoveTaskOrderID,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    sitDestinationOriginalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
			},
		}, nil)
		officeRemarks := "Changing to OCONUS address"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)
		suite.ErrorContains(err, "Shipment type must be HHG, NTSr or UB")
		suite.Nil(update)
	})

	suite.Run("TOO approves address change on invalid ShipmentType Boat", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeBoatHaulAway,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: models.StringPointer("Apt 2"),
					StreetAddress3: models.StringPointer("Suite 200"),
					City:           "New York",
					State:          "NY",
					PostalCode:     "10001",
				},
				Type: &factory.Addresses.OriginalAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "456 Northern Lights Blvd",
					StreetAddress2: models.StringPointer("Apt 5B"),
					StreetAddress3: models.StringPointer("Suite 300"),
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99503",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, nil)
		shipment := addressChange.Shipment
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		sitDestinationOriginalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: shipment.MoveTaskOrderID,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    sitDestinationOriginalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationOriginalAddress,
			},
		}, nil)
		officeRemarks := "Changing to OCONUS address"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)
		suite.ErrorContains(err, "Shipment type must be HHG, NTSr or UB")
		suite.Nil(update)
	})

	suite.Run("Successfully update estiamted pricing on service items when address update is approved by TOO", func() {
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		now := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: models.PoundPointer(4000),
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &now,
					RequestedPickupDate:  &now,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: models.CentPointer(1000),
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: models.CentPointer(1000),
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHUPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: models.CentPointer(1000),
				},
			},
		}, nil)
		// POEFSC needs a port location
		portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: models.CentPointer(1000),
				},
			},
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfEmbarkation,
			},
		}, nil)

		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(300, nil)

		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(300, nil)

		newDestUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99703")
		suite.FatalNoError(err)
		factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.Address{
					StreetAddress1:     "Cold Ave.",
					City:               "FORT WAINWRIGHT",
					State:              "AK",
					PostalCode:         "99703",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &newDestUSPRC.ID,
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		officeRemarks := "Changing to another OCONUS address"
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), shipment.ID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)

		// checking out the service items
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodePOEFSC,
			models.ReServiceCodeISLH,
			models.ReServiceCodeIHPK,
			models.ReServiceCodeIHUPK,
		}

		initialPrice := 1000
		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
			// pricing should not be nil
			suite.NotNil(serviceItems[i].PricingEstimate)
			// initially we set them all to 1000 and they should all be changed
			suite.NotEqual(serviceItems[i].PricingEstimate, &initialPrice)
		}
	})
	suite.Run("On approval - successfully update estimated pricing on all basic iHHG service items except for POEFSC when port isn't set", func() {
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		now := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: models.PoundPointer(4000),
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &now,
					RequestedPickupDate:  &now,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHUPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: models.CentPointer(1000),
				},
			},
		}, nil)
		// no port data
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)

		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(300, nil)

		newDestUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99703")
		suite.FatalNoError(err)
		factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.Address{
					StreetAddress1:     "Cold Ave.",
					City:               "FORT WAINWRIGHT",
					State:              "AK",
					PostalCode:         "99703",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &newDestUSPRC.ID,
				},
				Type: &factory.Addresses.NewAddress,
			},
		}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		officeRemarks := "Changing to another OCONUS address"
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), shipment.ID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)

		// checking out the service items
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			if serviceItems[i].ReService.Code != models.ReServiceCodePOEFSC {
				suite.NotNil(serviceItems[i].PricingEstimate)
			} else if serviceItems[i].ReService.Code == models.ReServiceCodePOEFSC {
				suite.Nil(serviceItems[i].PricingEstimate)
			}
		}
	})
}

func (suite *ShipmentAddressUpdateServiceSuite) TestTOOApprovedShipmentAddressUpdateRequestChangedPricing() {
	setupTestData := func() models.Move {
		originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		// Common ZIP3s used in these tests
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "895",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "902",
			},
		})
		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originalDomesticServiceArea.Contract,
				ContractID:          originalDomesticServiceArea.ContractID,
				DomesticServiceArea: originalDomesticServiceArea,
				Zip3:                "945",
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		return move
	}
	addressCreator := address.NewAddressCreator()
	mockPlanner := &routemocks.Planner{}
	moveRouter := moveservices.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter)

	suite.Run("Service items are rejected and regenerated when pricing type changes post TOO approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		move := setupTestData()

		suite.setupServiceItemData()

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
					IsOconus:   models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
					IsOconus:   models.BoolPointer(false),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		//Generate a couple of service items to test their status changes upon approval
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeMS, move, shipment, nil, nil)
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDLH, move, shipment, nil, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "RENO",
			State:          "NV",
			PostalCode:     "89503",
		}

		// Trigger the prime address update to get move in correct state for DLH -> DSH
		addressChange, _ := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		officeRemarks := "This is a TOO remark"

		// TOO Approves address change
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.ShipmentID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

		// Assert that the DLH service item was rejected and has the correct rejection reason
		rejectedServiceItems := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDLH)
		suite.Equal(rejectedServiceItems[0].Status, models.MTOServiceItemStatusRejected)
		autoRejectionRemark := "Automatically rejected due to change in destination address affecting the ZIP code qualification for short haul / line haul."
		suite.Equal(autoRejectionRemark, *rejectedServiceItems[0].RejectionReason)

		// Assert that the DSH service was created and is in an approved state
		approvedServiceItems := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDSH)
		suite.Equal(approvedServiceItems[0].Status, models.MTOServiceItemStatusApproved)

		// Should have an equal number of rejected and approved service items
		suite.Equal(len(approvedServiceItems), len(rejectedServiceItems))
	})

	suite.Run("Service items were already rejected are not regenerated when pricing type changes post TOO approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		move := setupTestData()

		suite.setupServiceItemData()

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		//Generate a couple of service items to test their status changes upon approval
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDLH, move, shipment, nil, nil)
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeMS, move, shipment, []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusRejected,
				},
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "RENO",
			State:          "NV",
			PostalCode:     "89503",
		}

		// Trigger the prime address update to get move in correct state for DLH -> DSH
		addressChange, _ := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		officeRemarks := "This is a TOO remark"

		// TOO Approves address change
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.ShipmentID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

		// Assert that only the service items that weren't already rejected were the ones regenerated
		rejectedServiceItems := suite.getServiceItemsByStatus(update.Shipment.MTOServiceItems, models.MTOServiceItemStatusRejected)
		approvedServiceItems := suite.getServiceItemsByStatus(update.Shipment.MTOServiceItems, models.MTOServiceItemStatusApproved)

		// Should have 2 rejected service items and only 1 approved
		suite.Equal(len(approvedServiceItems), 1)
		suite.Equal(len(rejectedServiceItems), 2)
	})

	suite.Run("Service items are not rejected when pricing type does not change post TOO approval", func() {
		move := setupTestData()

		suite.setupServiceItemData()

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		//Generate service items to test their statuses upon approval
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeMS, move, shipment, nil, nil)
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDLH, move, shipment, nil, nil)

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "FAIRFIELD",
			State:          "CA",
			PostalCode:     shipment.DestinationAddress.PostalCode,
		}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(2500, nil).Once()

		addressChange, _ := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		officeRemarks := "This is a TOO remark"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.ShipmentID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)
		// Assert that service item was not rejected
		rejectedServiceItems := suite.getServiceItemsByStatus(update.Shipment.MTOServiceItems, models.MTOServiceItemStatusRejected)
		suite.Equal(len(rejectedServiceItems), 0)
	})

	suite.Run("Linehaul to shorthaul generates appropriate service items post TOO approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		move := setupTestData()

		suite.setupServiceItemData()

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
				},
				Type: &factory.Addresses.PickupAddress,
			},
		}, nil)
		//Generate a couple of service items to test their status changes upon approval
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeMS, move, shipment, nil, nil)
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDLH, move, shipment, nil, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "RENO",
			State:          "NV",
			PostalCode:     "89503",
		}

		// Trigger the prime address update to get move in correct state for DLH -> DSH
		addressChange, _ := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		officeRemarks := "This is a TOO remark"

		// TOO Approves address change
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.ShipmentID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

		// Confirm that DLH service item is rejected and DSH service item is created and approved
		linehaul := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDLH)
		shorthaul := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDSH)
		autoRejectionRemark := "Automatically rejected due to change in destination address affecting the ZIP code qualification for short haul / line haul."

		suite.NotNil(shorthaul)
		suite.Equal(linehaul[0].Status, models.MTOServiceItemStatusRejected)
		suite.Equal(autoRejectionRemark, *linehaul[0].RejectionReason)
		suite.Equal(shorthaul[0].Status, models.MTOServiceItemStatusApproved)
	})

	suite.Run("Shorthaul to linehaul generates appropriate service items post TOO approval", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"89503",
		).Return(2500, nil).Once()
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"89523",
			"90210",
		).Return(2500, nil).Once()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}
		move := setupTestData()

		suite.setupServiceItemData()

		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
					City:       "RENO",
					State:      "NV",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "89503",
					City:       "RENO",
					State:      "NV",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		//Generate a couple of service items to test their status changes upon approval
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeMS, move, shipment, nil, nil)
		factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDSH, move, shipment, nil, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)

		// Trigger the prime address update to get move in correct state for DSH -> DLH
		addressChange, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		officeRemarks := "This is a TOO remark"
		suite.NoError(err)

		// TOO Approves address change
		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.ShipmentID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

		// Confirm that DSH service item is rejected and DLH service item is created and approved
		linehaul := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDLH)
		shorthaul := suite.getServiceItemsByCode(update.Shipment.MTOServiceItems, models.ReServiceCodeDSH)
		autoRejectionRemark := "Automatically rejected due to change in destination address affecting the ZIP code qualification for short haul / line haul."

		suite.NotNil(shorthaul)
		suite.Equal(shorthaul[0].Status, models.MTOServiceItemStatusRejected)
		suite.Equal(autoRejectionRemark, *shorthaul[0].RejectionReason)
		suite.Equal(linehaul[0].Status, models.MTOServiceItemStatusApproved)
	})

	suite.Run("Successfully update shipment and its service items without error", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"94535",
			"94535",
		).Return(30, nil)
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		//Generate a couple of service items to test their status changes upon approval
		serviceItem1 := factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDOASIT, move, shipment, nil, nil)

		var serviceItems models.MTOServiceItems
		shipment.MTOServiceItems = append(serviceItems, serviceItem1)

		newAddress := models.Address{
			StreetAddress1: shipment.DestinationAddress.StreetAddress1,
			City:           shipment.DestinationAddress.City,
			State:          shipment.DestinationAddress.State,
			PostalCode:     shipment.DestinationAddress.PostalCode,
			Country:        shipment.DestinationAddress.Country,
		}

		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "Submitting same address", etag.GenerateEtag(shipment.UpdatedAt))

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
	})
}

func (suite *ShipmentAddressUpdateServiceSuite) getServiceItemsByStatus(items models.MTOServiceItems, status models.MTOServiceItemStatus) models.MTOServiceItems {
	itemsWithStatus := models.MTOServiceItems{}

	for _, si := range items {
		if si.Status == status {
			itemsWithStatus = append(itemsWithStatus, si)
		}
	}

	return itemsWithStatus
}

func (suite *ShipmentAddressUpdateServiceSuite) getServiceItemsByCode(items models.MTOServiceItems, code models.ReServiceCode) models.MTOServiceItems {
	itemsWithCode := models.MTOServiceItems{}

	for _, si := range items {
		if si.ReService.Code == code {
			itemsWithCode = append(itemsWithCode, si)
		}
	}

	return itemsWithCode
}
