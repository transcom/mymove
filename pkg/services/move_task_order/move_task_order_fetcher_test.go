package movetaskorder_test

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	m "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderFetcher() {
	waf := entitlements.NewWeightAllotmentFetcher()

	setupTestData := func() (models.Move, models.MTOShipment) {

		expectedMTO := factory.BuildMove(suite.DB(), nil, nil)

		// Make a couple of shipments for the move; one prime, one external
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DeletedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)

		return expectedMTO, primeShipment
	}

	mtoFetcher := m.NewMoveTaskOrderFetcher(waf)

	suite.Run("Success with fetching a MTO that has a shipment address update", func() {
		traits := []factory.Trait{factory.GetTraitShipmentAddressUpdateApproved}
		expectedAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), nil, traits)

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:      false,
			IsAvailableToPrime: true,
			MoveTaskOrderID:    expectedAddressUpdate.Shipment.MoveTaskOrder.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// Validate MTO was fetched that includes expected shipment address update
		actualAddressUpdate := actualMTO.MTOShipments[0].DeliveryAddressUpdate
		suite.Equal(expectedAddressUpdate.ShipmentID, actualAddressUpdate.ShipmentID)
		suite.Equal(expectedAddressUpdate.Status, actualAddressUpdate.Status)
		suite.Equal(expectedAddressUpdate.OfficeRemarks, actualAddressUpdate.OfficeRemarks)
		suite.Equal(expectedAddressUpdate.ContractorRemarks, actualAddressUpdate.ContractorRemarks)
		suite.Equal(expectedAddressUpdate.NewAddressID, actualAddressUpdate.NewAddressID)
		suite.Equal(expectedAddressUpdate.OriginalAddressID, actualAddressUpdate.OriginalAddressID)
	})

	suite.Run("Success with fetching a MTO with a Shipment Address Update that has a customized Original Address", func() {
		traits := []factory.Trait{factory.GetTraitShipmentAddressUpdateApproved}

		expectedAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
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
		}, traits)

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:      false,
			IsAvailableToPrime: true,
			MoveTaskOrderID:    expectedAddressUpdate.Shipment.MoveTaskOrder.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		actualAddressUpdate := actualMTO.MTOShipments[0].DeliveryAddressUpdate

		// Validate MTO was fetched that includes expected shipment address update with customized original address
		suite.Equal(expectedAddressUpdate.ShipmentID, actualAddressUpdate.ShipmentID)
		suite.Equal(expectedAddressUpdate.Status, actualAddressUpdate.Status)
		suite.ElementsMatch(expectedAddressUpdate.OriginalAddressID, actualAddressUpdate.OriginalAddressID)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress1, actualAddressUpdate.OriginalAddress.StreetAddress1)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress2, actualAddressUpdate.OriginalAddress.StreetAddress2)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress3, actualAddressUpdate.OriginalAddress.StreetAddress3)
		suite.Equal(expectedAddressUpdate.OriginalAddress.City, actualAddressUpdate.OriginalAddress.City)
		suite.Equal(expectedAddressUpdate.OriginalAddress.State, actualAddressUpdate.OriginalAddress.State)
		suite.Equal(expectedAddressUpdate.OriginalAddress.PostalCode, actualAddressUpdate.OriginalAddress.PostalCode)
		suite.Equal(expectedAddressUpdate.OriginalAddress.CountryId, actualAddressUpdate.OriginalAddress.CountryId)
	})

	suite.Run("Success with fetching a MTO with a Shipment Address Update that has a customized Original Address and three addresses", func() {
		traits := []factory.Trait{factory.GetTraitShipmentAddressUpdateApproved}

		expectedAddressUpdate := factory.BuildShipmentAddressUpdate(suite.DB(), []factory.Customization{
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
		}, traits)

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:      false,
			IsAvailableToPrime: true,
			MoveTaskOrderID:    expectedAddressUpdate.Shipment.MoveTaskOrder.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		actualAddressUpdate := actualMTO.MTOShipments[0].DeliveryAddressUpdate

		// Validate MTO was fetched that includes expected shipment address update with customized original address
		suite.Equal(expectedAddressUpdate.ShipmentID, actualAddressUpdate.ShipmentID)
		suite.Equal(expectedAddressUpdate.Status, actualAddressUpdate.Status)
		suite.ElementsMatch(expectedAddressUpdate.OriginalAddressID, actualAddressUpdate.OriginalAddressID)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress1, actualAddressUpdate.OriginalAddress.StreetAddress1)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress2, actualAddressUpdate.OriginalAddress.StreetAddress2)
		suite.Equal(expectedAddressUpdate.OriginalAddress.StreetAddress3, actualAddressUpdate.OriginalAddress.StreetAddress3)
		suite.Equal(expectedAddressUpdate.OriginalAddress.City, actualAddressUpdate.OriginalAddress.City)
		suite.Equal(expectedAddressUpdate.OriginalAddress.State, actualAddressUpdate.OriginalAddress.State)
		suite.Equal(expectedAddressUpdate.OriginalAddress.PostalCode, actualAddressUpdate.OriginalAddress.PostalCode)
		suite.Equal(expectedAddressUpdate.OriginalAddress.CountryId, actualAddressUpdate.OriginalAddress.CountryId)
	})

	suite.Run("Success with Prime-available move by ID, fetch all non-deleted shipments", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.Nil(expectedMTO.ApprovedAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get two shipments back since we didn't set searchParams to exclude external ones.
		suite.Len(actualMTO.MTOShipments, 2)
	})

	suite.Run("Success with fetching contractor portion of a move", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotNil(expectedMTO.Contractor, actualMTO.Contractor)
		suite.NotNil(expectedMTO.Contractor.ContractNumber, actualMTO.Contractor.ContractNumber)
	})

	suite.Run("Success with fetching move with a related service item", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		address := factory.BuildAddress(suite.DB(), nil, nil)
		sitEntryDate := time.Now()
		customerContact := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{})
		serviceItemBasic := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &sitEntryDate,
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact},
				},
			},
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT, // DDFSIT - Domestic destination 1st day SIT
				},
			},
		}, nil)
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeICRT,
				},
			},
		}, nil)
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIUCRT,
				},
			},
		}, nil)
		serviceRequestDocumentUpload := factory.BuildServiceRequestDocumentUpload(suite.DB(), []factory.Customization{
			{
				Model:    serviceItemBasic,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
		}, nil)

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		found := false
		for _, serviceItem := range actualMTO.MTOServiceItems {
			if serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
				suite.Equal(address.StreetAddress1, serviceItem.SITDestinationFinalAddress.StreetAddress1)
				suite.Equal(address.State, serviceItem.SITDestinationFinalAddress.State)
				suite.Equal(address.City, serviceItem.SITDestinationFinalAddress.City)
				suite.Equal(1, len(serviceItem.CustomerContacts))

				if suite.Len(serviceItem.ServiceRequestDocuments, 1) {
					if suite.Len(serviceItem.ServiceRequestDocuments[0].ServiceRequestDocumentUploads, 1) {
						suite.Equal(serviceRequestDocumentUpload.ID, serviceItem.ServiceRequestDocuments[0].ServiceRequestDocumentUploads[0].ID)
					}
				}

				found = true
				break
			}
			if serviceItem.ReService.Code == models.ReServiceCodeICRT || serviceItem.ReService.Code == models.ReServiceCodeIUCRT {
				suite.NotNil(serviceItem.MTOShipment.PickupAddress)
				suite.NotNil(serviceItem.MTOShipment.DestinationAddress)
				break
			}
		}
		// Verify that the expected service item was found
		suite.True(found, "Expected service item with ReServiceCodeDDFSIT not found")
	})

	suite.Run("Success with Prime-available move by Locator, no deleted or external shipments", func() {
		expectedMTO, primeShipment := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			Locator:                  expectedMTO.Locator,
			ExcludeExternalShipments: true,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.Nil(expectedMTO.ApprovedAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		// Should get one shipment back since we requested no external shipments.
		if suite.Len(actualMTO.MTOShipments, 1) {
			suite.Equal(expectedMTO.ID.String(), actualMTO.ID.String())
			suite.Equal(primeShipment.ID.String(), actualMTO.MTOShipments[0].ID.String())
		}
	})

	suite.Run("Success with Prime available move, returns destination GBLOC in shipment dest address", func() {
		zone2UUID, err := uuid.FromString("66768964-e0de-41f3-b9be-7ef32e4ae2b4")
		suite.FatalNoError(err)
		army := models.AffiliationARMY
		postalCode := "99501"
		// since we truncate the test db, we need to add the postal_code_to_gbloc value
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "99744", "JEAT")

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode:         postalCode,
					UsPostRegionCityID: &zone2UUID,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			Locator:                  move.Locator,
			ExcludeExternalShipments: true,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(actualMTO)

		if suite.Len(actualMTO.MTOShipments, 1) {
			suite.Equal(move.ID.String(), actualMTO.ID.String())
			// the shipment should have a destination GBLOC value
			suite.NotNil(actualMTO.MTOShipments[0].DestinationAddress.DestinationGbloc)
		}
	})

	suite.Run("Success with Prime available move, returns USMC destination GBLOC for USMC move", func() {
		usmc := models.AffiliationMARINES

		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &usmc,
				},
			},
		}, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			Locator:                  move.Locator,
			ExcludeExternalShipments: true,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(actualMTO)

		if suite.Len(actualMTO.MTOShipments, 1) {
			suite.Equal(move.ID.String(), actualMTO.ID.String())
			suite.NotNil(actualMTO.MTOShipments[0].DestinationAddress.DestinationGbloc)
			suite.Equal(*actualMTO.MTOShipments[0].DestinationAddress.DestinationGbloc, "USMC")
		}
	})

	suite.Run("Success with move that has only deleted shipments", func() {
		mtoWithAllShipmentsDeleted := factory.BuildMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mtoWithAllShipmentsDeleted,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					DeletedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: mtoWithAllShipmentsDeleted.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Equal(mtoWithAllShipmentsDeleted.ID, actualMTO.ID)
		suite.Len(actualMTO.MTOShipments, 0)
	})

	suite.Run("Failure - nil searchParams", func() {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), nil)
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should not be nil")
	})

	suite.Run("Failure - searchParams with no ID/locator set", func() {
		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &services.MoveTaskOrderFetcherParams{})
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should have either a move ID or locator set")
	})

	suite.Run("Failure - Not Found with Bad ID", func() {
		badID, _ := uuid.NewV4()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: badID,
		}

		_, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.Error(err)
	})

	suite.Run("Success - found move with Port of Embarkation", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		poeId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					POELocationID: &poeId,
					Status:        models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC, // POEFSC, International POE Fuel Surcharge
				},
			},
		}, nil)

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		found := false
		for _, serviceItem := range actualMTO.MTOServiceItems {
			if serviceItem.ReService.Code == models.ReServiceCodePOEFSC {
				suite.Equal(poeId, serviceItem.POELocation.ID)
				found = true
				break
			}
		}
		// Verify that the expected service item was found
		suite.True(found, "Expected service item ReServiceCodePOEFSC")
	})

	suite.Run("Success - found move with Port of Debarkation", func() {
		expectedMTO, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: expectedMTO.ID,
		}

		podId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					PODLocationID: &podId,
					Status:        models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    expectedMTO,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePODFSC, // PODFSC, International POD Fuel Surcharge
				},
			},
		}, nil)

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		found := false
		for _, serviceItem := range actualMTO.MTOServiceItems {
			if serviceItem.ReService.Code == models.ReServiceCodePODFSC {
				suite.Equal(podId, serviceItem.PODLocation.ID)
				found = true
				break
			}
		}
		// Verify that the expected service item was found
		suite.True(found, "Expected service item ReServiceCodePODFSC")
	})

	suite.Run("Success - Move contains only deleted MTOAgents", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{Model: models.MTOAgent{
				MTOAgentType: models.MTOAgentReceiving,
				DeletedAt:    models.TimePointer(time.Now()),
			}},
			{Model: move, LinkOnly: true},
		}, nil)

		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: move.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.Equal(move.ID, actualMTO.ID)
		suite.Len(actualMTO.MTOShipments[0].MTOAgents, 0, "Expected no active agents since all are deleted")
	})

	suite.Run("Success - Move contains one MTOAgent", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{Model: move,
				LinkOnly: true},
		}, nil)

		factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{Model: models.MTOAgent{
				MTOAgentType: models.MTOAgentReceiving,
				DeletedAt:    models.TimePointer(time.Now()),
			}},
			{Model: move, LinkOnly: true},
			{Model: shipment, LinkOnly: true},
		}, nil)

		activeAgent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{Model: models.MTOAgent{
				FirstName:    models.StringPointer("John"),
				LastName:     models.StringPointer("Doe"),
				Email:        models.StringPointer("John.doe@example.com"),
				Phone:        models.StringPointer("222-222-2222"),
				MTOAgentType: models.MTOAgentReleasing,
			}},
			{Model: move, LinkOnly: true},
			{Model: shipment, LinkOnly: true},
		}, nil)

		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: move.ID,
		}

		actualMTO, err := mtoFetcher.FetchMoveTaskOrder(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.Equal(move.ID, actualMTO.ID)
		suite.Len(actualMTO.MTOShipments[0].MTOAgents, 1, "Expected only one active agent in the result")

		activeAgentReturned := actualMTO.MTOShipments[0].MTOAgents[0]
		suite.Equal(activeAgent.FirstName, activeAgentReturned.FirstName, "First names should match")
		suite.Equal(activeAgent.LastName, activeAgentReturned.LastName, "Last names should match")
		suite.Equal(activeAgent.Email, activeAgentReturned.Email, "Emails should match")
		suite.Equal(activeAgent.Phone, activeAgentReturned.Phone, "Phone numbers should match")
		suite.Equal(activeAgent.MTOAgentType, activeAgentReturned.MTOAgentType, "Agent types should match")
	})

}

func (suite *MoveTaskOrderServiceSuite) TestGetMoveTaskOrderFetcher() {
	setupTestData := func() models.Move {

		expectedMTO := factory.BuildMove(suite.DB(), nil, nil)

		return expectedMTO
	}
	waf := entitlements.NewWeightAllotmentFetcher()

	mtoFetcher := m.NewMoveTaskOrderFetcher(waf)

	suite.Run("success getting a move using GetMove for Prime user", func() {
		expectedMTO := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: expectedMTO.ID,
		}

		move, err := mtoFetcher.GetMove(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		suite.Equal(expectedMTO.ID, move.ID)
	})

	suite.Run("get an error if search params are not provided when using GetMove", func() {
		_, err := mtoFetcher.GetMove(suite.AppContextForTest(), &services.MoveTaskOrderFetcherParams{})
		suite.Error(err)
		suite.Contains(err.Error(), "searchParams should have either a move ID or locator set")
	})

	suite.Run("get an error if bad ID is provided when using GetMove", func() {
		badID, _ := uuid.NewV4()
		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: badID,
		}

		_, err := mtoFetcher.GetMove(suite.AppContextForTest(), &searchParams)
		suite.Error(err)
		suite.Contains(err.Error(), "not found")
	})

	suite.Run("Can fetch a move if it is a customer app request by the customer it belongs to", func() {
		expectedMTO := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		serviceMember := expectedMTO.Orders.ServiceMember
		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: expectedMTO.ID,
		}

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          serviceMember.User.ID,
			ServiceMemberID: serviceMember.ID,
		})

		moveReturned, err := mtoFetcher.GetMove(
			appCtx,
			&searchParams,
		)

		if suite.NoError(err) && suite.NotNil(moveReturned) {
			suite.Equal(expectedMTO.ID, moveReturned.ID)
		}
	})

	suite.Run("Returns a not found error if it is a customer app request by a customer that it does not belong to", func() {
		badUser := factory.BuildExtendedServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          badUser.User.ID,
			ServiceMemberID: badUser.ID,
		})

		expectedMTO := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: expectedMTO.ID,
		}

		moveReturned, err := mtoFetcher.GetMove(
			appCtx,
			&searchParams,
		)

		if suite.Error(err) && suite.Nil(moveReturned) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Contains(err.Error(), fmt.Sprintf("ID: %s not found", expectedMTO.ID))
		}
	})

	suite.Run("success getting a move for Office user", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), factory.GetTraitActiveOfficeUser(), nil)
		expectedMTO := factory.BuildMove(suite.DB(), nil, nil)

		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: expectedMTO.ID,
		}

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          officeUser.User.ID,
			OfficeUserID:    officeUser.ID,
		})

		moveReturned, err := mtoFetcher.GetMove(
			appCtx,
			&searchParams,
		)

		if suite.NoError(err) && suite.NotNil(moveReturned) {
			suite.Equal(expectedMTO.ID, moveReturned.ID)
		}
	})
}

// Checks that there are expectedMatchCount matches between the moves and move ID list
// Allows caller to check that expected moves are in the list
// Also returns a map from uuids → moves
func doIDsMatch(moves []models.Move, moveIDs []*uuid.UUID, expectedMatchCount int) (bool, map[uuid.UUID]models.Move) {
	matched := 0
	mapIDs := make(map[uuid.UUID]models.Move)
	for _, expectedID := range moveIDs {
		for _, actualMove := range moves {
			if expectedID != nil && actualMove.ID == *expectedID {
				mapIDs[*expectedID] = actualMove
				expectedID = nil // if found, clear so we don't match again
				matched++        // keep count so we match all
			}
		}
	}
	return expectedMatchCount == matched, mapIDs
}

func (suite *MoveTaskOrderServiceSuite) TestListAllMoveTaskOrdersFetcher() {
	// Set up a hidden move so we can check if it's in the output:
	now := time.Now()
	show := false
	waf := entitlements.NewWeightAllotmentFetcher()

	setupTestData := func() (models.Move, models.Move, models.MTOShipment) {
		hiddenMTO := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &show,
				},
			},
		}, nil)

		mto := factory.BuildMove(suite.DB(), nil, nil)

		// Make a couple of shipments for the default move; one prime, one external
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		return hiddenMTO, mto, primeShipment
	}

	mtoFetcher := m.NewMoveTaskOrderFetcher(waf)

	suite.Run("all move task orders", func() {
		hiddenMTO, mto, _ := setupTestData()
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: false,
			IncludeHidden:      true,
			Since:              nil,
		}

		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		move := moves[0]

		suite.Equal(2, len(moves))
		// These are the move ids we want to find in the return list
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &mto.ID}
		suite.True(doIDsMatch(moves, moveIDs, 2))
		suite.NotNil(move.Orders)
		suite.NotNil(move.Orders.OriginDutyLocation)
		suite.NotNil(move.Orders.NewDutyLocation)
	})

	suite.Run("default search - excludes hidden move task orders", func() {
		hiddenMTO, expectedMove, _ := setupTestData()
		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), nil)
		suite.NoError(err)

		// The hidden move should be nowhere in the output list:
		for _, move := range moves {
			suite.NotEqual(move.ID, hiddenMTO.ID)
		}

		// These are the move ids we have but we expect to find just one
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &expectedMove.ID}
		result, mapMoves := doIDsMatch(moves, moveIDs, 1)
		suite.True(result)                                   // exactly one ID should match
		suite.Contains(mapMoves, expectedMove.ID)            // it should be the expectedMove
		suite.Len(mapMoves[expectedMove.ID].MTOShipments, 2) // That move should have 2 shipments

	})

	suite.Run("returns shipments that respect the external searchParams flag", func() {
		hiddenMTO, expectedMove, primeShipment := setupTestData()

		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:            false,
			ExcludeExternalShipments: true,
		}

		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// We should only get back the one move that's not hidden.
		// These are the move ids we have but we expect to find just one
		moveIDs := []*uuid.UUID{&hiddenMTO.ID, &expectedMove.ID}
		result, mapMoves := doIDsMatch(moves, moveIDs, 1)
		suite.True(result)                                   // exactly one ID should match
		suite.Len(mapMoves[expectedMove.ID].MTOShipments, 1) // That move should have one shipment
		suite.Equal(primeShipment.ID.String(), mapMoves[expectedMove.ID].MTOShipments[0].ID.String())
	})

	suite.Run("all move task orders that are available to prime and using since", func() {
		// Under test: ListAllMoveOrders
		// Mocked:     None
		// Set up:     Create an old and new move, first check that both are returned.
		//             Then make one move "old" by moving the updatedAt time back
		//             Request all move orders again
		// Expected outcome: Only the new move should be returned
		now = time.Now()

		// Create the two moves
		newMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		oldMTO := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Check that they're both returned
		searchParams := services.MoveTaskOrderFetcherParams{
			IsAvailableToPrime: true,
			// IncludeHidden should be false by default
			Since: nil,
		}
		moves, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(moves, 2)

		// Put 1 Move updatedAt in the past
		suite.NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
			now.Add(-2*time.Hour), oldMTO.ID).Exec())

		// Make search params search for moves newer than timestamp
		searchParams.Since = &now
		mtosWithSince, err := mtoFetcher.ListAllMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)

		// Expect one move, the newer one
		suite.Equal(1, len(mtosWithSince))
		moveIDs := []*uuid.UUID{&oldMTO.ID, &newMove.ID}
		result, mapMoves2 := doIDsMatch(mtosWithSince, moveIDs, 1) // only one should match
		suite.True(result)
		suite.Contains(mapMoves2, newMove.ID)                                                                            // and it should be the new one
		suite.NotContains(mapMoves2, oldMTO.ID, "Returned moves should not contain the old move %s", oldMTO.ID.String()) // it's not the hiddenMTO
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListPrimeMoveTaskOrdersFetcher() {
	now := time.Now()
	waf := entitlements.NewWeightAllotmentFetcher()
	// Set up a hidden move so we can check if it's in the output:
	hiddenMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Show: models.BoolPointer(false),
			},
		},
	}, nil)
	// Make a default, not Prime-available move:
	nonPrimeMove := factory.BuildMove(suite.DB(), nil, nil)
	// Make some Prime moves:
	primeMove1 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	primeMove2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	primeMove3 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	factory.BuildMTOShipmentWithMove(&primeMove3, suite.DB(), nil, nil)
	primeMove4 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	shipmentForPrimeMove4 := factory.BuildMTOShipmentWithMove(&primeMove4, suite.DB(), nil, nil)
	reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
		MTOShipment: shipmentForPrimeMove4,
	})
	suite.Logger().Info(fmt.Sprintf("Reweigh %s", reweigh.ID))
	// Move primeMove1, primeMove3, and primeMove4 into the past so we can exclude them:
	suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=$1 WHERE id IN ($2, $3, $4);",
		now.Add(-10*time.Second), primeMove1.ID, primeMove3.ID, primeMove4.ID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id IN ($2, $3);",
		now.Add(-10*time.Second), primeMove1.OrdersID, primeMove4.OrdersID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE mto_shipments SET updated_at=$1 WHERE id=$2;",
		now.Add(-10*time.Second), shipmentForPrimeMove4.ID).Exec())

	fetcher := m.NewMoveTaskOrderFetcher(waf)
	page := int64(1)
	perPage := int64(20)
	// filling out search params to allow for pagination
	searchParams := services.MoveTaskOrderFetcherParams{Page: &page, PerPage: &perPage, MoveCode: nil, ID: nil}

	// Run the fetcher without `since` to get all Prime moves:
	primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(primeMoves, 4)

	moveIDs := []uuid.UUID{primeMoves[0].ID, primeMoves[1].ID, primeMoves[2].ID, primeMoves[3].ID}
	suite.NotContains(moveIDs, hiddenMove.ID)
	suite.NotContains(moveIDs, nonPrimeMove.ID)
	suite.Contains(moveIDs, primeMove1.ID)
	suite.Contains(moveIDs, primeMove2.ID)
	suite.Contains(moveIDs, primeMove3.ID)
	suite.Contains(moveIDs, primeMove4.ID)

	// Run the fetcher with `since` to get primeMove2, primeMove3 (because of the shipment), and primeMove4 (because of the reweigh)
	since := now.Add(-5 * time.Second)
	searchParams.Since = &since
	sinceMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
	suite.NoError(err)
	suite.Len(sinceMoves, 3)

	sinceMoveIDs := []uuid.UUID{sinceMoves[0].ID, sinceMoves[1].ID, sinceMoves[2].ID}
	suite.Contains(sinceMoveIDs, primeMove2.ID)
	suite.Contains(sinceMoveIDs, primeMove3.ID)
	suite.Contains(sinceMoveIDs, primeMove4.ID)
}

func (suite *MoveTaskOrderServiceSuite) TestListPrimeMoveTaskOrdersAmendmentsFetcher() {
	waf := entitlements.NewWeightAllotmentFetcher()

	suite.Run("Test with and without filter of moves containing amendments", func() {
		now := time.Now()
		// Set up a hidden move so we can check if it's in the output:
		hiddenMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: models.BoolPointer(false),
				},
			},
		}, nil)
		// Make a default, not Prime-available move:
		nonPrimeMove := factory.BuildMove(suite.DB(), nil, nil)
		// Make some Prime moves:
		primeMove1 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		primeMove2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		primeMove3 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		factory.BuildMTOShipmentWithMove(&primeMove3, suite.DB(), nil, nil)
		primeMove4 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentForPrimeMove4 := factory.BuildMTOShipmentWithMove(&primeMove4, suite.DB(), nil, nil)
		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipmentForPrimeMove4,
		})
		suite.Logger().Info(fmt.Sprintf("Reweigh %s", reweigh.ID))

		primeMove5 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		//////////////////////////////////////////
		// setup amendments for move
		//////////////////////////////////////////
		primeMoves := make([]models.Move, 0)
		primeMoves = append(primeMoves, primeMove1)
		primeMoves = append(primeMoves, primeMove2)
		primeMoves = append(primeMoves, primeMove3)
		primeMoves = append(primeMoves, primeMove4)
		primeMoves = append(primeMoves, primeMove5)

		docIDs := make([]uuid.UUID, 0)
		hasAmendmentsMap := make(map[uuid.UUID]bool)
		for i, pm := range primeMoves {
			document := factory.BuildDocumentLinkServiceMember(suite.DB(), primeMove1.Orders.ServiceMember)

			docIDs = append(docIDs, document.ID)

			suite.MustSave(&document)
			suite.Nil(document.DeletedAt)
			pm.Orders.UploadedOrders = document
			pm.Orders.UploadedOrdersID = document.ID

			if i != 4 {
				// set amendment for all except for one.
				pm.Orders.UploadedAmendedOrders = &document
				pm.Orders.UploadedAmendedOrdersID = &document.ID
				hasAmendmentsMap[pm.ID] = true
			} else {
				hasAmendmentsMap[pm.ID] = false
			}
			//nolint:gosec //G601
			suite.MustSave(&pm.Orders)
			upload := models.Upload{
				Filename:    "test.pdf",
				Bytes:       1048576,
				ContentType: uploader.FileTypePDF,
				Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
				UploadType:  models.UploadTypeUSER,
			}
			suite.MustSave(&upload)
			userUpload := models.UserUpload{
				DocumentID: &document.ID,
				UploaderID: document.ServiceMember.UserID,
				UploadID:   upload.ID,
				Upload:     upload,
			}
			suite.MustSave(&userUpload)
		}

		// Move primeMove1, primeMove3, and primeMove4 into the past so we can exclude them:
		suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=$1 WHERE id IN ($2, $3, $4);",
			now.Add(-10*time.Second), primeMove1.ID, primeMove3.ID, primeMove4.ID).Exec())
		suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id IN ($2, $3);",
			now.Add(-10*time.Second), primeMove1.OrdersID, primeMove4.OrdersID).Exec())
		suite.Require().NoError(suite.DB().RawQuery("UPDATE mto_shipments SET updated_at=$1 WHERE id=$2;",
			now.Add(-10*time.Second), shipmentForPrimeMove4.ID).Exec())

		fetcher := m.NewMoveTaskOrderFetcher(waf)
		page := int64(1)
		perPage := int64(20)
		// filling out search params to allow for pagination
		searchParams := services.MoveTaskOrderFetcherParams{Page: &page, PerPage: &perPage, MoveCode: nil, ID: nil}

		// Run the fetcher without `since` to get all Prime moves:
		primeMoves, amendmentCountInfo, err := fetcher.ListPrimeMoveTaskOrdersAmendments(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 5)

		moveIDs := []uuid.UUID{primeMoves[0].ID, primeMoves[1].ID, primeMoves[2].ID, primeMoves[3].ID, primeMoves[4].ID}
		suite.NotContains(moveIDs, hiddenMove.ID)
		suite.NotContains(moveIDs, nonPrimeMove.ID)
		suite.Contains(moveIDs, primeMove1.ID)
		suite.Contains(moveIDs, primeMove2.ID)
		suite.Contains(moveIDs, primeMove3.ID)
		suite.Contains(moveIDs, primeMove4.ID)
		suite.Contains(moveIDs, primeMove5.ID)

		// amendmentCountInfo should only contain moves that have amendments.
		suite.Len(amendmentCountInfo, 4)

		cnt := 0
		for _, value := range amendmentCountInfo {
			if hasAmendmentsMap[value.MoveID] {
				suite.Equal(1, value.Total)
				suite.Equal(1, value.AvailableSinceTotal)
				cnt++
			}
			// verify the Prime Moves without any amendments are NOT
			// in amendmentCountInfo
			for moveID, hasAmendment := range hasAmendmentsMap {
				if !hasAmendment {
					suite.False(value.MoveID == moveID)
				}
			}
		}
		suite.Equal(len(amendmentCountInfo), cnt)

		// Run the fetcher with `since` to get primeMove2, primeMove3 (because of the shipment), and primeMove4 (because of the reweigh)
		since := now.Add(-5 * time.Second)
		searchParams.Since = &since

		// fake out timestamp for new amendment upload by manually setting update column
		suite.Require().NoError(suite.DB().RawQuery("UPDATE user_uploads SET updated_at=$1 WHERE document_id IN ($2, $3, $4, $5, $6);",
			now.Add(-100*time.Second), docIDs[0], docIDs[1], docIDs[2], docIDs[3], docIDs[4]).Exec())

		sinceMoves, amendmentCountInfo, err := fetcher.ListPrimeMoveTaskOrdersAmendments(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(sinceMoves, 4)
		suite.Len(amendmentCountInfo, 3)

		sinceMoveIDs := []uuid.UUID{sinceMoves[0].ID, sinceMoves[1].ID, sinceMoves[2].ID, sinceMoves[3].ID}
		suite.Contains(sinceMoveIDs, primeMove2.ID)
		suite.Contains(sinceMoveIDs, primeMove3.ID)
		suite.Contains(sinceMoveIDs, primeMove4.ID)

		for _, value := range amendmentCountInfo {
			if hasAmendmentsMap[value.MoveID] {
				suite.Equal(1, value.Total)
				// verify sinceCount is filtering based on since parameter. Amendment was uploaded at an older date.
				suite.Equal(0, value.AvailableSinceTotal)
			}
		}
	})

	suite.Run("Test moves without any amendments", func() {
		now := time.Now()
		primeMove1 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		primeMove2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Move primeMove1, primeMove2 into the past so we can exclude them:
		suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=$1 WHERE id IN ($2, $3);",
			now.Add(-10*time.Second), primeMove1.ID, primeMove2.ID).Exec())

		fetcher := m.NewMoveTaskOrderFetcher(waf)
		page := int64(1)
		perPage := int64(20)
		// filling out search params to allow for pagination
		searchParams := services.MoveTaskOrderFetcherParams{Page: &page, PerPage: &perPage, MoveCode: nil, ID: nil}

		// Run the fetcher without `since` to get all Prime moves:
		primeMoves, amendmentCountInfo, err := fetcher.ListPrimeMoveTaskOrdersAmendments(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 2)
		suite.Len(amendmentCountInfo, 0)
		moveIDs := []uuid.UUID{primeMoves[0].ID, primeMoves[1].ID}
		suite.Contains(moveIDs, primeMove1.ID)
		suite.Contains(moveIDs, primeMove2.ID)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestListPrimeMoveTaskOrdersAcknowledgements() {
	waf := entitlements.NewWeightAllotmentFetcher()
	fetcher := m.NewMoveTaskOrderFetcher(waf)
	falseValue := false
	trueValue := true
	suite.Run("Retrieves only moves with prime acknowledgement", func() {
		now := time.Now()
		// Prime Move that has NOT been acknowledged
		factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: nil,
				},
			},
		}, nil)

		//Prime move that has already been acknowledged
		acknowledgedPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &now,
				},
			}}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			Acknowledged: &trueValue,
		}
		primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(primeMoves)
		suite.Len(primeMoves, 1, "Only the move that has already been acknowledged should be present.")
		suite.Equal(acknowledgedPrimeMove.ID, primeMoves[0].ID)
		suite.Equal(acknowledgedPrimeMove.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), primeMoves[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))

		//	Add an unacknowledged shipment to the acknowledged move
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    acknowledgedPrimeMove,
				LinkOnly: true,
			}}, nil)
		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Empty(primeMoves, "move no longer included since an unacknowledged shipment was added to it")
	})

	suite.Run("Retrieves moves without prime acknowledgement", func() {
		now := time.Now()
		// Prime Move that has NOT been acknowledged
		unacknowledgedPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: nil,
				},
			},
		}, nil)

		//Prime move that has already been acknowledged
		factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &now,
				},
			}}, nil)
		searchParams := services.MoveTaskOrderFetcherParams{
			Acknowledged: &falseValue,
		}
		primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(primeMoves)
		suite.Len(primeMoves, 1, "Only the move that has NOT been acknowledged should be present.")
		suite.Equal(unacknowledgedPrimeMove.ID, primeMoves[0].ID)
		suite.Nil(primeMoves[0].PrimeAcknowledgedAt)

		// Add prime acknowledgement to the move
		unacknowledgedPrimeMove.PrimeAcknowledgedAt = &now
		err = suite.DB().Save(&unacknowledgedPrimeMove)
		suite.NoError(err)
		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Empty(primeMoves, "move no longer included since it has been acknowledged")
	})

	suite.Run("Retrieves moves where prime acknowledgement on the move or any one of its shipments occurred after the acknowledegAfter date", func() {
		aWeekAgo := time.Now().AddDate(0, 0, -7)
		yesterday := time.Now().AddDate(0, 0, -1)
		today := time.Now()

		// Prime Move that was acknowledged today
		moveAcknowledgedToday := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &today,
				},
			},
		}, nil)

		// Prime Move that was acknowledged a week ago
		moveAcknowledgedAWeekAgo := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &aWeekAgo,
				},
			}}, nil)

		//Fetch moves acknowledged after yesterday
		searchParams := services.MoveTaskOrderFetcherParams{
			AcknowledgedAfter: &yesterday,
		}
		primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(primeMoves)
		suite.Len(primeMoves, 1, "Only the move that was acknowledged after yesterday should be present.")
		suite.Equal(moveAcknowledgedToday.ID, primeMoves[0].ID)
		suite.Equal(moveAcknowledgedToday.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), primeMoves[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))

		//Add a shipment that was acknowledged today to the older move
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeAcknowledgedAt: &today,
				},
			},
			{
				Model:    moveAcknowledgedAWeekAgo,
				LinkOnly: true,
			}}, nil)
		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 2, "Both moves now included since a recently acknowledged shipment was added to the older move")
	})

	suite.Run("Retrieves moves where prime acknowledgement on the move or any one of its shipments occurred before the acknowledegBefore date", func() {
		aWeekAgo := time.Now().AddDate(0, 0, -7)
		yesterday := time.Now().AddDate(0, 0, -1)
		today := time.Now()

		// Prime Move that was acknowledged today
		moveAcknowledgedToday := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &today,
				},
			},
		}, nil)

		// Prime Move that was acknowledged a week ago
		moveAcknowledgedAWeekAgo := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &aWeekAgo,
				},
			}}, nil)

		//Fetch moves acknowledged after yesterday
		searchParams := services.MoveTaskOrderFetcherParams{
			AcknowledgedBefore: &yesterday,
		}
		primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.NotNil(primeMoves)
		suite.Len(primeMoves, 1, "Only the move that was acknowledged before yesterday should be present.")
		suite.Equal(moveAcknowledgedAWeekAgo.ID, primeMoves[0].ID)
		suite.Equal(moveAcknowledgedAWeekAgo.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), primeMoves[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))

		//Add a shipment that was acknowledged a week ago to the move that was acknowledged today
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeAcknowledgedAt: &aWeekAgo,
				},
			},
			{
				Model:    moveAcknowledgedToday,
				LinkOnly: true,
			}}, nil)
		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 2, "Both moves now included since a shipment acknowledged a week ago was added to the move")
	})

	suite.Run("Retrieve moves when multiple parameters are included", func() {
		now := time.Now()
		aWeekAgo := now.AddDate(0, 0, -7)
		yesterday := now.AddDate(0, 0, -1)
		tomorrow := now.AddDate(0, 0, 1)

		// Prime Move that was acknowledged today
		move1 := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &now,
				},
			},
		}, nil)

		// Prime Move that was acknowledged a week ago
		move2 := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &aWeekAgo,
				},
			}}, nil)

		// Prime Move that is not acknowledged
		move3 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		//Fetch moves acknowledged after yesterday and updated after tomorrow
		searchParams := services.MoveTaskOrderFetcherParams{
			Since:             &tomorrow,
			AcknowledgedAfter: &yesterday,
		}

		primeMoves, err := fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Empty(primeMoves, "None of the moves are updated after tomorrow")

		//Fetch moves acknowledged after yesterday and updated after yesterday
		searchParams = services.MoveTaskOrderFetcherParams{
			Since:             &yesterday,
			AcknowledgedAfter: &yesterday,
		}

		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 1, "Move 1 was acknowledged today and udpated today")
		suite.Equal(move1.ID, primeMoves[0].ID)

		//Fetch moves acknowledged after yesterday and updated after yesterday
		searchParams = services.MoveTaskOrderFetcherParams{
			Since:              &yesterday,
			AcknowledgedBefore: &yesterday,
		}

		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 1, "Move 2 was acknowledged a week ago and udpated today")
		suite.Equal(move2.ID, primeMoves[0].ID)

		//Fetch moves acknowledged after yesterday and updated after yesterday
		searchParams = services.MoveTaskOrderFetcherParams{
			Since:        &yesterday,
			Acknowledged: &falseValue,
		}

		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 1, "Move 3 is not acknowledged and was udpated today")
		suite.Equal(move3.ID, primeMoves[0].ID)

		//Fetch moves acknowledged after yesterday and updated after yesterday
		searchParams = services.MoveTaskOrderFetcherParams{
			Since:        &yesterday,
			Acknowledged: &trueValue,
		}

		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 2, "Move 1 and 2 are acknowledged and were udpated today")

		//Fetch moves acknowledged after yesterday and updated after yesterday
		searchParams = services.MoveTaskOrderFetcherParams{
			Since:              &yesterday,
			AcknowledgedBefore: &tomorrow,
			AcknowledgedAfter:  &yesterday,
		}

		primeMoves, err = fetcher.ListPrimeMoveTaskOrders(suite.AppContextForTest(), &searchParams)
		suite.NoError(err)
		suite.Len(primeMoves, 1, "Move 1 was acknowledged today and was udpated today")
		suite.Equal(move1.ID, primeMoves[0].ID)

	})
}
