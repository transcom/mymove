package publicapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
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

func (suite *HandlerSuite) TestUpdateShipmentAccessorialTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment accessorials tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			ShipmentID: shipment.ID,
			Location:   models.ShipmentAccessorialLocationDESTINATION,
			Quantity1:  unit.BaseQuantity(int64(123456)),
			Quantity2:  unit.BaseQuantity(int64(654321)),
			Notes:      "",
		},
	})
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// create a new accessorial to test
	updateAcc1 := testdatagen.MakeDummyAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	updateShipmentAccessorial := apimessages.ShipmentAccessorial{
		ID:          handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:  handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:    apimessages.AccessorialLocationORIGIN,
		Quantity1:   handlers.FmtInt64(int64(1)),
		Quantity2:   handlers.FmtInt64(int64(2)),
		Notes:       "HELLO",
		Accessorial: payloadForAccessorialModel(&updateAcc1),
	}

	params := accessorialop.UpdateShipmentAccessorialParams{
		HTTPRequest:               req,
		ShipmentAccessorialID:     strfmt.UUID(shipAcc1.ID.String()),
		UpdateShipmentAccessorial: &updateShipmentAccessorial,
	}

	// And: get shipment is returned
	handler := UpdateShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.UpdateShipmentAccessorialOK{}, response)
	okResponse := response.(*accessorialop.UpdateShipmentAccessorialOK)

	// Payload should match the UpdateShipmentAccesorial
	suite.Equal(updateShipmentAccessorial.ID.String(), okResponse.Payload.ID.String())
	suite.Equal(updateShipmentAccessorial.ShipmentID.String(), okResponse.Payload.ShipmentID.String())
	suite.Equal(updateShipmentAccessorial.Location, okResponse.Payload.Location)
	suite.Equal(*updateShipmentAccessorial.Quantity1, *okResponse.Payload.Quantity1)
	suite.Equal(*updateShipmentAccessorial.Quantity2, *okResponse.Payload.Quantity2)
	suite.Equal(updateShipmentAccessorial.Notes, okResponse.Payload.Notes)
	suite.Equal(updateShipmentAccessorial.Accessorial.ID.String(), okResponse.Payload.Accessorial.ID.String())
}

func (suite *HandlerSuite) TestUpdateShipmentAccessorialOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// Two shipment accessorials tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			Location:  models.ShipmentAccessorialLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
	})
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// create a new accessorial to test
	updateAcc1 := testdatagen.MakeDummyAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	updateShipmentAccessorial := apimessages.ShipmentAccessorial{
		ID:          handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:  handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:    apimessages.AccessorialLocationORIGIN,
		Quantity1:   handlers.FmtInt64(int64(1)),
		Quantity2:   handlers.FmtInt64(int64(2)),
		Notes:       "HELLO",
		Accessorial: payloadForAccessorialModel(&updateAcc1),
	}

	params := accessorialop.UpdateShipmentAccessorialParams{
		HTTPRequest:               req,
		ShipmentAccessorialID:     strfmt.UUID(shipAcc1.ID.String()),
		UpdateShipmentAccessorial: &updateShipmentAccessorial,
	}

	// And: get shipment is returned
	handler := UpdateShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.UpdateShipmentAccessorialOK{}, response)
	okResponse := response.(*accessorialop.UpdateShipmentAccessorialOK)

	// Payload should match the UpdateShipmentAccesorial
	suite.Equal(updateShipmentAccessorial.ID.String(), okResponse.Payload.ID.String())
	suite.Equal(updateShipmentAccessorial.ShipmentID.String(), okResponse.Payload.ShipmentID.String())
	suite.Equal(updateShipmentAccessorial.Location, okResponse.Payload.Location)
	suite.Equal(*updateShipmentAccessorial.Quantity1, *okResponse.Payload.Quantity1)
	suite.Equal(*updateShipmentAccessorial.Quantity2, *okResponse.Payload.Quantity2)
	suite.Equal(updateShipmentAccessorial.Notes, okResponse.Payload.Notes)
	suite.Equal(updateShipmentAccessorial.Accessorial.ID.String(), okResponse.Payload.Accessorial.ID.String())
}

func (suite *HandlerSuite) TestDeleteShipmentAccessorialTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment accessorials tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			ShipmentID: shipment.ID,
			Location:   models.ShipmentAccessorialLocationDESTINATION,
			Quantity1:  unit.BaseQuantity(int64(123456)),
			Quantity2:  unit.BaseQuantity(int64(654321)),
			Notes:      "",
		},
	})
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.DeleteShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(shipAcc1.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.DeleteShipmentAccessorialOK{}, response)

	// Check if we actually deleted the shipment accessorial
	err = suite.TestDB().Find(&shipAcc1, shipAcc1.ID)
	suite.NotNil(err)
}

func (suite *HandlerSuite) TestDeleteShipmentAccessorialOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// Two shipment accessorials tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			Location:  models.ShipmentAccessorialLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
	})
	testdatagen.MakeDefaultShipmentAccessorial(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.DeleteShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(shipAcc1.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.DeleteShipmentAccessorialOK{}, response)

	// Check if we actually deleted the shipment accessorial
	err := suite.TestDB().Find(&shipAcc1, shipAcc1.ID)
	suite.NotNil(err)
}
