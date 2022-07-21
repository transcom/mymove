package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/testdatagen"
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
		availableMove := testdatagen.MakeAvailableMove(suite.DB())

		// Create handler
		handler := UpdateMTOShipmentAddressHandler{
			suite.HandlerConfig(),
			mtoshipment.NewMTOShipmentAddressUpdater(),
		}
		return handler, availableMove
	}

	newAddress := models.Address{
		StreetAddress1: "7 Q St",
		City:           "Framington",
		State:          "MA",
		PostalCode:     "94055",
	}

	suite.Run("Success updating address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		// Make a shipment on the available MTO
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})

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
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)

		// Check values
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)
		isAddressEqual(suite, payload, shipmentOk.Payload)
	})

	suite.Run("Success updating full address", func() {
		// Testcase:   address is updated on a shipment that's available to MTO, all fields in address provided
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		newAddress2 := models.Address{
			StreetAddress1: "7 Q St",
			StreetAddress2: swag.String("6622 Airport Way S #1430"),
			StreetAddress3: swag.String("441 SW Río de la Plata Drive"),
			City:           "Alameda",
			State:          "CA",
			PostalCode:     "94055",
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
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)

		// Check values
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)
		isAddressEqual(suite, payload, shipmentOk.Payload)

	})

	suite.Run("Fail - NotFound due to unavailable MTO", func() {
		// Testcase:   address is updated on a shipment that's on an MTO that is NOT available to Prime
		// Expected:   NotFound error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, _ := setupTestData()
		// Make a shipment with an unavailable MTO
		pickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PickupAddress: &pickupAddress,
			},
		})

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
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressNotFound{}, response)

	})
	suite.Run("Fail - ConflictError due to unassociated mtoShipment", func() {
		// Testcase:   address is updated on a shipment that it's not associated with
		// Expected:   Conflict error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		// Make a random address that is not associated
		randomAddress := testdatagen.MakeDefaultAddress(suite.DB())

		payload := payloads.Address(&randomAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", shipment.ID.String(), randomAddress.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(randomAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(randomAddress.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressConflict{}, response)

	})
	suite.Run("Fail - PreconditionFailed due to wrong etag", func() {
		// Testcase:   address is updated on a shipment, but etag for address is wrong
		// Expected:   PreconditionFailed error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object
		handler, availableMove := setupTestData()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
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
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressPreconditionFailed{}, response)

	})

}
