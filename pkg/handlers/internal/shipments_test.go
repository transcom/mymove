package internal

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) verifyAddressFields(expected, actual *internalmessages.Address) {
	suite.parent.T().Helper()
	suite.parent.Equal(expected.StreetAddress1, actual.StreetAddress1, "Street1 did not match")
	suite.parent.Equal(expected.StreetAddress2, actual.StreetAddress2, "Street2 did not match")
	suite.parent.Equal(expected.StreetAddress3, actual.StreetAddress3, "Street3 did not match")
	suite.parent.Equal(expected.City, actual.City, "City did not match")
	suite.parent.Equal(expected.State, actual.State, "State did not match")
	suite.parent.Equal(expected.PostalCode, actual.PostalCode, "PostalCode did not match")
	suite.parent.Equal(expected.Country, actual.Country, "Country did not match")
}

func (suite *HandlerSuite) TestCreateShipmentHandlerAllValues() {
	move := testdatagen.MakeMove(suite.parent.Db, testdatagen.Assertions{})
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
	req = suite.parent.AuthenticateRequest(req, sm)

	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)
	market := "dHHG"
	codeOfService := "D"

	suite.parent.Equal(strfmt.UUID(move.ID.String()), unwrapped.Payload.MoveID)
	suite.parent.Equal(strfmt.UUID(sm.ID.String()), unwrapped.Payload.ServiceMemberID)
	suite.parent.Equal("DRAFT", unwrapped.Payload.Status)
	suite.parent.Equal(&codeOfService, unwrapped.Payload.CodeOfService)
	suite.parent.Equal(&market, unwrapped.Payload.Market)
	suite.parent.Equal(swag.Int64(2), unwrapped.Payload.EstimatedPackDays)
	suite.parent.Equal(swag.Int64(5), unwrapped.Payload.EstimatedTransitDays)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.PickupAddress)
	suite.parent.Equal(true, unwrapped.Payload.HasSecondaryPickupAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.SecondaryPickupAddress)
	suite.parent.Equal(true, unwrapped.Payload.HasDeliveryAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.DeliveryAddress)
	suite.parent.Equal(true, unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.PartialSitDeliveryAddress)
	suite.parent.Equal(swag.Int64(4500), unwrapped.Payload.WeightEstimate)
	suite.parent.Equal(swag.Int64(325), unwrapped.Payload.ProgearWeightEstimate)
	suite.parent.Equal(swag.Int64(120), unwrapped.Payload.SpouseProgearWeightEstimate)

	count, err := suite.parent.Db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.parent.Nil(err, "could not count shipments")
	suite.parent.Equal(1, count)
}

func (suite *HandlerSuite) TestCreateShipmentHandlerEmpty() {
	move := testdatagen.MakeMove(suite.parent.Db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	req := httptest.NewRequest("POST", "/moves/move_id/shipment", nil)
	req = suite.parent.AuthenticateRequest(req, sm)

	newShipment := internalmessages.Shipment{}
	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	market := "dHHG"
	codeOfService := "D"
	suite.parent.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)

	count, err := suite.parent.Db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.parent.Nil(err, "could not count shipments")
	suite.parent.Equal(1, count)

	suite.parent.Equal(strfmt.UUID(move.ID.String()), unwrapped.Payload.MoveID)
	suite.parent.Equal(strfmt.UUID(sm.ID.String()), unwrapped.Payload.ServiceMemberID)
	suite.parent.Equal("DRAFT", unwrapped.Payload.Status)
	suite.parent.Equal(&market, unwrapped.Payload.Market)
	suite.parent.Equal(&codeOfService, unwrapped.Payload.CodeOfService)
	suite.parent.Nil(unwrapped.Payload.EstimatedPackDays)
	suite.parent.Nil(unwrapped.Payload.EstimatedTransitDays)
	suite.parent.Nil(unwrapped.Payload.PickupAddress)
	suite.parent.Equal(false, unwrapped.Payload.HasSecondaryPickupAddress)
	suite.parent.Nil(unwrapped.Payload.SecondaryPickupAddress)
	suite.parent.Equal(false, unwrapped.Payload.HasDeliveryAddress)
	suite.parent.Nil(unwrapped.Payload.DeliveryAddress)
	suite.parent.Equal(false, unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.parent.Nil(unwrapped.Payload.PartialSitDeliveryAddress)
	suite.parent.Nil(unwrapped.Payload.WeightEstimate)
	suite.parent.Nil(unwrapped.Payload.ProgearWeightEstimate)
	suite.parent.Nil(unwrapped.Payload.SpouseProgearWeightEstimate)
}

func (suite *HandlerSuite) TestPatchShipmentsHandlerHappyPath() {
	move := testdatagen.MakeMove(suite.parent.Db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	addressPayload := testdatagen.MakeAddress(suite.parent.Db, testdatagen.Assertions{})

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
		WeightEstimate:               utils.PoundPtrFromInt64Ptr(swag.Int64(4500)),
		ProgearWeightEstimate:        utils.PoundPtrFromInt64Ptr(swag.Int64(325)),
		SpouseProgearWeightEstimate:  utils.PoundPtrFromInt64Ptr(swag.Int64(120)),
		ServiceMemberID:              sm.ID,
	}
	suite.parent.MustSave(&shipment1)

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.parent.AuthenticateRequest(req, sm)

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

	handler := PatchShipmentHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(patchShipmentParams)

	// assert we got back the 200 response
	okResponse := response.(*shipmentop.PatchShipmentOK)
	patchShipmentPayload := okResponse.Payload

	suite.parent.Equal(patchShipmentPayload.HasDeliveryAddress, true, "HasDeliveryAddress should have been updated.")
	suite.verifyAddressFields(newAddress, patchShipmentPayload.DeliveryAddress)

	suite.parent.Equal(patchShipmentPayload.HasSecondaryPickupAddress, false, "HasSecondaryPickupAddress should have been updated.")
	suite.parent.Nil(patchShipmentPayload.SecondaryPickupAddress, "SecondaryPickupAddress should have been updated to nil.")

	suite.parent.Equal(*patchShipmentPayload.EstimatedPackDays, int64(15), "EstimatedPackDays should have been set to 15")
	suite.parent.Equal(*patchShipmentPayload.SpouseProgearWeightEstimate, int64(100), "SpouseProgearWeightEstimate should have been set to 100")
}

func (suite *HandlerSuite) TestPatchShipmentHandlerNoMove() {
	t := suite.parent.T()
	move := testdatagen.MakeMove(suite.parent.Db, testdatagen.Assertions{})
	sm := move.Orders.ServiceMember
	badMoveID := uuid.Must(uuid.NewV4())

	addressPayload := testdatagen.MakeAddress(suite.parent.Db, testdatagen.Assertions{})

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
		WeightEstimate:               utils.PoundPtrFromInt64Ptr(swag.Int64(4500)),
		ProgearWeightEstimate:        utils.PoundPtrFromInt64Ptr(swag.Int64(325)),
		SpouseProgearWeightEstimate:  utils.PoundPtrFromInt64Ptr(swag.Int64(120)),
		ServiceMemberID:              sm.ID,
	}
	suite.parent.MustSave(&shipment1)

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.parent.AuthenticateRequest(req, sm)

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

	handler := PatchShipmentHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(patchShipmentParams)

	// assert we got back the badrequest response
	_, ok := response.(*shipmentop.PatchShipmentBadRequest)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
}
