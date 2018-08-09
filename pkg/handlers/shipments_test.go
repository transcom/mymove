package handlers

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	publicshipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

/*
 * ------------------------------------------
 * The code below is for the INTERNAL REST API.
 * ------------------------------------------
 */

func (suite *HandlerSuite) verifyAddressFields(expected, actual *internalmessages.Address) {
	suite.T().Helper()
	suite.Equal(expected.StreetAddress1, actual.StreetAddress1, "Street1 did not match")
	suite.Equal(expected.StreetAddress2, actual.StreetAddress2, "Street2 did not match")
	suite.Equal(expected.StreetAddress3, actual.StreetAddress3, "Street3 did not match")
	suite.Equal(expected.City, actual.City, "City did not match")
	suite.Equal(expected.State, actual.State, "State did not match")
	suite.Equal(expected.PostalCode, actual.PostalCode, "PostalCode did not match")
	suite.Equal(expected.Country, actual.Country, "Country did not match")
}

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
	market := "dHHG"
	codeOfService := "D"

	suite.Equal(strfmt.UUID(move.ID.String()), unwrapped.Payload.MoveID)
	suite.Equal(strfmt.UUID(sm.ID.String()), unwrapped.Payload.ServiceMemberID)
	suite.Equal("DRAFT", unwrapped.Payload.Status)
	suite.Equal(&codeOfService, unwrapped.Payload.CodeOfService)
	suite.Equal(&market, unwrapped.Payload.Market)
	suite.Equal(swag.Int64(2), unwrapped.Payload.EstimatedPackDays)
	suite.Equal(swag.Int64(5), unwrapped.Payload.EstimatedTransitDays)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.PickupAddress)
	suite.Equal(true, unwrapped.Payload.HasSecondaryPickupAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.SecondaryPickupAddress)
	suite.Equal(true, unwrapped.Payload.HasDeliveryAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.DeliveryAddress)
	suite.Equal(true, unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.verifyAddressFields(addressPayload, unwrapped.Payload.PartialSitDeliveryAddress)
	suite.Equal(swag.Int64(4500), unwrapped.Payload.WeightEstimate)
	suite.Equal(swag.Int64(325), unwrapped.Payload.ProgearWeightEstimate)
	suite.Equal(swag.Int64(120), unwrapped.Payload.SpouseProgearWeightEstimate)

	count, err := suite.db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)
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

	market := "dHHG"
	codeOfService := "D"
	suite.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)

	count, err := suite.db.Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)

	suite.Equal(strfmt.UUID(move.ID.String()), unwrapped.Payload.MoveID)
	suite.Equal(strfmt.UUID(sm.ID.String()), unwrapped.Payload.ServiceMemberID)
	suite.Equal("DRAFT", unwrapped.Payload.Status)
	suite.Equal(&market, unwrapped.Payload.Market)
	suite.Equal(&codeOfService, unwrapped.Payload.CodeOfService)
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
		ServiceMemberID:              sm.ID,
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

	// assert we got back the 200 response
	okResponse := response.(*shipmentop.PatchShipmentOK)
	patchShipmentPayload := okResponse.Payload

	suite.Equal(patchShipmentPayload.HasDeliveryAddress, true, "HasDeliveryAddress should have been updated.")
	suite.verifyAddressFields(newAddress, patchShipmentPayload.DeliveryAddress)

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
		ServiceMemberID:              sm.ID,
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

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

// TestPublicIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
	}

	// And: an index of shipments is returned
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.Equal(1, len(okResponse.Payload))
	if len(okResponse.Payload) == 1 {
		responsePayload := okResponse.Payload[0]
		// And: Payload is equivalent to original shipment
		suite.Equal(strfmt.UUID(shipment.ID.String()), responsePayload.ID)
		suite.Equal(apimessages.SelectedMoveType(*shipment.Move.SelectedMoveType), *responsePayload.Move.SelectedMoveType)
		suite.Equal(shipment.TrafficDistributionList.SourceRateArea, *responsePayload.TrafficDistributionList.SourceRateArea)
	}
}

// TestPublicIndexShipmentsHandlerPaginated tests the api endpoint with pagination query parameters
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser1 := tspUsers[0]
	tspUser2 := tspUsers[1]

	// Constants
	limit := int64(25)
	offset := int64(1)

	// Handler to Test
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))

	// Test query with first user
	req1 := httptest.NewRequest("GET", "/shipments", nil)
	req1 = suite.authenticateTspRequest(req1, tspUser1)
	params1 := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req1,
		Limit:       &limit,
		Offset:      &offset,
	}

	response1 := handler.Handle(params1)
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response1)
	okResponse1 := response1.(*publicshipmentop.IndexShipmentsOK)
	suite.Equal(15, len(okResponse1.Payload))

	// Test query with second user
	req2 := httptest.NewRequest("GET", "/shipments", nil)
	req2 = suite.authenticateTspRequest(req2, tspUser2)
	params2 := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req2,
		Limit:       &limit,
		Offset:      &offset,
	}

	response2 := handler.Handle(params2)
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response2)
	okResponse2 := response2.(*publicshipmentop.IndexShipmentsOK)
	suite.Equal(10, len(okResponse2.Payload))
}

// TestPublicIndexShipmentsHandlerSortShipmentsPickupAsc sorts returned shipments
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsPickupAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "PICKUP_DATE_ASC"
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.Equal(3, len(okResponse.Payload))

	var pickupDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if pickupDate == empty {
			pickupDate = time.Time(responsePayload.PickupDate)
		} else {
			newDT := time.Time(responsePayload.PickupDate)
			suite.True(newDT.After(pickupDate))
			pickupDate = newDT
		}
	}
}

// TestPublicIndexShipmentsHandlerSortShipmentsPickupDesc sorts returned shipments
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsPickupDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "PICKUP_DATE_DESC"
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.Equal(3, len(okResponse.Payload))

	var pickupDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if pickupDate == empty {
			pickupDate = time.Time(responsePayload.PickupDate)
		} else {
			newDT := time.Time(responsePayload.PickupDate)
			suite.True(newDT.Before(pickupDate))
			pickupDate = newDT
		}
	}
}

// TestPublicIndexShipmentsHandlerSortShipmentsDeliveryAsc sorts returned shipments
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsDeliveryAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "DELIVERY_DATE_ASC"
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.Equal(3, len(okResponse.Payload))

	var deliveryDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if deliveryDate == empty {
			deliveryDate = time.Time(responsePayload.DeliveryDate)
		} else {
			newDT := time.Time(responsePayload.DeliveryDate)
			suite.True(newDT.After(deliveryDate))
			deliveryDate = newDT
		}
	}
}

// TestPublicIndexShipmentsHandlerSortShipmentsDeliveryDesc sorts returned shipments
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsDeliveryDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "DELIVERY_DATE_DESC"
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.Equal(3, len(okResponse.Payload))

	var deliveryDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if deliveryDate == empty {
			deliveryDate = time.Time(responsePayload.DeliveryDate)
		} else {
			newDT := time.Time(responsePayload.DeliveryDate)
			suite.True(newDT.Before(deliveryDate))
			deliveryDate = newDT
		}
	}
}

// TestPublicIndexShipmentsHandlerFilterByStatus tests the api endpoint with defined status query param
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerFilterByStatus() {

	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// Constants
	status := []string{"DEFAULT"}

	// Handler to Test
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      status,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)
	suite.Equal(25, len(okResponse.Payload))
}

// TestPublicIndexShipmentsHandlerFilterByStatusNoResults tests the api endpoint with defined status query param that returns nothing
func (suite *HandlerSuite) TestPublicIndexShipmentsHandlerFilterByStatusNoResults() {

	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// Constants
	status := []string{"NOTASTATUS"}

	// Handler to Test
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      status,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)
	suite.Equal(0, len(okResponse.Payload))
}
