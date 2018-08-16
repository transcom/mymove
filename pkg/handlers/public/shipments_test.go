package public

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// TestIndexShipmentsHandlerAllShipments tests the api endpoint with no query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerAllShipments() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []string{"DEFAULT"}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)

	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
	}

	// And: an index of shipments is returned
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.parent.Equal(1, len(okResponse.Payload))
	if len(okResponse.Payload) == 1 {
		responsePayload := okResponse.Payload[0]
		// And: Payload is equivalent to original shipment
		suite.parent.Equal(strfmt.UUID(shipment.ID.String()), responsePayload.ID)
		suite.parent.Equal(apimessages.SelectedMoveType(*shipment.Move.SelectedMoveType), *responsePayload.Move.SelectedMoveType)
		suite.parent.Equal(shipment.TrafficDistributionList.SourceRateArea, *responsePayload.TrafficDistributionList.SourceRateArea)
	}
}

// TestIndexShipmentsHandlerPaginated tests the api endpoint with pagination query parameters
func (suite *HandlerSuite) TestIndexShipmentsHandlerPaginated() {

	numTspUsers := 2
	numShipments := 25
	numShipmentOfferSplit := []int{15, 10}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser1 := tspUsers[0]
	tspUser2 := tspUsers[1]

	// Constants
	limit := int64(25)
	offset := int64(1)

	// Handler to Test
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))

	// Test query with first user
	req1 := httptest.NewRequest("GET", "/shipments", nil)
	req1 = suite.parent.AuthenticateTspRequest(req1, tspUser1)
	params1 := shipmentop.IndexShipmentsParams{
		HTTPRequest: req1,
		Limit:       &limit,
		Offset:      &offset,
	}

	response1 := handler.Handle(params1)
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response1)
	okResponse1 := response1.(*shipmentop.IndexShipmentsOK)
	suite.parent.Equal(15, len(okResponse1.Payload))

	// Test query with second user
	req2 := httptest.NewRequest("GET", "/shipments", nil)
	req2 = suite.parent.AuthenticateTspRequest(req2, tspUser2)
	params2 := shipmentop.IndexShipmentsParams{
		HTTPRequest: req2,
		Limit:       &limit,
		Offset:      &offset,
	}

	response2 := handler.Handle(params2)
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response2)
	okResponse2 := response2.(*shipmentop.IndexShipmentsOK)
	suite.parent.Equal(10, len(okResponse2.Payload))
}

// TestIndexShipmentsHandlerSortShipmentsPickupAsc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsPickupAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.parent.Equal(3, len(okResponse.Payload))

	var pickupDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if pickupDate == empty {
			pickupDate = time.Time(responsePayload.PickupDate)
		} else {
			newDT := time.Time(responsePayload.PickupDate)
			suite.parent.True(newDT.After(pickupDate))
			pickupDate = newDT
		}
	}
}

// TestIndexShipmentsHandlerSortShipmentsPickupDesc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsPickupDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.parent.Equal(3, len(okResponse.Payload))

	var pickupDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if pickupDate == empty {
			pickupDate = time.Time(responsePayload.PickupDate)
		} else {
			newDT := time.Time(responsePayload.PickupDate)
			suite.parent.True(newDT.Before(pickupDate))
			pickupDate = newDT
		}
	}
}

// TestIndexShipmentsHandlerSortShipmentsDeliveryAsc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsDeliveryAsc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.parent.Equal(3, len(okResponse.Payload))

	var deliveryDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if deliveryDate == empty {
			deliveryDate = time.Time(responsePayload.DeliveryDate)
		} else {
			newDT := time.Time(responsePayload.DeliveryDate)
			suite.parent.True(newDT.After(deliveryDate))
			deliveryDate = newDT
		}
	}
}

// TestIndexShipmentsHandlerSortShipmentsDeliveryDesc sorts returned shipments
func (suite *HandlerSuite) TestIndexShipmentsHandlerSortShipmentsDeliveryDesc() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)

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
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)

	// And: Returned query to have at least one shipment in the list
	suite.parent.Equal(3, len(okResponse.Payload))

	var deliveryDate time.Time
	empty := time.Time{}
	for _, responsePayload := range okResponse.Payload {
		if deliveryDate == empty {
			deliveryDate = time.Time(responsePayload.DeliveryDate)
		} else {
			newDT := time.Time(responsePayload.DeliveryDate)
			suite.parent.True(newDT.Before(deliveryDate))
			deliveryDate = newDT
		}
	}
}

// TestIndexShipmentsHandlerFilterByStatus tests the api endpoint with defined status query param
func (suite *HandlerSuite) TestIndexShipmentsHandlerFilterByStatus() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      status,
	}

	response := handler.Handle(params)
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)
	suite.parent.Equal(25, len(okResponse.Payload))
}

// TestIndexShipmentsHandlerFilterByStatusNoResults tests the api endpoint with defined status query param that returns nothing
func (suite *HandlerSuite) TestIndexShipmentsHandlerFilterByStatusNoResults() {
	numTspUsers := 1
	numShipments := 25
	numShipmentOfferSplit := []int{25}
	status := []string{"DEFAULT"}
	tspUsers, _, _, err := testdatagen.CreateShipmentOfferData(suite.parent.Db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	if err != nil {
		fmt.Println(err)
	}

	tspUser := tspUsers[0]

	// Handler to Test
	handler := IndexShipmentsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	statusFilter := []string{"NOTASTATUS"}

	// Test query with first user
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.parent.AuthenticateTspRequest(req, tspUser)
	params := shipmentop.IndexShipmentsParams{
		HTTPRequest: req,
		Status:      statusFilter,
	}

	response := handler.Handle(params)
	suite.parent.Assertions.IsType(&shipmentop.IndexShipmentsOK{}, response)
	okResponse := response.(*shipmentop.IndexShipmentsOK)
	suite.parent.Equal(0, len(okResponse.Payload))
}
