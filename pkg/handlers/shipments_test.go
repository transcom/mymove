package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateShipmentHandlerAllValues() {
	move := testdatagen.MakeMove(suite.db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	addressPayload := fakeAddressPayload()

	newShipment := internalmessages.Shipment{
		EstimatedPackDays:            swag.Int64(2),
		EstimatedTransitDays:         swag.Int64(5),
		PickupAddress:                addressPayload,
		HasSecondaryPickupAddress:    true,
		SecondaryPickupAddress:       addressPayload,
		HasDeliveryAddress:           true,
		DeliveryAddress:              addressPayload,
		HasPartialSitDeliveryAddress: true,
		PartialSitDeliveryAddress:    addressPayload,
		WeightEstimate:               swag.Int64(4500),
		ProgearWeightEstimate:        swag.Int64(325),
		SpouseProgearWeightEstimate:  swag.Int64(120),
	}

	req := httptest.NewRequest("POST", "/moves/move_id/shipment", nil)
	req = suite.authenticateRequest(req, sm)

	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)

	count, err := suite.db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)

	suite.Equal("DRAFT", unwrapped.Payload.Status)
	suite.Equal(swag.Int64(2), unwrapped.Payload.EstimatedPackDays)
	suite.Equal(swag.Int64(5), unwrapped.Payload.EstimatedTransitDays)
	suite.Equal(addressPayload, unwrapped.Payload.PickupAddress)
	suite.Equal(true, unwrapped.Payload.HasSecondaryPickupAddress)
	suite.Equal(addressPayload, unwrapped.Payload.SecondaryPickupAddress)
	suite.Equal(true, unwrapped.Payload.HasDeliveryAddress)
	suite.Equal(addressPayload, unwrapped.Payload.DeliveryAddress)
	suite.Equal(true, unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.Equal(addressPayload, unwrapped.Payload.PartialSitDeliveryAddress)
	suite.Equal(swag.Int64(4500), unwrapped.Payload.WeightEstimate)
	suite.Equal(swag.Int64(325), unwrapped.Payload.ProgearWeightEstimate)
	suite.Equal(swag.Int64(120), unwrapped.Payload.SpouseProgearWeightEstimate)
}

func (suite *HandlerSuite) TestCreateShipmentHandlerEmpty() {
	move := testdatagen.MakeMove(suite.db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	req := httptest.NewRequest("POST", "/moves/move_id/shipment", nil)
	req = suite.authenticateRequest(req, sm)

	newShipment := internalmessages.Shipment{}
	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)

	count, err := suite.db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)

	suite.Equal("DRAFT", unwrapped.Payload.Status)
	suite.Nil(unwrapped.Payload.EstimatedPackDays)
	suite.Nil(unwrapped.Payload.EstimatedTransitDays)
	suite.Nil(unwrapped.Payload.PickupAddress)
	suite.Equal(false, unwrapped.Payload.HasSecondaryPickupAddress)
	suite.Nil(unwrapped.Payload.SecondaryPickupAddress)
	suite.Equal(false, unwrapped.Payload.HasDeliveryAddress)
	suite.Nil(unwrapped.Payload.DeliveryAddress)
	suite.Equal(false, unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.Nil(unwrapped.Payload.PartialSitDeliveryAddress)
	suite.Nil(unwrapped.Payload.WeightEstimate)
	suite.Nil(unwrapped.Payload.ProgearWeightEstimate)
	suite.Nil(unwrapped.Payload.SpouseProgearWeightEstimate)
}

func (suite *HandlerSuite) TestPatchShipmentsHandlerHappyPath() {
	move := testdatagen.MakeMove(suite.db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	addressPayload := testdatagen.MakeAddress(suite.db, testdatagen.Assertions{})

	shipment1 := models.Shipment{
		MoveID:                       move.ID,
		Status:                       "DRAFT",
		EstimatedPackDays:            swag.Int64(2),
		EstimatedTransitDays:         swag.Int64(5),
		PickupAddress:                &addressPayload,
		HasSecondaryPickupAddress:    true,
		SecondaryPickupAddress:       &addressPayload,
		HasDeliveryAddress:           false,
		HasPartialSITDeliveryAddress: true,
		PartialSITDeliveryAddress:    &addressPayload,
		WeightEstimate:               poundPtrFromInt64Ptr(swag.Int64(4500)),
		ProgearWeightEstimate:        poundPtrFromInt64Ptr(swag.Int64(325)),
		SpouseProgearWeightEstimate:  poundPtrFromInt64Ptr(swag.Int64(120)),
	}
	suite.mustSave(&shipment1)

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.authenticateRequest(req, sm)

	newAddress := otherFakeAddressPayload()

	payload := internalmessages.Shipment{
		EstimatedPackDays:           swag.Int64(15),
		HasSecondaryPickupAddress:   false,
		HasDeliveryAddress:          true,
		DeliveryAddress:             newAddress,
		SpouseProgearWeightEstimate: swag.Int64(100),
	}

	patchShipmentParams := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		ShipmentID:  strfmt.UUID(shipment1.ID.String()),
		Shipment:    &payload,
	}

	handler := PatchShipmentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchShipmentParams)

	// assert we got back the 201 response
	okResponse := response.(*shipmentop.PatchShipmentCreated)
	patchShipmentPayload := okResponse.Payload

	suite.Equal(patchShipmentPayload.HasDeliveryAddress, true, "HasDeliveryAddress should have been updated.")
	suite.Equal(patchShipmentPayload.DeliveryAddress, newAddress, "DeliveryAddress should have been updated.")

	suite.Equal(patchShipmentPayload.HasSecondaryPickupAddress, false, "HasSecondaryPickupAddress should have been updated.")
	suite.Nil(patchShipmentPayload.SecondaryPickupAddress, "SecondaryPickupAddress should have been updated to nil.")

	suite.Equal(*patchShipmentPayload.EstimatedPackDays, int64(15), "EstimatedPackDays should have been set to 15")
	suite.Equal(*patchShipmentPayload.SpouseProgearWeightEstimate, int64(100), "SpouseProgearWeightEstimate should have been set to 100")
}

func (suite *HandlerSuite) TestPatchShipmentHandlerNoMove() {
	t := suite.T()
	move := testdatagen.MakeMove(suite.db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember
	badMoveID := uuid.Must(uuid.NewV4())

	addressPayload := testdatagen.MakeAddress(suite.db, testdatagen.Assertions{})

	shipment1 := models.Shipment{
		MoveID:                       move.ID,
		Status:                       "DRAFT",
		EstimatedPackDays:            swag.Int64(2),
		EstimatedTransitDays:         swag.Int64(5),
		PickupAddress:                &addressPayload,
		HasSecondaryPickupAddress:    true,
		SecondaryPickupAddress:       &addressPayload,
		HasDeliveryAddress:           false,
		HasPartialSITDeliveryAddress: true,
		PartialSITDeliveryAddress:    &addressPayload,
		WeightEstimate:               poundPtrFromInt64Ptr(swag.Int64(4500)),
		ProgearWeightEstimate:        poundPtrFromInt64Ptr(swag.Int64(325)),
		SpouseProgearWeightEstimate:  poundPtrFromInt64Ptr(swag.Int64(120)),
	}
	suite.mustSave(&shipment1)

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.authenticateRequest(req, sm)

	newAddress := otherFakeAddressPayload()

	payload := internalmessages.Shipment{
		EstimatedPackDays:           swag.Int64(15),
		HasSecondaryPickupAddress:   false,
		HasDeliveryAddress:          true,
		DeliveryAddress:             newAddress,
		SpouseProgearWeightEstimate: swag.Int64(100),
	}

	patchShipmentParams := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(badMoveID.String()),
		ShipmentID:  strfmt.UUID(shipment1.ID.String()),
		Shipment:    &payload,
	}

	handler := PatchShipmentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchShipmentParams)

	// assert we got back the badrequest response
	_, ok := response.(*shipmentop.PatchShipmentBadRequest)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

}
