package publicapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/apimessages"

	"github.com/go-openapi/strfmt"

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

func (suite *HandlerSuite) TestCreateShipmentAccessorialHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// Two shipment accessorials tied to two different shipments
	shipment := testdatagen.MakeDefaultShipment(suite.TestDB())
	acc := testdatagen.MakeDefaultTariff400ngItem(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	payload := apimessages.ShipmentAccessorial{
		AccessorialID: handlers.FmtUUID(acc.ID),
		Location:      apimessages.ShipmentAccessorialLocationORIGIN,
		Notes:         "Some notes",
		Quantity1:     handlers.FmtInt64(int64(5)),
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
		suite.Equal("Some notes", okResponse.Payload.Notes)
	}
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
	updateAcc1 := testdatagen.MakeDefaultTariff400ngItem(suite.TestDB())
	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	updateShipmentAccessorial := apimessages.ShipmentAccessorial{
		ID:            *handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:    *handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:      apimessages.ShipmentAccessorialLocationORIGIN,
		Quantity1:     handlers.FmtInt64(int64(1)),
		Quantity2:     handlers.FmtInt64(int64(2)),
		Notes:         "HELLO",
		AccessorialID: handlers.FmtUUID(updateAcc1.ID),
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
	suite.Equal(updateShipmentAccessorial.AccessorialID.String(), okResponse.Payload.AccessorialID.String())
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
	updateAcc1 := testdatagen.MakeDefaultTariff400ngItem(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	updateShipmentAccessorial := apimessages.ShipmentAccessorial{
		ID:            *handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:    *handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:      apimessages.ShipmentAccessorialLocationORIGIN,
		Quantity1:     handlers.FmtInt64(int64(1)),
		Quantity2:     handlers.FmtInt64(int64(2)),
		Notes:         "HELLO",
		AccessorialID: handlers.FmtUUID(updateAcc1.ID),
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
	suite.Equal(updateShipmentAccessorial.AccessorialID.String(), okResponse.Payload.AccessorialID.String())
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
	suite.Error(err)
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
	suite.Error(err)
}

func (suite *HandlerSuite) TestApproveShipmentAccessorialHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// A shipment accessorial with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(acc1.ID.String()),
	}

	// And: get shipment is returned
	handler := ApproveShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.ApproveShipmentAccessorialOK{}, response)
	okResponse := response.(*accessorialop.ApproveShipmentAccessorialOK)

	// And: Payload is equivalent to original shipment accessorial
	suite.Equal(acc1.ID.String(), okResponse.Payload.ID.String())
	suite.Equal(apimessages.AccessorialStatusAPPROVED, okResponse.Payload.Status)
}

func (suite *HandlerSuite) TestApproveShipmentAccessorialNotRequired() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// A shipment accessorial with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: false,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect user to be forbidden from approving an item that doesn't require pre-approval
	suite.Assertions.IsType(&accessorialop.ApproveShipmentAccessorialForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveShipmentAccessorialAlreadyApproved() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.TestDB())

	// A shipment accessorial with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			Status: models.ShipmentAccessorialStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect user to be forbidden from approving an item that is already approved
	suite.Assertions.IsType(&accessorialop.ApproveShipmentAccessorialForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveShipmentAccessorialTSPUser() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// A shipment accessorial claimed by the tspUser's TSP, and item requires pre-approval
	acc1 := testdatagen.MakeShipmentAccessorial(suite.TestDB(), testdatagen.Assertions{
		ShipmentAccessorial: models.ShipmentAccessorial{
			ShipmentID: shipment.ID,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.ApproveShipmentAccessorialParams{
		HTTPRequest:           req,
		ShipmentAccessorialID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentAccessorialHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect TSP user to be forbidden from approving
	suite.Assertions.IsType(&accessorialop.ApproveShipmentAccessorialForbidden{}, response)
}
