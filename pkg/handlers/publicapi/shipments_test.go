package publicapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"DRAFT"}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	params := shipmentop.GetShipmentParams{
		HTTPRequest:  req,
		ShipmentUUID: strfmt.UUID(shipment.ID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.GetShipmentOK{}, response)
	okResponse := response.(*shipmentop.GetShipmentOK)

	// And: Payload is equivalent to original shipment
	suite.Equal(strfmt.UUID(shipment.ID.String()), okResponse.Payload.ID)
}

func (suite *HandlerSuite) TestPatchShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"AWARDED"}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	genericDate := time.Now()
	UpdatePayload := apimessages.Shipment{
		PmSurveyPackDate:                    handlers.FmtDatePtr(&genericDate),
		PmSurveyPickupDate:                  handlers.FmtDatePtr(&genericDate),
		PmSurveyEarliestDeliveryDate:        handlers.FmtDatePtr(&genericDate),
		PmSurveyLatestDeliveryDate:          handlers.FmtDatePtr(&genericDate),
		PmSurveyWeightEstimate:              swag.Int64(33),
		PmSurveyProgearWeightEstimate:       swag.Int64(53),
		PmSurveySpouseProgearWeightEstimate: swag.Int64(54),
		PmSurveyNotes:                       swag.String("Unsure about pickup date."),
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest:  req,
		ShipmentUUID: strfmt.UUID(shipment.ID.String()),
		Update:       &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentOK{}, response)
	okResponse := response.(*shipmentop.PatchShipmentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(shipment.ID.String()), okResponse.Payload.ID)
	suite.Equal(*UpdatePayload.PmSurveyNotes, *okResponse.Payload.PmSurveyNotes)
	suite.Equal(int64(54), *okResponse.Payload.PmSurveySpouseProgearWeightEstimate)
	suite.Equal(int64(53), *okResponse.Payload.PmSurveyProgearWeightEstimate)
	suite.Equal(int64(33), *okResponse.Payload.PmSurveyWeightEstimate)
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyLatestDeliveryDate))
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyEarliestDeliveryDate))
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyPickupDate))
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyPackDate))
}

func (suite *HandlerSuite) TestPatchShipmentHandlerWrongTSP() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"AWARDED"}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]

	otherTspUser := testdatagen.MakeDefaultTspUser(suite.db)

	// And: the context contains the auth values for the wrong tsp
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, otherTspUser)

	genericDate := time.Now()
	UpdatePayload := apimessages.Shipment{
		PmSurveyPackDate:                    handlers.FmtDatePtr(&genericDate),
		PmSurveyPickupDate:                  handlers.FmtDatePtr(&genericDate),
		PmSurveyEarliestDeliveryDate:        handlers.FmtDatePtr(&genericDate),
		PmSurveyLatestDeliveryDate:          handlers.FmtDatePtr(&genericDate),
		PmSurveyWeightEstimate:              swag.Int64(33),
		PmSurveyProgearWeightEstimate:       swag.Int64(53),
		PmSurveySpouseProgearWeightEstimate: swag.Int64(54),
		PmSurveyNotes:                       swag.String("Unsure about pickup date."),
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest:  req,
		ShipmentUUID: strfmt.UUID(shipment.ID.String()),
		Update:       &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentBadRequest{}, response)
}

// TestIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"DRAFT"}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

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

// TestIndexShipmentsHandlerPaginated tests the api endpoint with pagination query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser1 := tspUsers[0]
	tspUser2 := tspUsers[1]

	// Constants
	limit := int64(25)
	offset := int64(1)

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}

	// Test query with first user
	req1 := httptest.NewRequest("GET", "/shipments", nil)
	req1 = suite.authenticateTspRequest(req1, tspUser1)
	params1 := shipmentop.IndexShipmentsParams{
		HTTPRequest: req1,
		Limit:       &limit,
		Offset:      &offset,
	}

	response1 := handler.Handle(params1)
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response1)
	okResponse1 := response1.(*shipmentop.IndexShipmentsOK)
	suite.Equal(15, len(okResponse1.Payload))

	// Test query with second user
	req2 := httptest.NewRequest("GET", "/shipments", nil)
	req2 = suite.authenticateTspRequest(req2, tspUser2)
	params2 := shipmentop.IndexShipmentsParams{
		HTTPRequest: req2,
		Limit:       &limit,
		Offset:      &offset,
	}

	response2 := handler.Handle(params2)
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response2)
	okResponse2 := response2.(*shipmentop.IndexShipmentsOK)
	suite.Equal(10, len(okResponse2.Payload))
}

// TestIndexShipmentsHandlerSortShipmentsPickupAsc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsPickupAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
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
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

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

// TestIndexShipmentsHandlerSortShipmentsPickupDesc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsPickupDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
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
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

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

// TestIndexShipmentsHandlerSortShipmentsDeliveryAsc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsDeliveryAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "DELIVERY_DATE_ASC"
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

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

// TestIndexShipmentsHandlerSortShipmentsDeliveryDesc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsDeliveryDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	limit := int64(25)
	offset := int64(1)
	orderBy := "DELIVERY_DATE_DESC"
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Limit:       &limit,
		Offset:      &offset,
		OrderBy:     &orderBy,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

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

// TestIndexShipmentsHandlerFilterByStatus tests the api endpoint with defined status query param
func (suite *HandlerSuite) TestIndexShipmentsHandlerFilterByStatus() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      status,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)
	suite.Equal(25, len(okResponse.Payload))
}

// TestIndexShipmentsHandlerFilterByStatusNoResults tests the api endpoint with defined status query param that returns nothing
func (suite *HandlerSuite) TestIndexShipmentsHandlerFilterByStatusNoResults() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DRAFT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.db, suite.logger)}
	statusFilter := []string{"NOTASTATUS"}

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      statusFilter,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)
	suite.Equal(0, len(okResponse.Payload))
}
