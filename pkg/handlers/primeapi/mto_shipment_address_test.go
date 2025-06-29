package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	servicemocks "github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// isAddressEqual compares 2 addresses
func isAddressEqual(suite *HandlerSuite, reqAddress *primemessages.Address, respAddress *primemessages.Address) {
	if reqAddress.StreetAddress1 != nil && respAddress.StreetAddress1 != nil {
		suite.Equal(*reqAddress.StreetAddress1, *respAddress.StreetAddress1)
	}
	if reqAddress.StreetAddress2 != nil && respAddress.StreetAddress2 != nil {
		suite.Equal(*reqAddress.StreetAddress2, *respAddress.StreetAddress2)
	}
	if reqAddress.StreetAddress3 != nil && respAddress.StreetAddress3 != nil {
		suite.Equal(*reqAddress.StreetAddress3, *respAddress.StreetAddress3)
	}
	suite.Equal(*reqAddress.PostalCode, *respAddress.PostalCode)
	suite.Equal(*reqAddress.State, *respAddress.State)
	suite.Equal(*reqAddress.City, *respAddress.City)

}

func (suite *HandlerSuite) TestUpdateMTOShipmentAddressHandler() {
	setupTestData := func() (UpdateMTOShipmentAddressHandler, models.Move) {
		// Make an available MTO
		availableMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		planner := &mocks.Planner{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		vLocationServices := address.NewVLocation()
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)

		// Create handler
		handler := UpdateMTOShipmentAddressHandler{
			suite.NewHandlerConfig(),
			mtoshipment.NewMTOShipmentAddressUpdater(planner, addressCreator, addressUpdater),
			vLocationServices,
		}
		return handler, availableMove
	}

	newAddress := models.Address{
		StreetAddress1: "7 Q St",
		City:           "Acmar",
		State:          "AL",
		PostalCode:     "35004",
	}

	suite.Run("Success updating address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		// Make a shipment on the available MTO
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)

		// Update with new address
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)

		// Validate outgoing payload
		suite.NoError(shipmentOk.Payload.Validate(strfmt.Default))

		// Check values
		isAddressEqual(suite, payload, shipmentOk.Payload)
	})

	suite.Run("Success updating full address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO, all fields in address provided
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: models.StringPointer("6622 Airport Way S #1430"),
			StreetAddress3: models.StringPointer("441 SW Río de la Plata Drive"),
			City:           "Alameda",
			State:          "CA",
			PostalCode:     "94502",
		}

		// Update with new address
		payload := payloads.Address(&newAddress2)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)

		// Validate outgoing payload
		suite.NoError(shipmentOk.Payload.Validate(strfmt.Default))

		// Check values
		isAddressEqual(suite, payload, shipmentOk.Payload)

	})

	suite.Run("Fail - NotFound due to unavailable MTO", func() {
		// Testcase:   address is updated on a shipment that's on an MTO that is NOT available to Prime
		// Expected:   NotFound error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, _ := setupTestData()
		// Make a shipment with an unavailable MTO
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, nil)

		// Update with new address
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentAddressNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Fail - ConflictError due to unassociated mtoShipment", func() {
		// Testcase:   address is updated on a shipment that it's not associated with
		// Expected:   Conflict error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		// Make a random address that is not associated
		randomAddress := factory.BuildAddress(suite.DB(), nil, nil)

		payload := payloads.Address(&randomAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), randomAddress.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(randomAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(randomAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressConflict{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentAddressConflict).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Fail - PreconditionFailed due to wrong etag", func() {
		// Testcase:   address is updated on a shipment, but etag for address is wrong
		// Expected:   PreconditionFailed error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		// Update with new address with a bad etag
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       "bad-etag",
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressPreconditionFailed{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentAddressPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Fail - Unprocessable due to dest address being updated for approved shipment", func() {
		// Testcase:   destination address is updated on a shipment, but shipment is approved
		// Expected:   Conflict error is returned
		// Under Test: UpdateMTOShipmentAddress handler
		handler, availableMove := setupTestData()
		destAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    destAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		// Try to update destination address for approved shipment
		payload := payloads.Address(&destAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(destAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.DestinationAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		resp, ok := response.(*mtoshipmentops.UpdateMTOShipmentAddressConflict)
		suite.True(ok, "Expected response to be of type UpdateMTOShipmentAddressConflict")
		suite.Contains(*resp.Payload.Detail, "This shipment has already been approved, please use the updateShipmentDestinationAddress endpoint / ShipmentAddressUpdateRequester service to update the destination address")
	})

	suite.Run("Fail - Conflict due to updating pickup address on NTS-Release shipment", func() {
		// Testcase:   destination address is updated on a shipment, but shipment is approved
		// Expected:   Conflict error is returned
		// Under Test: UpdateMTOShipmentAddress handler
		handler, availableMove := setupTestData()
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		shipment := factory.BuildNTSRShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, nil)
		// Try to update destination address for approved shipment
		payload := payloads.Address(&pickupAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.DestinationAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		resp, ok := response.(*mtoshipmentops.UpdateMTOShipmentAddressConflict)
		suite.True(ok, "Expected response to be of type UpdateMTOShipmentAddressConflict")
		suite.Contains(*resp.Payload.Detail, "please update the storage facility address instead")
	})

	suite.Run("Failure - Unprocessable when updating address with invalid data", func() {
		// Testcase:   address is updated on a shipment that's available to MTO with invalid address
		// Expected:   Failure response 422
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: models.StringPointer("6622 Airport Way S #1430"),
			StreetAddress3: models.StringPointer("441 SW Río de la Plata Drive"),
			City:           "Bad City",
			State:          "CA",
			PostalCode:     "99999", // invalid postal code
		}

		// Update with new address
		payload := payloads.Address(&newAddress2)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressUnprocessableEntity{}, response)
	})

	suite.Run("Failure - Unprocessable with AK FF off and valid AK address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO with AK address but FF off
		// Expected:   Failure response 422
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: models.StringPointer("6622 Airport Way S #1430"),
			StreetAddress3: models.StringPointer("441 SW Río de la Plata Drive"),
			City:           "JUNEAU",
			State:          "AK",
			PostalCode:     "99801",
		}

		// setting the AK flag to false and use a valid address
		handlerConfig := suite.NewHandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: false,
		}

		mockFeatureFlagFetcher := &servicemocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig

		// Update with new address
		payload := payloads.Address(&newAddress2)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressUnprocessableEntity{}, response)
	})

	suite.Run("Failure - Unprocessable with HI FF off and valid HI address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO with HI address but FF off
		// Expected:   Failure response 422
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: models.StringPointer("6622 Airport Way S #1430"),
			StreetAddress3: models.StringPointer("441 SW Río de la Plata Drive"),
			City:           "HONOLULU",
			State:          "HI",
			PostalCode:     "96835",
		}

		// setting the HI flag to false and use a valid address
		handlerConfig := suite.NewHandlerConfig()

		expectedFeatureFlag := services.FeatureFlag{
			Key:   "enable_alaska",
			Match: false,
		}

		mockFeatureFlagFetcher := &servicemocks.FeatureFlagFetcher{}
		mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
			mock.Anything,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(expectedFeatureFlag, nil)
		handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		handler.HandlerConfig = handlerConfig

		// Update with new address
		payload := payloads.Address(&newAddress2)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressUnprocessableEntity{}, response)
	})

	suite.Run("Failure - Internal Error mock GetLocationsByZipCityState return error", func() {
		// Testcase:   address is updated on a shipment that's available to MTO with invalid address
		// Expected:   Failure response 422
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableMove,
				LinkOnly: true,
			},
		}, nil)
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: models.StringPointer("6622 Airport Way S #1430"),
			StreetAddress3: models.StringPointer("441 SW Río de la Plata Drive"),
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
		}

		// Update with new address
		payload := payloads.Address(&newAddress2)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(shipment.PickupAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(shipment.PickupAddress.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		expectedError := models.ErrFetchNotFound
		vLocationFetcher := &servicemocks.VLocation{}
		vLocationFetcher.On("GetLocationsByZipCityState",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()

		handler.VLocation = vLocationFetcher

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressInternalServerError{}, response)
	})
}
