package publicapi

import (
	"net/http/httptest"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/apimessages"

	"github.com/go-openapi/strfmt"

	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func makePreApprovalItem(db *pop.Connection) models.Tariff400ngItem {
	item := testdatagen.MakeDefaultTariff400ngItem(db)
	item.RequiresPreApproval = true
	db.Save(&item)
	return item
}

func (suite *HandlerSuite) TestGetShipmentLineItemTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment line items tied to two different shipments
	acc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			ShipmentID: shipment.ID,
		},
	})
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.GetShipmentLineItemsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(acc1.ShipmentID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentLineItemsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.GetShipmentLineItemsOK{}, response) {
		okResponse := response.(*accessorialop.GetShipmentLineItemsOK)

		// And: Payload is equivalent to original shipment line item
		suite.Len(okResponse.Payload, 1)
		suite.Equal(acc1.ID.String(), okResponse.Payload[0].ID.String())
	}
}

func (suite *HandlerSuite) TestGetShipmentLineItemOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	acc1 := testdatagen.MakeDefaultShipmentLineItem(suite.DB())
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.GetShipmentLineItemsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(acc1.ShipmentID.String()),
	}

	// And: get shipment is returned
	handler := GetShipmentLineItemsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.GetShipmentLineItemsOK{}, response) {
		okResponse := response.(*accessorialop.GetShipmentLineItemsOK)

		// And: Payload is equivalent to original shipment line item
		suite.Len(okResponse.Payload, 1)
		suite.Equal(acc1.ID.String(), okResponse.Payload[0].ID.String())
	}
}

func (suite *HandlerSuite) TestGetShipmentLineItemRecalculateHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	acc1 := testdatagen.MakeDefaultShipmentLineItem(suite.DB())
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.GetShipmentLineItemsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(acc1.ShipmentID.String()),
	}

	// Create date range
	recalculateRange := models.ShipmentRecalculate{
		ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
		ShipmentUpdatedBefore: time.Now(),
		Active:                true,
	}
	suite.MustCreate(suite.DB(), &recalculateRange)

	// And: get shipment is returned
	handler := GetShipmentLineItemsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.GetShipmentLineItemsOK{}, response) {
		okResponse := response.(*accessorialop.GetShipmentLineItemsOK)

		// And: Payload is equivalent to original shipment line item
		suite.Len(okResponse.Payload, 1)
		suite.Equal(acc1.ID.String(), okResponse.Payload[0].ID.String())
	}
}

func (suite *HandlerSuite) TestCreateShipmentLineItemHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	tariffItem := makePreApprovalItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	payload := apimessages.ShipmentLineItem{
		Tariff400ngItemID: handlers.FmtUUID(tariffItem.ID),
		Location:          apimessages.ShipmentLineItemLocationORIGIN,
		Notes:             "Some notes",
		Quantity1:         handlers.FmtInt64(int64(5)),
	}

	params := accessorialop.CreateShipmentLineItemParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Payload:     &payload,
	}

	// And: get shipment is returned
	handler := CreateShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.CreateShipmentLineItemCreated{}, response) {
		okResponse := response.(*accessorialop.CreateShipmentLineItemCreated)
		// And: Payload is equivalent to original shipment line
		if suite.NotNil(okResponse.Payload.Notes) {
			suite.Equal("Some notes", okResponse.Payload.Notes)
		}
	}
}

func (suite *HandlerSuite) TestCreateShipmentLineItemForbidden() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	tariffItem := testdatagen.MakeDefaultTariff400ngItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	payload := apimessages.ShipmentLineItem{
		Tariff400ngItemID: handlers.FmtUUID(tariffItem.ID),
		Location:          apimessages.ShipmentLineItemLocationORIGIN,
		Notes:             "Some notes",
		Quantity1:         handlers.FmtInt64(int64(5)),
	}

	params := accessorialop.CreateShipmentLineItemParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Payload:     &payload,
	}

	// And: get shipment is returned
	handler := CreateShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 403 status code
	suite.Assertions.IsType(&accessorialop.CreateShipmentLineItemForbidden{}, response)
}

func (suite *HandlerSuite) TestUpdateShipmentLineItemTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)
	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			ShipmentID: shipment.ID,
			Location:   models.ShipmentLineItemLocationDESTINATION,
			Quantity1:  unit.BaseQuantity(int64(123456)),
			Quantity2:  unit.BaseQuantity(int64(654321)),
			Notes:      "",
		},
	})

	testdatagen.MakeDefaultShipmentLineItem(suite.DB())
	// create a new tariff400ngitem to test
	updateAcc1 := makePreApprovalItem(suite.DB())
	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	updateShipmentLineItem := apimessages.ShipmentLineItem{
		ID:                *handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:        *handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:          apimessages.ShipmentLineItemLocationORIGIN,
		Quantity1:         handlers.FmtInt64(int64(1)),
		Quantity2:         handlers.FmtInt64(int64(2)),
		Notes:             "HELLO",
		Tariff400ngItemID: handlers.FmtUUID(updateAcc1.ID),
	}
	params := accessorialop.UpdateShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
		Payload:            &updateShipmentLineItem,
	}

	// And: get shipment is returned
	handler := UpdateShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.UpdateShipmentLineItemOK{}, response) {
		okResponse := response.(*accessorialop.UpdateShipmentLineItemOK)

		// Payload should match the UpdateShipmentLineItem
		suite.Equal(updateShipmentLineItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(updateShipmentLineItem.ShipmentID.String(), okResponse.Payload.ShipmentID.String())
		suite.Equal(updateShipmentLineItem.Location, okResponse.Payload.Location)
		suite.Equal(*updateShipmentLineItem.Quantity1, *okResponse.Payload.Quantity1)
		suite.Equal(*updateShipmentLineItem.Quantity2, *okResponse.Payload.Quantity2)
		suite.Equal(updateShipmentLineItem.Notes, okResponse.Payload.Notes)
		suite.Equal(updateShipmentLineItem.Tariff400ngItemID.String(), okResponse.Payload.Tariff400ngItemID.String())
	}
}

func (suite *HandlerSuite) TestUpdateShipmentLineItemOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Location:  models.ShipmentLineItemLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
	})
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// create a new tariff400ngItem to test
	updateAcc1 := makePreApprovalItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	updateShipmentLineItem := apimessages.ShipmentLineItem{
		ID:                *handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:        *handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:          apimessages.ShipmentLineItemLocationORIGIN,
		Quantity1:         handlers.FmtInt64(int64(1)),
		Quantity2:         handlers.FmtInt64(int64(2)),
		Notes:             "HELLO",
		Tariff400ngItemID: handlers.FmtUUID(updateAcc1.ID),
	}
	params := accessorialop.UpdateShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
		Payload:            &updateShipmentLineItem,
	}

	// And: get shipment is returned
	handler := UpdateShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.UpdateShipmentLineItemOK{}, response) {
		okResponse := response.(*accessorialop.UpdateShipmentLineItemOK)

		// Payload should match the UpdateShipmentLineItem
		suite.Equal(updateShipmentLineItem.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(updateShipmentLineItem.ShipmentID.String(), okResponse.Payload.ShipmentID.String())
		suite.Equal(updateShipmentLineItem.Location, okResponse.Payload.Location)
		suite.Equal(*updateShipmentLineItem.Quantity1, *okResponse.Payload.Quantity1)
		suite.Equal(*updateShipmentLineItem.Quantity2, *okResponse.Payload.Quantity2)
		suite.Equal(updateShipmentLineItem.Notes, okResponse.Payload.Notes)
		suite.Equal(updateShipmentLineItem.Tariff400ngItemID.String(), okResponse.Payload.Tariff400ngItemID.String())
	}
}

func (suite *HandlerSuite) TestUpdateShipmentLineItemForbidden() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Location:  models.ShipmentLineItemLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
	})
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// create a new tariff400ngItem to test
	updateAcc1 := testdatagen.MakeDefaultTariff400ngItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PUT", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	updateShipmentLineItem := apimessages.ShipmentLineItem{
		ID:                *handlers.FmtUUID(shipAcc1.ID),
		ShipmentID:        *handlers.FmtUUID(shipAcc1.ShipmentID),
		Location:          apimessages.ShipmentLineItemLocationORIGIN,
		Quantity1:         handlers.FmtInt64(int64(1)),
		Quantity2:         handlers.FmtInt64(int64(2)),
		Notes:             "HELLO",
		Tariff400ngItemID: handlers.FmtUUID(updateAcc1.ID),
	}
	params := accessorialop.UpdateShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
		Payload:            &updateShipmentLineItem,
	}

	// And: get shipment is returned
	handler := UpdateShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 403 status code
	suite.Assertions.IsType(&accessorialop.UpdateShipmentLineItemForbidden{}, response)
}

func (suite *HandlerSuite) TestDeleteShipmentLineItemTSPHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			ShipmentID: shipment.ID,
			Location:   models.ShipmentLineItemLocationDESTINATION,
			Quantity1:  unit.BaseQuantity(int64(123456)),
			Quantity2:  unit.BaseQuantity(int64(654321)),
			Notes:      "",
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.DeleteShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.DeleteShipmentLineItemOK{}, response) {
		// Check if we actually deleted the shipment line
		err = suite.DB().Find(&shipAcc1, shipAcc1.ID)
		suite.Error(err)
	}
}

func (suite *HandlerSuite) TestDeleteShipmentLineItemCode105BE() {

	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	acc105B := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			ShipmentID: shipment.ID,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                "105B",
			RequiresPreApproval: true,
		},
	})

	notes := "It's a giant moose head named Fred he seemed rather pleasant"
	baseParams := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   acc105B.ID,
		Tariff400ngItemCode: acc105B.Code,
		Location:            "ORIGIN",
		Notes:               &notes,
	}
	additionalParams := models.AdditionalShipmentLineItemParams{
		ItemDimensions: &models.AdditionalLineItemDimensions{
			Length: 100,
			Width:  100,
			Height: 100,
		},
		CrateDimensions: &models.AdditionalLineItemDimensions{
			Length: 100,
			Width:  100,
			Height: 100,
		},
	}
	// Given: Create 105B preapproval
	shipmentLineItem, _, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.DeleteShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipmentLineItem.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code, and Crate Dimensions item should be deleted.
	if suite.Assertions.IsType(&accessorialop.DeleteShipmentLineItemOK{}, response) {
		// Check if we actually deleted the shipment line
		err = suite.DB().Find(&shipmentLineItem.CrateDimensions, shipmentLineItem.CrateDimensions.ID)
		suite.Error(err)
	}
}

func (suite *HandlerSuite) TestDeleteShipmentLineItemOfficeHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Location:  models.ShipmentLineItemLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.DeleteShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.DeleteShipmentLineItemOK{}, response) {
		// Check if we actually deleted the shipment line item
		err := suite.DB().Find(&shipAcc1, shipAcc1.ID)
		suite.Error(err)
	}
}

func (suite *HandlerSuite) TestDeleteShipmentLineItemWithoutPreapprovalForbidden() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Two shipment line items tied to two different shipments
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Location:  models.ShipmentLineItemLocationDESTINATION,
			Quantity1: unit.BaseQuantity(int64(123456)),
			Quantity2: unit.BaseQuantity(int64(654321)),
			Notes:     "",
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: false,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.DeleteShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
	}

	// And: get shipment is returned
	handler := DeleteShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 403 status code
	suite.Assertions.IsType(&accessorialop.DeleteShipmentLineItemForbidden{}, response)
}

func (suite *HandlerSuite) TestDeleteShipmentLineItemWithInvoiceBadRequest() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// A ShipmentLineItem tied to an invoice
	invoice := testdatagen.MakeDefaultInvoice(suite.DB())
	shipAcc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			InvoiceID: &invoice.ID,
			Status:    models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("DELETE", "/shipments", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.DeleteShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(shipAcc1.ID.String()),
	}

	handler := DeleteShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.CheckResponseBadRequest(response)
}

func (suite *HandlerSuite) TestApproveShipmentLineItemHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// A shipment line item with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(acc1.ID.String()),
	}

	// And: get shipment is returned
	handler := ApproveShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.ApproveShipmentLineItemOK{}, response) {
		okResponse := response.(*accessorialop.ApproveShipmentLineItemOK)

		// And: Payload is equivalent to original shipment line item
		suite.Equal(acc1.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(apimessages.ShipmentLineItemStatusAPPROVED, okResponse.Payload.Status)
	}
}

func (suite *HandlerSuite) TestApproveShipmentLineItemHandlerShipmentDelivered() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	rateCents := unit.Cents(1000)
	item := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status:    models.ShipmentLineItemStatusSUBMITTED,
			Shipment:  shipments[0],
			Quantity1: unit.BaseQuantityFromInt(1),
		},
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                "130A",
			RequiresPreApproval: true,
			DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		},
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			RateCents: rateCents,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(item.ID.String()),
	}

	// And: get shipment line item is returned
	handler := ApproveShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	if suite.Assertions.IsType(&accessorialop.ApproveShipmentLineItemOK{}, response) {
		okResponse := response.(*accessorialop.ApproveShipmentLineItemOK)

		// And: Payload is equivalent to original shipment line item
		suite.Equal(item.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(apimessages.ShipmentLineItemStatusAPPROVED, okResponse.Payload.Status)

		// And: Rate and amount have been assigned to item
		suite.NotNil(okResponse.Payload.AppliedRate)
		if suite.NotNil(okResponse.Payload.AmountCents) {
			discountRate := shipments[0].ShipmentOffers[0].TransportationServiceProviderPerformance.LinehaulRate
			suite.NotEqual(discountRate, 0)
			// There should be a discount rate applied for code 130A
			suite.Equal(int64(discountRate.Apply(rateCents)), *okResponse.Payload.AmountCents)
		}
	}
}

func (suite *HandlerSuite) TestApproveShipmentLineItemNotRequired() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// A shipment line item with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: false,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect user to be forbidden from approving an item that doesn't require pre-approval
	suite.Assertions.IsType(&accessorialop.ApproveShipmentLineItemForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveShipmentLineItemAlreadyApproved() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// A shipment line item with an item that requires pre-approval
	acc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status: models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := accessorialop.ApproveShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect user to be forbidden from approving an item that is already approved
	suite.Assertions.IsType(&accessorialop.ApproveShipmentLineItemForbidden{}, response)
}

func (suite *HandlerSuite) TestApproveShipmentLineItemTSPUser() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// A shipment line item claimed by the tspUser's TSP, and item requires pre-approval
	acc1 := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			ShipmentID: shipment.ID,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/accessorials/some_id/approve", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := accessorialop.ApproveShipmentLineItemParams{
		HTTPRequest:        req,
		ShipmentLineItemID: strfmt.UUID(acc1.ID.String()),
	}

	handler := ApproveShipmentLineItemHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect TSP user to be forbidden from approving
	suite.Assertions.IsType(&accessorialop.ApproveShipmentLineItemForbidden{}, response)
}
