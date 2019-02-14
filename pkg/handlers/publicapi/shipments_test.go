package publicapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/paperwork"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/route"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	paperworkservice "github.com/transcom/mymove/pkg/services/paperwork"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestPayloadForShipmentModelWhenTspIDIsNotPresent() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	shipmentPayload := payloadForShipmentModel(shipment)
	suite.Equal(shipmentPayload.TransportationServiceProviderID, strfmt.UUID(""))
}

func (suite *HandlerSuite) TestPayloadForShipmentModelWhenTspIDIsPresent() {
	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{})
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tsp.ID,
			ShipmentID:                      shipment.ID,
		},
	})
	reloadShipment, err := models.FetchShipmentByTSP(suite.DB(), tsp.ID, shipment.ID)
	suite.Nil(err)

	shipmentPayload := payloadForShipmentModel(*reloadShipment)
	expectedTspID := *handlers.FmtUUID(tsp.ID)
	suite.Equal(shipmentPayload.TransportationServiceProviderID, expectedTspID)
}

func (suite *HandlerSuite) TestPayloadForShipmentWithNullValues() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	// Set all values that can be nil to nil
	shipment.SourceGBLOC = nil
	shipment.DestinationGBLOC = nil
	shipment.GBLNumber = nil
	shipment.Market = nil
	shipment.TrafficDistributionListID = nil
	shipment.TrafficDistributionList = nil
	shipment.ActualPickupDate = nil
	shipment.ActualPackDate = nil
	shipment.ActualDeliveryDate = nil
	shipment.BookDate = nil
	shipment.RequestedPickupDate = nil
	shipment.OriginalDeliveryDate = nil
	shipment.OriginalPackDate = nil
	shipment.EstimatedPackDays = nil
	shipment.EstimatedTransitDays = nil
	shipment.PickupAddressID = nil
	shipment.PickupAddress = nil
	shipment.SecondaryPickupAddressID = nil
	shipment.SecondaryPickupAddress = nil
	shipment.DeliveryAddressID = nil
	shipment.DeliveryAddress = nil
	shipment.PartialSITDeliveryAddressID = nil
	shipment.PartialSITDeliveryAddress = nil
	shipment.PmSurveyConductedDate = nil
	shipment.PmSurveyCompletedAt = nil
	shipment.PmSurveyPlannedPackDate = nil
	shipment.PmSurveyPlannedPickupDate = nil
	shipment.PmSurveyPlannedDeliveryDate = nil
	shipment.PmSurveyWeightEstimate = nil
	shipment.PmSurveyProgearWeightEstimate = nil
	shipment.PmSurveySpouseProgearWeightEstimate = nil
	shipment.PmSurveyNotes = nil

	shipmentPayload := payloadForShipmentModel(shipment)
	suite.Nil(shipmentPayload.SourceGbloc)
	suite.Nil(shipmentPayload.DestinationGbloc)
	suite.Nil(shipmentPayload.GblNumber)
	suite.Nil(shipmentPayload.Market)
	suite.Nil(shipmentPayload.TrafficDistributionList)
	suite.Nil(shipmentPayload.ActualPickupDate)
	suite.Nil(shipmentPayload.ActualPackDate)
	suite.Nil(shipmentPayload.ActualDeliveryDate)
	suite.Nil(shipmentPayload.BookDate)
	suite.Nil(shipmentPayload.RequestedPickupDate)
	suite.Nil(shipmentPayload.OriginalDeliveryDate)
	suite.Nil(shipmentPayload.OriginalPackDate)
	suite.Nil(shipmentPayload.EstimatedPackDays)
	suite.Nil(shipmentPayload.EstimatedTransitDays)
	suite.Nil(shipmentPayload.PickupAddress)
	suite.Nil(shipmentPayload.SecondaryPickupAddress)
	suite.Nil(shipmentPayload.DeliveryAddress)
	suite.Nil(shipmentPayload.PartialSitDeliveryAddress)
	suite.Nil(shipmentPayload.PmSurveyConductedDate)
	suite.Nil(shipmentPayload.PmSurveyCompletedAt)
	suite.Nil(shipmentPayload.PmSurveyPlannedPackDate)
	suite.Nil(shipmentPayload.PmSurveyPlannedPickupDate)
	suite.Nil(shipmentPayload.PmSurveyPlannedDeliveryDate)
	suite.Nil(shipmentPayload.PmSurveyWeightEstimate)
	suite.Nil(shipmentPayload.PmSurveyProgearWeightEstimate)
	suite.Nil(shipmentPayload.PmSurveySpouseProgearWeightEstimate)
	suite.Nil(shipmentPayload.PmSurveyNotes)
}

func (suite *HandlerSuite) TestGetShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := GetShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.GetShipmentOK{}, response)
	okResponse := response.(*shipmentop.GetShipmentOK)

	// And: Payload is equivalent to original shipment
	suite.Equal(strfmt.UUID(shipment.ID.String()), okResponse.Payload.ID)
	suite.Equal(apimessages.AffiliationARMY, *okResponse.Payload.ServiceMember.Affiliation)
}

func (suite *HandlerSuite) TestGetShipmentHandlerWhereSessionServiceMemberIDDoesNotEqualShipmentServiceMemberID() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	params := shipmentop.GetShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	handler := GetShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.NotEqual(session.ServiceMemberID, shipment.ServiceMemberID)
	suite.Assertions.IsType(&shipmentop.GetShipmentForbidden{}, response)
}

func (suite *HandlerSuite) TestPatchShipmentHandlerNetWeight() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.Shipment{
		NetWeight:   swag.Int64(17500),
		GrossWeight: swag.Int64(12500),
	}

	params := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Update:      &UpdatePayload,
	}

	// And: patch shipment is returned
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentOK{}, response)
	okResponse := response.(*shipmentop.PatchShipmentOK)

	// And: Payload has new values
	suite.Equal(int64(17500), *okResponse.Payload.NetWeight)
	suite.Equal(int64(12500), *okResponse.Payload.GrossWeight)
}

func (suite *HandlerSuite) TestPatchShipmentHandlerPmSurvey() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	calendar := dates.NewUSCalendar()
	genericDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	UpdatePayload := apimessages.Shipment{
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(&genericDate),
		PmSurveyConductedDate:               handlers.FmtDatePtr(&genericDate),
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
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
	suite.Equal(genericDate, *(*time.Time)(okResponse.Payload.PmSurveyConductedDate))
}

func (suite *HandlerSuite) TestPatchShipmentHandlerPmSurveyWrongTSP() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	shipment := shipments[0]

	otherTspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	// And: the context contains the auth values for the wrong tsp
	req := httptest.NewRequest("PATCH", "/shipments/shipmentId", nil)
	req = suite.AuthenticateTspRequest(req, otherTspUser)

	calendar := dates.NewUSCalendar()
	genericDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	UpdatePayload := apimessages.Shipment{
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(&genericDate),
		PmSurveyConductedDate:               handlers.FmtDatePtr(&genericDate),
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
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&shipmentop.PatchShipmentBadRequest{}, response)
}

// TestIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

	// When: There is the full set of data for a GBL to be generated successfully
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	unauthedTSPUser := testdatagen.MakeTspUser(suite.DB(), testdatagen.Assertions{
		TspUser: models.TspUser{
			Email: "unauthorized@example.com",
		},
		User: models.User{
			LoginGovEmail: "unauthorized@example.com",
		},
	})
	shipment := scenario.MakeHhgFromAwardedToAcceptedGBLReady(suite.DB(), tspUser)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.CreateGovBillOfLadingParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetFileStorer(fakeS3)

	// And: the Orders are missing required data
	shipment.Move.Orders.TAC = nil
	suite.MustSave(&shipment.Move.Orders)

	formCreator := paperworkservice.NewFormCreator(context.FileStorer().TempFileSystem(), paperwork.NewFormFiller())
	// And: the create gbl handler is called
	handler := CreateGovBillOfLadingHandler{context, formCreator}
	response := handler.Handle(params)

	// Then: expect a 417 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)
	suite.Assertions.Equal(http.StatusExpectationFailed, errResponse.Code)

	// When: the Orders have all required data
	shipment.Move.Orders.TAC = models.StringPointer("NTA4")
	suite.MustSave(&shipment.Move.Orders)

	// And: the create gbl handler is called
	handler = CreateGovBillOfLadingHandler{context, formCreator}
	response = handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&shipmentop.CreateGovBillOfLadingCreated{}, response)

	// When: there is an existing GBL for a shipment and handler is called
	handler = CreateGovBillOfLadingHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger()), formCreator}
	response = handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse = response.(*handlers.ErrResponse)
	suite.Assertions.Equal(http.StatusBadRequest, errResponse.Code)

	// When: an unauthed TSP user hits the handler
	req = suite.AuthenticateTspRequest(req, unauthedTSPUser)
	params.HTTPRequest = req
	response = handler.Handle(params)

	// Then: expect a 400 status code
	suite.CheckResponseForbidden(response)
}

// TestIndexShipmentsHandlerPaginated tests the api endpoint with pagination query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser1 := tspUsers[0]
	tspUser2 := tspUsers[1]

	// Constants
	limit := int64(25)
	offset := int64(1)

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

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
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
			pickupDate = time.Time(*responsePayload.ActualPickupDate)
		} else {
			newDT := time.Time(*responsePayload.ActualPickupDate)
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
	status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
			pickupDate = time.Time(*responsePayload.ActualPickupDate)
		} else {
			newDT := time.Time(*responsePayload.ActualPickupDate)
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
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
			deliveryDate = time.Time(*responsePayload.ActualDeliveryDate)
		} else {
			newDT := time.Time(*responsePayload.ActualDeliveryDate)
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
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
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
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
			deliveryDate = time.Time(*responsePayload.ActualDeliveryDate)
		} else {
			newDT := time.Time(*responsePayload.ActualDeliveryDate)
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
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	// The params expect statuses in strings, so they have to be cast from ShipmentStatus types
	stringStatus := []string{string(models.ShipmentStatusSUBMITTED)}
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
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Handler to Test
	handler := AcceptShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

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

// TestTransportShipmentHandler tests the api endpoint that transports a shipment
func (suite *HandlerSuite) TestTransportShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAPPROVED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Handler to Test
	handler := TransportShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	// Test query with first user
	path := fmt.Sprintf("/shipments/%s/transport", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	actualPackDate := shipment.BookDate.AddDate(0, 0, 1)
	actualPickupDate := actualPackDate.AddDate(0, 0, 1)

	body := apimessages.TransportPayload{
		ActualPackDate:   handlers.FmtDatePtr(&actualPackDate),
		ActualPickupDate: handlers.FmtDatePtr(&actualPickupDate),
		NetWeight:        swag.Int64(2000),
		GrossWeight:      swag.Int64(3000),
		TareWeight:       swag.Int64(1000),
	}
	params := shipmentop.TransportShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
		Payload:     &body,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.TransportShipmentOK{}, response)
	okResponse := response.(*shipmentop.TransportShipmentOK)
	suite.Equal("IN_TRANSIT", string(okResponse.Payload.Status))
	suite.Equal(actualPickupDate, time.Time(*okResponse.Payload.ActualPickupDate))
	suite.Equal(actualPackDate, time.Time(*okResponse.Payload.ActualPackDate))
	suite.Equal(int64(2000), *okResponse.Payload.NetWeight)
	suite.Equal(int64(3000), *okResponse.Payload.GrossWeight)
	suite.Equal(int64(1000), *okResponse.Payload.TareWeight)
}

// TestDeliverShipmentHandler tests the api endpoint that delivers a shipment
func (suite *HandlerSuite) TestDeliverShipmentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Add a line item that's ready to be priced
	preApproval := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// Make sure there's a FuelEIADieselPrice to use
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// Handler to Test
	handler := DeliverShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.SetPlanner(route.NewTestingPlanner(1044))

	// Test query with first user
	path := fmt.Sprintf("/shipments/%s/deliver", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	calendar := dates.NewUSCalendar()
	actualDeliveryDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	body := apimessages.ActualDeliveryDate{
		ActualDeliveryDate: handlers.FmtDatePtr(&actualDeliveryDate),
	}
	params := shipmentop.DeliverShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
		Payload:     &body,
	}

	response := handler.Handle(params)
	suite.CheckNotErrorResponse(response)
	suite.Assertions.IsType(&shipmentop.DeliverShipmentOK{}, response)
	okResponse := response.(*shipmentop.DeliverShipmentOK)
	suite.Equal("DELIVERED", string(okResponse.Payload.Status))
	suite.Equal(actualDeliveryDate, time.Time(*okResponse.Payload.ActualDeliveryDate))

	// Check for ShipmentLineItems
	addedLineItems, _ := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)

	// The details of the line items are tested in the rateengine package.  We just
	// check the count here.
	suite.Len(addedLineItems, 7)

	updatedPreApproval, err := models.FetchShipmentLineItemByID(suite.DB(), &preApproval.ID)
	if suite.NoError(err) {
		suite.NotNil(updatedPreApproval.AmountCents)
	}
}

// TestCompletePmSurveyHandler tests the api endpoint that saves a shipment's Pm Survey
func (suite *HandlerSuite) TestCompletePmSurveyHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAPPROVED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Handler to Test
	handler := CompletePmSurveyHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	// Test query with first user
	path := fmt.Sprintf("/shipments/%s/completePmSurvey", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.CompletePmSurveyParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.CompletePmSurveyOK{}, response)
}

// TestGetShipmentInvoicesHandler tests the api endpoint that saves a shipment's Pm Survey
func (suite *HandlerSuite) TestGetShipmentInvoicesHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusAPPROVED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Given 2 invoices
	testdatagen.MakeInvoice(suite.DB(), testdatagen.Assertions{
		Invoice: models.Invoice{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
		},
	})
	testdatagen.MakeInvoice(suite.DB(), testdatagen.Assertions{
		Invoice: models.Invoice{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
		},
	})

	// Handler to Test
	handler := GetShipmentInvoicesHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	path := fmt.Sprintf("/shipments/%s/invoices", shipment.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	// An unauthorized TSP user should be forbidden
	otherTspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	req = suite.AuthenticateTspRequest(req, otherTspUser)
	params := shipmentop.GetShipmentInvoicesParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
	}
	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.GetShipmentInvoicesForbidden{}, response)

	// A valid TSP user should retrieve invoices
	req = suite.AuthenticateTspRequest(req, tspUser)
	params = shipmentop.GetShipmentInvoicesParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
	}

	response = handler.Handle(params)
	if suite.Assertions.IsType(&shipmentop.GetShipmentInvoicesOK{}, response) {
		okResponse := response.(*shipmentop.GetShipmentInvoicesOK)
		suite.Len(okResponse.Payload, 2)
	}

	// A valid office user should retrieve invoices
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params = shipmentop.GetShipmentInvoicesParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(shipment.ID),
	}

	response = handler.Handle(params)
	if suite.Assertions.IsType(&shipmentop.GetShipmentInvoicesOK{}, response) {
		okResponse := response.(*shipmentop.GetShipmentInvoicesOK)
		suite.Len(okResponse.Payload, 2)
	}
}
