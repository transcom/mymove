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
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestGetShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.GetShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.GetShipmentOK{}, response)
	okResponse := response.(*shipmentop.GetShipmentOK)

	// And: Payload is equivalent to original shipment
	suite.Equal(strfmt.UUID(shipment.ID.String()), okResponse.Payload.ID)
}

func (suite *HandlerSuite) TestPatchShipmentHandlerActualWeight() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.Shipment{
		ActualWeight: swag.Int64(17500),
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Update:      &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentOK{}, response)
	okResponse := response.(*shipmentop.PatchShipmentOK)

	// And: Payload has new values
	suite.Equal(int64(17500), *okResponse.Payload.ActualWeight)
}

func (suite *HandlerSuite) TestPatchShipmentHandlerPmSurvey() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	genericDate := time.Now()
	UpdatePayload := apimessages.Shipment{
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(&genericDate),
		PmSurveyPlannedPickupDate:           handlers.FmtDatePtr(&genericDate),
		PmSurveyPlannedDeliveryDate:         handlers.FmtDatePtr(&genericDate),
		PmSurveyWeightEstimate:              swag.Int64(33),
		PmSurveyProgearWeightEstimate:       swag.Int64(53),
		PmSurveySpouseProgearWeightEstimate: swag.Int64(54),
		PmSurveyNotes:                       swag.String("Unsure about pickup date."),
		PmSurveyMethod:                      "PHONE",
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Update:      &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentOK{}, response)
	okResponse := response.(*shipmentop.PatchShipmentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(shipment.ID.String()), okResponse.Payload.ID)
	suite.Equal(*UpdatePayload.PmSurveyNotes, *okResponse.Payload.PmSurveyNotes)
	suite.Equal(UpdatePayload.PmSurveyMethod, okResponse.Payload.PmSurveyMethod)
	suite.Equal(int64(54), *okResponse.Payload.PmSurveySpouseProgearWeightEstimate)
	suite.Equal(int64(53), *okResponse.Payload.PmSurveyProgearWeightEstimate)
	suite.Equal(int64(33), *okResponse.Payload.PmSurveyWeightEstimate)
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyPlannedDeliveryDate))
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyPlannedPickupDate))
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyPlannedPackDate))
}

func (suite *HandlerSuite) TestPatchShipmentHandlerPmSurveyWrongTSP() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]

	otherTspUser := testdatagen.MakeDefaultTspUser(suite.TestDB())

	// And: the context contains the auth values for the wrong tsp
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, otherTspUser)

	genericDate := time.Now()
	UpdatePayload := apimessages.Shipment{
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(&genericDate),
		PmSurveyPlannedPickupDate:           handlers.FmtDatePtr(&genericDate),
		PmSurveyPlannedDeliveryDate:         handlers.FmtDatePtr(&genericDate),
		PmSurveyWeightEstimate:              swag.Int64(33),
		PmSurveyProgearWeightEstimate:       swag.Int64(53),
		PmSurveySpouseProgearWeightEstimate: swag.Int64(54),
		PmSurveyNotes:                       swag.String("Unsure about pickup date."),
		PmSurveyMethod:                      "PHONE",
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Update:      &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentBadRequest{}, response)
}

// TestIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
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

// TestCreateGovBillOfLadingHandler
func (suite *HandlerSuite) TestCreateGovBillOfLadingHandler() {

	// When: There is a move, shipment, and move document of type GBL
	tspUser := testdatagen.MakeDefaultTspUser(suite.TestDB())
	shipment := scenario.MakeHhgFromAwardedToAccepted(suite.TestDB(), tspUser)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.CreateGovBillOfLadingParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// And: an index of shipments is returned
	// TODO: split out os calls from handler and/or mock OS to be able to test the handler
	handler := CreateGovBillOfLadingHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.CreateGovBillOfLadingCreated{}, response)

	// // When: there is an existing GBL for a shipment and handler is called
	// handler = CreateGovBillOfLadingHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	// response = handler.Handle(params)

	// // Then: expect a 400 status code
	// suite.Assertions.IsType(&shipmentop.CreateGovBillOfLadingBadRequest{}, response)

	// // When: an unauthed TSP user hits the handler
	// req = suite.AuthenticateTspRequest(req, unauthedTSPUser)
	// params.HTTPRequest = req
	// response = handler.Handle(params)

	// // Then: expect a 400 status code
	// suite.Assertions.IsType(&shipmentop.CreateGovBillOfLadingBadRequest{}, response)
}

// TestIndexShipmentsHandlerPaginated tests the api endpoint with pagination query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser1 := tspUsers[0]
	tspUser2 := tspUsers[1]

	// Constants
	limit := int64(25)
	offset := int64(1)

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}

	// Test query with first user
	req1 := httptest.NewRequest("GET", "/shipments", nil)
	req1 = suite.AuthenticateTspRequest(req1, tspUser1)
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
	req2 = suite.AuthenticateTspRequest(req2, tspUser2)
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}

	// The params expect statuses in strings, so they have to be cast from ShipmentStatus types
	stringStatus := []string{string(models.ShipmentStatusDRAFT)}
	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      stringStatus,
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
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	statusFilter := []string{"NOTASTATUS"}

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      statusFilter,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)
	suite.Equal(0, len(okResponse.Payload))
}

// TestAcceptShipmentHandler tests the api endpoint that accepts a shipment
func (suite *HandlerSuite) TestAcceptShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Handler to Test
	handler := AcceptShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}

	// Test query with first user
	path := fmt.Sprintf("/shipments/%s/accept", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := shipmentop.AcceptShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.AcceptShipmentOK{}, response)
	okResponse := response.(*shipmentop.AcceptShipmentOK)
	suite.Equal("ACCEPTED", string(okResponse.Payload.Status))
}

// TestRejectShipmentHandler tests the api endpoint that rejects a shipment
func (suite *HandlerSuite) TestRejectShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Handler to Test
	handler := RejectShipmentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}

	// Test query with first user
	path := fmt.Sprintf("/shipments/%s/reject", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	reason := "To Test Rejection"
	body := apimessages.RejectShipment{
		Reason: &reason,
	}
	params := shipmentop.RejectShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
		Payload:     &body,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.RejectShipmentOK{}, response)
	okResponse := response.(*shipmentop.RejectShipmentOK)
	suite.Equal("SUBMITTED", string(okResponse.Payload.Status))

	shipmentOffer, err := models.FetchShipmentOfferByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.NoError(err)

	suite.Equal(false, *shipmentOffer.Accepted)
	suite.Equal(reason, *shipmentOffer.RejectionReason)

}
