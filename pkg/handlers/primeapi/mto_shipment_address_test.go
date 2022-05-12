package primeapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
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
	// Make an available MTO
	mto := testdatagen.MakeAvailableMove(suite.DB())

	// Make a shipment on the available MTO
	mtoShipment1 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})
	pickupAddress1 := mtoShipment1.PickupAddress

	newAddress := models.Address{
		StreetAddress1: "7 Q St",
		City:           "Framington",
		State:          "MA",
		PostalCode:     "94055",
	}

	// Create handler
	handler := UpdateMTOShipmentAddressHandler{
		handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
		mtoshipment.NewMTOShipmentAddressUpdater(),
	}

	var updatedETag string

	suite.T().Run("Success updating address", func(t *testing.T) {
		// Testcase:   address is updated on a shipment that's available to MTO
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object

		// Update with new address
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress1.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment1.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(pickupAddress1.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)

		// Check values
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)
		isAddressEqual(suite, payload, shipmentOk.Payload)
		updatedETag = shipmentOk.Payload.ETag
	})

	suite.T().Run("Success updating full address", func(t *testing.T) {
		// Testcase:   address is updated on a shipment that's available to MTO, all fields in address provided
		// Expected:   Success response 200
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object

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
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress1.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment1.ID),
			Body:          payload,
			IfMatch:       updatedETag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)

		// Check values
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)
		isAddressEqual(suite, payload, shipmentOk.Payload)
		updatedETag = shipmentOk.Payload.ETag

	})

	suite.T().Run("Fail - NotFound due to unavailable MTO", func(t *testing.T) {
		// Testcase:   address is updated on a shipment that's on an MTO that is NOT available to Prime
		// Expected:   NotFound error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object

		// Make a shipment with an unavailable MTO
		pickupAddress2 := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
		mtoShipment2 := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PickupAddress: &pickupAddress2,
			},
		})

		// Update with new address
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress2.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
			Body:          payload,
			IfMatch:       updatedETag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressNotFound{}, response)

	})
	suite.T().Run("Fail - ConflictError due to unassociated mtoShipment", func(t *testing.T) {
		// Testcase:   address is updated on a shipment that it's not associated with
		// Expected:   Conflict error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object

		// Make another shipment with an available MTO
		mto3 := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment3 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto3,
		})
		// Make a random address that is not associated
		randomAddress := testdatagen.MakeDefaultAddress(suite.DB())

		payload := payloads.Address(&randomAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment3.ID.String(), randomAddress.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(randomAddress.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment3.ID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(randomAddress.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressConflict{}, response)

	})
	suite.T().Run("Fail - PreconditionFailed due to wrong etag", func(t *testing.T) {
		// Testcase:   address is updated on a shipment, but etag for address is wrong
		// Expected:   PreconditionFailed error is returned
		// Under Test: UpdateMTOShipmentAddress handler code and mtoShipmentAddressUpdater service object

		// Update with new address with a bad etag
		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress1.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment1.ID),
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
