package public

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	publicshipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// TestPublicIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"DEFAULT"}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	status := []string{"DEFAULT"}
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsPickupAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsPickupDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsDeliveryAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerSortShipmentsDeliveryDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerFilterByStatus() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

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
func (suite *utils.HandlerSuite) TestPublicIndexShipmentsHandlerFilterByStatusNoResults() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// Handler to Test
	handler := PublicIndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	statusFilter := []string{"NOTASTATUS"}

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.authenticateTspRequest(req, tspUser)
	params := publicshipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      statusFilter,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&publicshipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*publicshipmentop.IndexShipmentsOK)
	suite.Equal(0, len(okResponse.Payload))
}
