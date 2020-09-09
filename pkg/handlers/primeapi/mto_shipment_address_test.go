package primeapi

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
func isAddressEqual(suite *HandlerSuite, reqAddress *primemessages.Address, respAddress *primemessages.Address) {
	if reqAddress.StreetAddress1 != nil {
		suite.Equal(*reqAddress.StreetAddress1, *respAddress.StreetAddress1)
	}
	if reqAddress.StreetAddress2 != nil {
		suite.Equal(*reqAddress.StreetAddress2, *respAddress.StreetAddress2)
	}
	if reqAddress.StreetAddress3 != nil {
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
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewMTOShipmentAddressUpdater(suite.DB()),
	}

	suite.T().Run("Successful case updating address", func(t *testing.T) {

		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress1.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment1.ID),
			Body:          payload,
			IfMatch:       "who-cares",
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressOK{}, response)

		// Check values
		shipmentOk := response.(*mtoshipmentops.UpdateMTOShipmentAddressOK)
		respPayload := shipmentOk.Payload
		isAddressEqual(suite, payload, respPayload)
		suite.True(true)
	})
	suite.T().Run("Fail - NotFound due to unavailable MTO", func(t *testing.T) {

		// Make a shipment with an unavailable MTO
		mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})
		pickupAddress2 := mtoShipment1.PickupAddress

		payload := payloads.Address(&newAddress)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/addresses/%s", mtoShipment1.ID.String(), mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentAddressParams{
			HTTPRequest:   req,
			AddressID:     *handlers.FmtUUID(pickupAddress2.ID),
			MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
			Body:          payload,
			IfMatch:       "who-cares",
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressNotFound{}, response)

	})
	suite.T().Run("Fail - ConflictError due to unassociated mtoShipment", func(t *testing.T) {

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
			IfMatch:       "who-cares",
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentAddressConflict{}, response)

	})

}
