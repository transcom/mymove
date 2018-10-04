package publicapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/apimessages"

	"github.com/go-openapi/strfmt"

	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetShipmentAccessorialTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment accessorials tied to two different shipments
	acc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			ShipmentID: shipment.ID,
		},
	})
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.GetShipmentAccessorialsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(acc1.ShipmentID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentAccessorialsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.GetShipmentAccessorialsOK{}, response)
	okResponse := response.(*accessorialop.GetShipmentAccessorialsOK)

	// And: Payload is equivalent to original shipment accessorial
	suite.Len(okResponse.Payload, 1)
	suite.Equal(acc1.ID.String(), okResponse.Payload[0].ID.String())
}

func (suite *HandlerSuite) TestGetShipmentAccessorialOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// Two shipment accessorials tied to two different shipments
	acc1 := testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.GetShipmentAccessorialsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(acc1.ShipmentID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentAccessorialsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.GetShipmentAccessorialsOK{}, response)
	okResponse := response.(*accessorialop.GetShipmentAccessorialsOK)

	// And: Payload is equivalent to original shipment accessorial
	suite.Len(okResponse.Payload, 1)
	suite.Equal(acc1.ID.String(), okResponse.Payload[0].ID.String())
}

func (suite *HandlerSuite) TestCreateShipmentAccessorialHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// Two shipment accessorials tied to two different shipments
	shipment := testdatagen.MakeDefaultShipment(suite.TestDB())
	acc := testdatagen.MakeDummyAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	payload := apimessages.ShipmentAccessorial{
		Accessorial: payloadForAccessorialModel(&acc),
		Location:    apimessages.AccessorialLocationORIGIN,
		Notes:       handlers.FmtString("Some notes"),
		Quantity1:   handlers.FmtInt64(int64(5)),
	}

	params := accessorialop.CreateShipmentAccessorialParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Payload:     &payload,
	}

	// And: get shipment is returned
	handler := CreateShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.CreateShipmentAccessorialCreated{}, response)
	okResponse := response.(*accessorialop.CreateShipmentAccessorialCreated)

	// And: Payload is equivalent to original shipment accessorial
	if suite.NotNil(okResponse.Payload.Notes) {
		suite.Equal("Some notes", *okResponse.Payload.Notes)
	}
}
