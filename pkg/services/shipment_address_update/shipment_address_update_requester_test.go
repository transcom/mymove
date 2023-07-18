package shipmentaddressupdate

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ShipmentAddressUpdateServiceSuite) TestCreateApprovedShipmentAddressUpdate() {
	setupTestData := func() models.Move {
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
	moveRouter := moveservices.NewMoveRouter()
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter)

	suite.Run("Successfully create ShipmentAddressUpdate", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Twice()
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)

		// New destination address with same postal code should not change pricing
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     shipment.DestinationAddress.PostalCode,
			Country:        models.StringPointer("United States"),
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
		suite.Equal(newAddress.State, updatedShipment.DestinationAddress.State)
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
			Country:        models.StringPointer("United States"),
		}
		suite.NotEmpty(move.MTOShipments)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt.Add(-1)))
		suite.Error(err)
		suite.Nil(update)
	})

	suite.Run("Failed distance calculation should error", func() {
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("90210"),
			mock.AnythingOfType("94535"),
		).Return(0, fmt.Errorf("error calculating distance 2")).Once()

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.Error(err)
		suite.Nil(update)
	})

	suite.Run("Should be able to use this service to update a shipment with SIT", func() {
		move := setupTestData()
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
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
					Code: models.ReServiceCodeDOPSIT,
				},
			},
		}, nil)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
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
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.Error(err)
		suite.Nil(update)
	})
	suite.Run("Request destination address changes on the same shipment multiple times", func() {
		move := setupTestData()
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), nil, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     shipment.DestinationAddress.PostalCode,
			Country:        models.StringPointer("United States"),
		}
		mockPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"90210",
			"94535",
		).Return(2500, nil).Times(4)
		update, err := addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address", etag.GenerateEtag(shipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("we really need to change the address", update.ContractorRemarks)

		// Need to re-request the shipment to get the updated etag
		var updatedShipment models.MTOShipment
		err = suite.DB().Find(&updatedShipment, shipment.ID)
		suite.NoError(err)

		update, err = addressUpdateRequester.RequestShipmentDeliveryAddressUpdate(suite.AppContextForTest(), shipment.ID, newAddress, "we really need to change the address again", etag.GenerateEtag(updatedShipment.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
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
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("United States"),
		}
		move := setupTestData()
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
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "89503",
			Country:        models.StringPointer("United States"),
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
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "89503",
			Country:        models.StringPointer("United States"),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentWithMove(&move, suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "94535",
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
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		newAddress := models.Address{
			StreetAddress1: "123 Any St",
			City:           "Albuquerque",
			State:          "NM",
			PostalCode:     "87053",
			Country:        models.StringPointer("United States"),
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
}

func (suite *ShipmentAddressUpdateServiceSuite) TestTOOApprovedShipmentAddressUpdateRequest() {

	addressCreator := address.NewAddressCreator()
	mockPlanner := &routemocks.Planner{}
	moveRouter := moveservices.NewMoveRouter()
	addressUpdateRequester := NewShipmentAddressUpdateRequester(mockPlanner, addressCreator, moveRouter)

	suite.Run("TOO approves address change", func() {

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)
		officeRemarks := "This is a TOO remark"

		update, err := addressUpdateRequester.ReviewShipmentAddressChange(suite.AppContextForTest(), addressChange.Shipment.ID, "APPROVED", officeRemarks)

		suite.NoError(err)
		suite.NotNil(update)
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, update.Status)
		suite.Equal("This is a TOO remark", *update.OfficeRemarks)

	})

	suite.Run("TOO rejects address change", func() {

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)
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
		fmt.Println(err)
		suite.Nil(update)
	})

}
