package shipmentaddressupdate

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Test cases we should cover:
// test auto approve
// test update existing record
// test create new address update record
// test flag due to service area
// test flag due to mileage bracket
// test flag due to shorthaul to linehaul
// test flag due to linehaul to shorthaul

// test update non hhg shipment - error
// test update shipment with SIT - error
// test unable to look up distance error
// do we auto reject in this case? or return an error to the prime without creating address update?

func (suite *ShipmentAddressUpdateServiceSuite) TestCreateApprovedShipmentAddressUpdate() {
	addressCreator := address.NewAddressCreator()
	shipmentSITStatus := mtoshipment.NewShipmentSITStatus()
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		"89503",
		"90210",
	).Return(450, nil)
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		"90210",
		"30905",
	).Return(2500, nil)
	moveRouter := moveservices.NewMoveRouter()
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter, shipmentSITStatus)

	suite.Run("Successfully create ShipmentAddressUpdate", func() {
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)

		// Make sure the destination address on the shipment was updated
		var updatedShipment models.MTOShipment
		err = suite.DB().EagerPreload("DestinationAddress").Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(newAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(newAddress.PostalCode, updatedShipment.DestinationAddress.PostalCode)
		suite.Equal(newAddress.State, updatedShipment.DestinationAddress.State)
		suite.Equal(newAddress.City, updatedShipment.DestinationAddress.City)
	})

	suite.Run("Failed distance calculation should error", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(0, fmt.Errorf("error calculating distance")).Once()

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.Error(err)
		suite.Nil(update)
	})

	suite.Run("Should not be able to use this service to update a shipment with SIT", func() {
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
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
					Code: models.ReServiceCodeDOPSIT,
				},
			},
		}, nil)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.Error(err)
		suite.Nil(update)
	})
	suite.Run("Should not be able to update NTS shipment", func() {
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		shipment := factory.BuildNTSShipment(suite.DB(), nil, nil)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.Error(err)
		suite.Nil(update)
	})
	suite.Run("Request destination address changes on the same shipment multiple times", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)
		update, err = addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address again")
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("we really need to change the address again", update.ContractorRemarks)
	})
	suite.Run("Shorthaul to linehaul should be flagged", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "89523",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "89503",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "89503",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
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
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
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
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "89503",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
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
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
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

		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Albuquerque",
			State:          "NM",
			PostalCode:     "87053",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "87108",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			"87108",
		).Return(500, nil)
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			"87053",
		).Return(501, nil)
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address")
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		var updatedMove models.Move
		err = suite.DB().Find(&updatedMove, shipment.MoveTaskOrderID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, updatedMove.Status)
	})
}
